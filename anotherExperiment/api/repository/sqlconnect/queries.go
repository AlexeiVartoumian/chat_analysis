package sqlconnect

import (
	"api/utils"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

type LastThreeDays struct {
	Date_posted      string
	Day_posted       string
	Title            string
	Company_name     string
	Location         string
	Salary           string
	Applicants_count string
	Apply_url        string
	Job_url          string
}

type UrlOnCompanySite struct {
	Date_posted      string
	Title            string
	Company_name     string
	Location         string
	Salary           string
	Applicants_count string
	Apply_url        string
	Job_url          string
}

type QueryResponse struct {
	Columns []ColumnMeta             `json:"columns"`
	Count   int                      `json:"count"`
	Rows    []map[string]interface{} `json:"rows"`
}
type ColumnMeta struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SearchTerm struct {
	Search_term_id int    `json:"search_term_id"`
	Search_term    string `json:"search_term"`
}

// TODO!!! dep injection on db and close db con
func SearchSimilarJobs(query string) error {

	db, err := ConnectDb(1)

	if err != nil {
		return utils.ErrorHandler(err, "db conn error")
	}
	defer db.Close()

	embedding, err := GetEmbedding(query) //api call

	if err != nil {
		return utils.ErrorHandler(err, "failed to embed query")
	}

	rows, err := db.Query(`
    SELECT 
        jd.job_id,
        jd.job_description,
        1 - (jd.embedding <=> $1) AS cosine_similarity,
        j.title,
        c.name AS company_name,
        j.job_url
    FROM job_description jd
    JOIN jobs j ON jd.job_id = j.job_id
    JOIN company c ON j.company_id = c.company_id
    WHERE jd.embedding IS NOT NULL
    ORDER BY jd.embedding <=> $1
    LIMIT 5`,
		pgvector.NewVector(embedding),
	)
	if err != nil {
		return utils.ErrorHandler(err, "similarity search failed")
	}
	defer rows.Close()

	for rows.Next() {
		var JobId int
		var JobDescription, title, companyName, job_url string
		var similarity float64

		rows.Scan(&JobId, &JobDescription, &similarity, &title, &companyName, &job_url)
		log.Printf("%.4f | %s | %s | job_id: %d | View: %s ", similarity, title, companyName, JobId, job_url)
	}

	return nil
}

func LastThreeDaysJobs() ([]LastThreeDays, error) {

	db, err := ConnectDb(1)

	if err != nil {
		return nil, utils.ErrorHandler(err, "db conn error")
	}
	defer db.Close()

	rows, err := db.Query(`
SELECT
		j.date_posted,
		TRIM(TO_CHAR(j.date_posted, 'Day')) as day_posted,
		j.title,
		c.name as company_name,
		j.location,
		COALESCE(j.salary, 'Not Specified') as salary,
		jm.applicants_count,
		jm.company_apply_url,
		j.job_url
	FROM jobs j
	JOIN company c ON j.company_id = c.company_id
	JOIN job_metadata jm ON j.job_id = jm.job_id
	WHERE j.date_posted >= CURRENT_DATE - INTERVAL '3 days'
	AND jm.company_apply_url NOT LIKE 'https://www.linkedin%'
	AND jm.applicants_count NOT LIKE '%Reposted%'
	ORDER BY j.date_posted DESC;
	`)
	if err != nil {
		return nil, utils.ErrorHandler(err, "its funky")
	}
	defer rows.Close()

	var output []LastThreeDays

	for rows.Next() {
		var Out LastThreeDays
		//rows.Scan(&date_posted, &day_posted, &title, &company_name, &location, &salary, &applicants_count, &apply_url, &job_url)
		// log.Printf("DatePosted: %s | %s | %s | %s | %s | %s | View: %s",
		// 	date_posted, day_posted, title, company_name, location, applicants_count, apply_url)

		rows.Scan(&Out.Date_posted, &Out.Day_posted,
			&Out.Title, &Out.Company_name, &Out.Location, &Out.Salary,
			&Out.Applicants_count, &Out.Apply_url, &Out.Job_url)
		// log.Printf("DatePosted: %s | %s | %s | %s | %s | %s | View: %s",
		// 	Out.date_posted, Out.day_posted, Out.title, Out.company_name, Out.location,
		// 	Out.applicants_count, Out.apply_url)

		output = append(output, Out)
	}
	return output, nil
}

func OnlyUrlOnCompanySite() ([]UrlOnCompanySite, error) {

	db, err := ConnectDb(1)

	if err != nil {
		return nil, utils.ErrorHandler(err, "db conn error")
	}
	defer db.Close()

	rows, err := db.Query(`
	SELECT 
    j.date_posted,
    j.title,
    c.name as company_name,
    j.location,
    j.salary,
    jm.applicants_count,
    jm.company_apply_url,
    j.job_url
FROM JOBS j
JOIN COMPANY c ON j.company_id = c.company_id
LEFT JOIN JOB_METADATA jm ON j.job_id = jm.job_id
WHERE j.date_posted IS NOT NULL
  AND j.salary IS NOT NULL
  AND j.salary != 'Not specified'
  AND (jm.company_apply_url NOT LIKE 'https://www.linkedin%'
       OR jm.company_apply_url IS NULL)
ORDER BY j.date_posted DESC;
	`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "hoo hhaa ")
	}
	defer rows.Close()

	var output []UrlOnCompanySite

	for rows.Next() {

		var Out UrlOnCompanySite

		rows.Scan(&Out.Date_posted, &Out.Title, &Out.Company_name,
			&Out.Location, &Out.Salary, &Out.Applicants_count,
			&Out.Apply_url, &Out.Job_url)

		output = append(output, Out)
	}
	return output, nil
}

