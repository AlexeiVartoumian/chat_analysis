package cmd

import (
	"bytes"
	"cli/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func AddNewRow(model interface{}, tablename string) (int, error) {
	db, err := ConnectDb()

	if err != nil {
		return 0, ErrorHandler(err, "db conn error")
	}

	defer db.Close()

	stmt, err := db.Prepare(GenerateInsertQuery(tablename, model))

	if err != nil {
		return 0, ErrorHandler(err, "SQL prep statement err")
	}
	defer stmt.Close()
	values := getStructValues(model)

	fmt.Println("Args:", values)
	res, err := stmt.Exec(values...)

	if err != nil {
		//return nil, ErrorHandler(err, "db insertion error")
		//if strings.Contains(err.Error(), "Duplicate entry") {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			log.Printf("Skipping duplicate entry: %v", values[0])
			return 1, nil
		}
		return 0, ErrorHandler(err, "db insertion error")
	}

	fmt.Println(res.RowsAffected())
	//lastid, err := res.LastInsertId()

	// if err != nil {
	// 	return nil, ErrorHandler(err, "err getting last id")
	// }

	// fmt.Printf("job successful!  job insertd with id %d ", lastid)

	return 0, nil

}

// generic inserter
func GenerateInsertQuery(tableName string, model interface{}) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string
	paramindex := 1 // postgres way of doing things
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		fmt.Println("dbTag", dbTag)
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")

		//if dbTag != "" && dbTag != "job_id" {
		if dbTag != "" {
			if columns != "" {
				columns += ", "
				placeholders += ", "
			}
			columns += dbTag
			//placeholders += "?" mysqlway
			placeholders += fmt.Sprintf("$%d", paramindex)

			paramindex += 1
		}
	}
	//mysql way
	// fmt.Printf("INSERT INTO %s (%s) VALUES (%s)\n", tableName, columns, placeholders)
	// return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)\n", tableName, columns, placeholders)

	//postgresway
	fmt.Printf("INSERT INTO %s (%s) VALUES (%s)\n", tableName, columns, placeholders)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, placeholders)
}

func getStructValues(model interface{}) []interface{} {
	modelValue := reflect.ValueOf(model)
	modelType := modelValue.Type()
	values := []interface{}{}
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")

		if dbTag != "" {
			fmt.Println("Processing ", modelValue.Field(i).Interface())
			values = append(values, modelValue.Field(i).Interface())
		}
	}
	log.Println("Values", values)
	return values
}

// api call to openai embedding model
func GetEmbedding(text string) ([]float32, error) {
	err := godotenv.Load("../../../.env")

	if err != nil {
		return nil, ErrorHandler(err, "env variables did not load for embedding call")
	}
	apiKey := os.Getenv("OPEN_API_KEY")

	body, err := json.Marshal(map[string]string{
		"input": text,
		"model": "text-embedding-3-small",
	})
	if err != nil {
		return nil, ErrorHandler(err, "failed to marshal embedding request")
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(body))
	if err != nil {
		return nil, ErrorHandler(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, ErrorHandler(err, "embedding api call failed")
	}
	defer resp.Body.Close()

	var result EmbeddingResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, ErrorHandler(err, "failed to decode embedding response")
	}
	return result.Data[0].Embedding, nil
}

func BackfillEmbeddings() error {
	db, err := ConnectDb()

	if err != nil {
		return ErrorHandler(err, "db conn error")
	}

	defer db.Close()

	// can do in mem only hundreds of jobs at a time
	rows, err := db.Query("SELECT job_id , job_description FROM job_description WHERE embedding IS NULL")

	if err != nil {
		return ErrorHandler(err, "Failed to query job description")
	}
	defer rows.Close()

	for rows.Next() {
		var jobId int
		var JobDescription string

		if err := rows.Scan(&jobId, &JobDescription); err != nil {
			log.Printf("scan failed for job_id %d , skipping", jobId)
			continue
		}
		embedding, err := GetEmbedding(JobDescription)
		_, err = db.Exec(
			"UPDATE job_description SET embedding = $1 WHERE job_id = $2",
			pgvector.NewVector(embedding), jobId,
		)
		if err != nil {
			log.Printf("update failed for job_id %d , skipping", jobId)
			continue
		}

		log.Printf("embedded job_id %d", jobId)
	}
	return nil
}

func upsertDepartment(db *sql.DB, node models.DeptNode, parentID *int64, companyID int64) error {
	name := strings.TrimSpace(node.Name) // source data has trailing spaces e.g. "Operations "

	_, err := db.Exec(`
		INSERT INTO DEPARTMENT_GREEN (department_id, department_name, parent_id, company_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (department_id) DO UPDATE
		SET department_name = EXCLUDED.department_name,
		    parent_id       = EXCLUDED.parent_id,
		    company_id      = EXCLUDED.company_id
	`, node.ID, name, parentID, companyID)
	if err != nil {
		return fmt.Errorf("upsert department %d (%s): %w", node.ID, name, err)
	}

	for _, child := range node.Children {
		// pass THIS node's own ID as the child's parent — recursion handles
		// arbitrary depth, not just one level of children.
		if err := upsertDepartment(db, child, &node.ID, companyID); err != nil {
			return err
		}
	}
	return nil
}

func upsertTeamAsh(db *sql.DB, node models.DeptNodeAsh, parentID *string, companyID string, departmentID string) error {
	name := strings.TrimSpace(node.Name)

	var externalName *string
	if node.ExternalName != nil {
		trimmed := strings.TrimSpace(*node.ExternalName)
		externalName = &trimmed
	}

	_, err := db.Exec(`
		INSERT INTO TEAM_ASH (team_id, team_name, department_id, parent_team_id, team_external_name, company_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (team_id) DO UPDATE
		SET team_name          = EXCLUDED.team_name,
		    department_id      = EXCLUDED.department_id,
		    parent_team_id     = EXCLUDED.parent_team_id,
		    team_external_name = EXCLUDED.team_external_name,
		    company_id         = EXCLUDED.company_id
	`, node.ID, name, departmentID, parentID, externalName, companyID)
	if err != nil {
		return fmt.Errorf("upsert team %s (%s): %w", node.ID, name, err)
	}
	return nil
}

func walkTeamsAsh(db *sql.DB, nodeID string, byID map[string]models.DeptNodeAsh, childrenOf map[string][]string, companyID, departmentID string) error {
	node := byID[nodeID]

	if err := upsertTeamAsh(db, node, node.ParentTeamId, companyID, departmentID); err != nil {
		return err
	}

	for _, childID := range childrenOf[nodeID] {
		if err := walkTeamsAsh(db, childID, byID, childrenOf, companyID, departmentID); err != nil {
			return err
		}
	}
	return nil
}
