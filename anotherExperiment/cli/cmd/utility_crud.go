package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pgvector/pgvector-go"
)

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func AddNewRow(model interface{}, tablename string) (error, error) {
	db, err := ConnectDb()

	if err != nil {
		return nil, ErrorHandler(err, "db conn error")
	}

	defer db.Close()

	stmt, err := db.Prepare(GenerateInsertQuery(tablename, model))

	if err != nil {
		return nil, ErrorHandler(err, "SQL prep statement err")
	}
	defer stmt.Close()
	values := getStructValues(model)

	fmt.Println("Args:", values)
	res, err := stmt.Exec(values...)

	if err != nil {
		//return nil, ErrorHandler(err, "db insertion error")
		if strings.Contains(err.Error(), "Duplicate entry") {
			log.Printf("Skipping duplicate entry: %v", values[0])
			return nil, nil
		}
		return nil, ErrorHandler(err, "db insertion error")
	}

	fmt.Println(res.RowsAffected())
	//lastid, err := res.LastInsertId()

	// if err != nil {
	// 	return nil, ErrorHandler(err, "err getting last id")
	// }

	// fmt.Printf("job successful!  job insertd with id %d ", lastid)

	return nil, nil

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