func SeekExpired() ([]string, error) {

	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "db conn error")
	}
	rows, err := db.Query(`
		SELECT job_id FROM JOB_LIFECYCLE WHERE job_state LIKE 'LISTED' ORDER BY last_seen_listed_at ASC limit 1000;
	`)
	if err != nil {
		return nil, utils.ErrorHandler(err, "yep yep but no")
	}
	defer rows.Close()

	var output []string

	for rows.Next() {

		var res string

		rows.Scan(&res)
		output = append(output, res)
	}
	return output, nil
}

func SeekExpiredAuto(filetype string, firstrun bool) ([]string, error) {

	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "db conn error")
	}

	if firstrun == true {
		_, err := db.Exec(`
		UPDATE JOB_LIFECYCLE SET visited = FALSE`)

		if err != nil {
			return nil, utils.ErrorHandler(err, "first run errored on db output yes")
		}

	}
	var workflow string
	if filetype == "live" {
		workflow = "LISTED"
	} else {
		workflow = "SUSPENDED"
	}
	// UPDATE JOB_LIFECYCLE SET visited = TRUE where job_id IN (SELECT job_id FROM JOB_LIFECYCLE where job_state LIKE 'LISTED' and visited = FALSE limit 1000);
	// order of ops: check if first run with bool flag . if yes then all open are in unvisited state .
	// select 1000 listed roles where listed and unvisited
	// update table as visited
	// send.
	//UPDATE JOB_LIFECYCLE SET visited = TRUE where job_id IN (SELECT job_id FROM JOB_LIFECYCLE where job_state LIKE 'LISTED' and visited = FALSE limit 1000);

	rows, err := db.Query(`
		SELECT job_id FROM JOB_LIFECYCLE WHERE job_state LIKE $1 and visited = FALSE order by last_seen_listed_at ASC limit 1000;
	`, workflow)
	if err != nil {
		return nil, utils.ErrorHandler(err, "yep yep but no")
	}
	defer rows.Close()

	var output []string

	for rows.Next() {

		var res string

		rows.Scan(&res)
		output = append(output, res)
	}

	// then last run has been reached.
	fmt.Println(len(output))
	if len(output) == 0 {

		_, err := db.Exec(`
		UPDATE JOB_LIFECYCLE SET visited = FALSE`)

		if err != nil {
			return nil, utils.ErrorHandler(err, " last run errored on db output yes")
		}
		return nil, nil
	}

	for index := range output {

		// _, err := db.Exec(`
		// UPDATE JOB_LIFECYCLE SET visited = TRUE where job_id IN (SELECT job_id FROM JOB_LIFECYCLE where job_state LIKE 'LISTED' and visited = FALSE
		// `)
		job_id := output[index]
		//fmt.Println(job_id)
		_, err := db.Exec(`
		 UPDATE JOB_LIFECYCLE SET visited = TRUE where job_id = $1
		 `, job_id)

		if err != nil {
			return nil, utils.ErrorHandler(err, "Update error on backfill visited tracker")
		}

	}
	return output, nil
}

func SeekReopened() ([]string, error) {

	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "db conn error")
	}
	rows, err := db.Query(`
		SELECT job_id FROM JOB_LIFECYCLE WHERE job_state LIKE 'SUSPENDED' ORDER BY last_seen_listed_at ASC limit 1000;
	`)
	if err != nil {
		return nil, utils.ErrorHandler(err, "no no but yes")
	}
	defer rows.Close()

	var output []string

	for rows.Next() {

		var res string

		rows.Scan(&res)
		output = append(output, res)
	}
	return output, nil

}
func GetKeys() ([]string, error) {

	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "db conn error")
	}
	rows, err := db.Query(`
		SELECT * from FILE_KEYS;
	`)
	if err != nil {
		return nil, utils.ErrorHandler(err, "no no but yes")
	}
	defer rows.Close()

	var output []string

	for rows.Next() {

		var res string

		rows.Scan(&res)
		output = append(output, res)
	}
	return output, nil

}

func GetSearchTerms(first_run bool, number_accounts int) ([]SearchTerm, error) {
	db, err := ConnectDb()

	if err != nil {
		return nil, utils.ErrorHandler(err, "db conn error")
	}

	if first_run == true {
		rows, err := db.Query(`
		SELECT search_term_id ,term from SEARCH_TERM LIMIT $1;
	`, number_accounts)
		if err != nil {
			return nil, utils.ErrorHandler(err, "first run failure in GetSearchTerms function")
		}
		defer rows.Close()
		var output []SearchTerm

		for rows.Next() {
			var res SearchTerm

			err = rows.Scan(&res.Search_term_id, &res.Search_term)

			if err != nil {
				return nil, utils.ErrorHandler(err, "Scann load error on first run in GetSearchTerms function")
			}
			_, err := db.Exec(`
			UPDATE SEARCH_TERM SET mid_run = TRUE where search_term_id = $1;
				`, res.Search_term_id)

			if err != nil {
				return nil, utils.ErrorHandler(err, "Update error on first run in GetSearchTerms function")
			}
			output = append(output, res)
		}
		return output, nil
	} else {

		// then we know that its automation that triggered this and number_accounts is actually searchterm id

		//order of ops : 1. update run count by search term id.
		// 2. select search term by id order by min
		// 3. check if min is equal to max then we know we are done else load and return

		_, err := db.Exec(`
		UPDATE SEARCH_TERM SET run_count = 1 , mid_run = FALSE where search_term_id = $1;
	`, number_accounts)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Update error on auto in GetSearchTerms function")
		}

		row := db.QueryRow(`
			SELECT search_term_id , term , min(run_count) , (SELECT max(run_count) FROM SEARCH_TERM WHERE mid_run = FALSE) AS max FROM SEARCH_TERM WHERE mid_run = FALSE GROUP BY search_term_id HAVING run_count=(SELECT min(run_count) FROM SEARCH_TERM)  limit 1;
		`)
		var output []SearchTerm

		var res SearchTerm
		var min int
		var max int

		err = row.Scan(&res.Search_term_id, &res.Search_term, &min, &max)
		if err == sql.ErrNoRows {
			//at n last runs will be True for mid_run but have min count not equal to max which will return no rows
			fmt.Println("job done , waiting for remaining jobs to complete ")
			return nil, nil
		} else if err != nil {
			fmt.Println("something unexpected happened ")
			return nil, err
		}
		if min == max {
			//SHOULD only be reached if all all search terms have had a run . mid_run is false

			fmt.Println(res)
			_, err := db.Exec(`
			UPDATE SEARCH_TERM SET run_count = 0 , total_run_count = total_run_count +1;
			`)
			if err != nil {
				return nil, utils.ErrorHandler(err, "no no but yes")
			}

			return nil, nil
		} else {
			_, err := db.Exec(`
			UPDATE SEARCH_TERM SET mid_run = TRUE where search_term_id = $1
			`, res.Search_term_id)
			if err != nil {
				return nil, utils.ErrorHandler(err, "no no but yes")
			}

		}
		output = append(output, res)

		return output, nil

	}

}
