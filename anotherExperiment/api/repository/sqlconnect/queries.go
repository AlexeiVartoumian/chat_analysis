package sqlconnect

import (
	"api/models"
	"api/utils"
	"database/sql"
	"fmt"
	"log"
	"reflect"

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

func (s *PostgresStore) SeekExpired() ([]string, error) {

	rows, err := s.db.Query(`
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

func (s *PostgresStore) SeekExpiredAuto(filetype string, firstrun bool) ([]string, error) {

	if firstrun == true {
		_, err := s.db.Exec(`
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

	rows, err := s.db.Query(`
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
		// leave it as visited .
		// _, err := s.db.Exec(`
		// UPDATE JOB_LIFECYCLE SET visited = FALSE`)

		// if err != nil {
		// 	return nil, utils.ErrorHandler(err, " last run errored on db output yes")
		// }
		return nil, nil
	}

	for index := range output {

		// _, err := db.Exec(`
		// UPDATE JOB_LIFECYCLE SET visited = TRUE where job_id IN (SELECT job_id FROM JOB_LIFECYCLE where job_state LIKE 'LISTED' and visited = FALSE
		// `)
		job_id := output[index]
		//fmt.Println(job_id)
		_, err := s.db.Exec(`
		 UPDATE JOB_LIFECYCLE SET visited = TRUE where job_id = $1
		 `, job_id)

		if err != nil {
			return nil, utils.ErrorHandler(err, "Update error on backfill visited tracker")
		}

	}
	return output, nil
}

func (s *PostgresStore) SeekReopened() ([]string, error) {

	rows, err := s.db.Query(`
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
func (s *PostgresStore) GetKeys() ([]string, error) {

	rows, err := s.db.Query(`
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

func (s *PostgresStore) GetSearchTerms(first_run bool, number_accounts int, tablename string) ([]SearchTerm, error) {

	if first_run == true {

		query := fmt.Sprintf("SELECT search_term_id , term from %s LIMIT $1", tablename)

		// 	rows, err := db.Query(`
		// 	SELECT search_term_id ,term from SEARCH_TERM LIMIT $1;
		// `, number_accounts)
		rows, err := s.db.Query(query, number_accounts)

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
			query := fmt.Sprintf("UPDATE %s SET mid_run = TRUE where search_term_id = $1;", tablename)

			_, err := s.db.Exec(query, res.Search_term_id)
			// _, err := db.Exec(`
			// UPDATE SEARCH_TERM SET mid_run = TRUE where search_term_id = $1;
			// 	`, res.Search_term_id)

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

		if number_accounts != -1 { // backoff file will send -1 since they already updated table

			query := fmt.Sprintf("UPDATE %s SET run_count = 1 , mid_run = FALSE where search_term_id =$1", tablename)
			_, err := s.db.Exec(query, number_accounts)

			// 	_, err := db.Exec(`
			// 	UPDATE SEARCH_TERM SET run_count = 1 , mid_run = FALSE where search_term_id = $1;
			// `, number_accounts)
			if err != nil {
				return nil, utils.ErrorHandler(err, "Update error on auto in GetSearchTerms function")
			}

		}

		query := fmt.Sprintf("SELECT search_term_id , term , min(run_count) , (SELECT max(run_count) FROM %s WHERE mid_run = FALSE) AS max FROM %s WHERE mid_run = FALSE GROUP BY search_term_id HAVING run_count=(SELECT min(run_count) FROM %s) limit 1;", tablename, tablename, tablename)

		row := s.db.QueryRow(query)
		// row := db.QueryRow(`
		// 	SELECT search_term_id , term , min(run_count) , (SELECT max(run_count) FROM SEARCH_TERM WHERE mid_run = FALSE) AS max FROM SEARCH_TERM WHERE mid_run = FALSE GROUP BY search_term_id HAVING run_count=(SELECT min(run_count) FROM SEARCH_TERM)  limit 1;
		// `)
		var output []SearchTerm

		var res SearchTerm
		var min int
		var max int

		err := row.Scan(&res.Search_term_id, &res.Search_term, &min, &max)
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
			query := fmt.Sprintf("UPDATE %s SET run_count = 0 , total_run_count = total_run_count +1;", tablename)
			_, err := s.db.Exec(query)
			// _, err := db.Exec(`
			// UPDATE SEARCH_TERM SET run_count = 0 , total_run_count = total_run_count +1;
			// `)
			if err != nil {
				return nil, utils.ErrorHandler(err, "no no but yes")
			}

			return nil, nil
		} else {
			query := fmt.Sprintf("UPDATE %s SET mid_run = TRUE where search_term_id = $1", tablename)
			_, err := s.db.Exec(query, res.Search_term_id)
			// _, err := db.Exec(`
			// UPDATE SEARCH_TERM SET mid_run = TRUE where search_term_id = $1
			// `, res.Search_term_id)
			if err != nil {
				return nil, utils.ErrorHandler(err, "no no but yes")
			}

		}
		output = append(output, res)

		return output, nil

	}

}

func (s *PostgresStore) BackoffUpdate(searchtermid int) error {

	// _, err = db.Exec(`
	//     UPDATE SEARCH_TERM SET mid_run = False, run_count = 0 WHERE search_term LIKE $1;
	// `, "%"+searchterm+"%")
	_, err := s.db.Exec(`
        UPDATE SEARCH_TERM SET mid_run = False, run_count = 0 WHERE search_term_id LIKE $1;
    `, searchtermid)
	if err != nil {
		return utils.ErrorHandler(err, "update error")
	}

	return nil
}

func (s *PostgresStore) SeekCompanyChecker() ([]models.COMPANY_DEED, error) {

	rows, err := s.db.Query(`SELECT * FROM COMPANY_DEED WHERE NOT EXISTS (SELECT * FROM COMPANY_METADATA_DEED where COMPANY_DEED.company_id = COMPANY_METADATA_deed.company_id);
`)
	if err != nil {
		return nil, utils.ErrorHandler(err, "yep yep but no")
	}
	defer rows.Close()

	var results []models.COMPANY_DEED

	for rows.Next() {

		var res models.COMPANY_DEED

		rows.Scan(&res.CompanyId, &res.Name, &res.Employer_url)
		results = append(results, res)

	}

	return results, nil
}

func (s *PostgresStore) GetCompanyDeets() ([]models.CompanyDetail, error) {

	//best effort if visited fails for whatever reason then can fall back on set difference
	rows, err := s.db.Query(`
	SELECT company_id, name FROM COMPANY WHERE NOT EXISTS ( SELECT 1 FROM COMPANY_DETAIL WHERE COMPANY_DETAIL.company_id = COMPANY.company_id) AND VISITED = FALSE LIMIT 100;`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "error on CompanyDeets")
	}
	defer rows.Close()

	var output []models.CompanyDetail

	for rows.Next() {
		var res models.CompanyDetail
		err = rows.Scan(&res.Company_id, &res.Company_name)
		_, err := s.db.Exec(
			`UPDATE COMPANY SET VISITED = TRUE where company_id = $1`, res.Company_id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Scann load error on companydetails")
		}

		output = append(output, res)
	}

	fmt.Println(len(output))

	return output, nil
}
func (s *PostgresStore) RedirectIndAuto(firstrun bool) ([]models.JobRedirect_DEED, error) {

	if firstrun == true {
		_, err := s.db.Exec(`
		UPDATE JOBS_DEED SET VISITED = FALSE`)

		if err != nil {
			return nil, utils.ErrorHandler(err, "first run update error")
		}

	}
	// what we say below is give me jobs that point to real page that are still live
	// and the work for it has not been done for it
	rows, err := s.db.Query(`WITH candidate AS ( 
    SELECT job_id , job_url , JOBS_DEED.visited from JOBS_DEED where NOT EXISTS (SELECT REDIRECT_DEED.job_id from REDIRECT_DEED WHERE JOBS_DEED.job_id = REDIRECT_DEED.job_id) LIMIT 100
) SELECT candidate.job_id , candidate.job_url  FROM candidate 
JOIN JOB_LIFECYCLE_DEED on JOB_LIFECYCLE_DEED.job_id = candidate.job_id  where job_state LIKE 'False' and candidate.visited = False; 
`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "query on payload messed up ")
	}
	defer rows.Close()

	var output []models.JobRedirect_DEED

	for rows.Next() {

		var res models.JobRedirect_DEED
		rows.Scan(&res.JobId, &res.JobUrl)

		_, err := s.db.Exec(
			`UPDATE JOBS_DEED SET VISITED = TRUE where job_id = $1`, res.JobId)

		if err != nil {
			fmt.Println(utils.ErrorHandler(err, "bad job"))
			continue
		}
		output = append(output, res)

	}
	fmt.Println(len(output))
	if len(output) == 0 {

		_, err := s.db.Exec(`
		UPDATE JOBS_DEED SET visited = FALSE`)

		if err != nil {
			return nil, utils.ErrorHandler(err, " last run errored on db output yes")
		}
		return nil, nil
	}
	//todo refrence the jobs-deed table to update the visited
	// then work on script
	return output, nil
}

func (s *PostgresStore) RedirectInd() ([]models.JobRedirect_DEED, error) {

	// rows, err := db.Query(`
	// SELECT job_id , job_url FROM jobs_deed limit 100;
	// `)
	rows, err := s.db.Query(`
		SELECT job_id , job_url from JOBS_DEED where NOT EXISTS (SELECT * from REDIRECT_DEED WHERE JOBS_DEED.job_id = REDIRECT_DEED.job_id) LIMIT 100;
	`)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error on upload")
	}

	defer rows.Close()

	var output []models.JobRedirect_DEED
	for rows.Next() {

		var res models.JobRedirect_DEED

		err = rows.Scan(&res.JobId, &res.JobUrl)

		if err != nil {
			return nil, utils.ErrorHandler(err, "Scann load error on redirectlink")
		}

		output = append(output, res)
	}
	return output, nil

}

func (s *PostgresStore) SeekExpiredAutoDeed(firstrun bool) ([]string, error) {

	if firstrun == true {
		_, err := s.db.Exec(`
		UPDATE JOB_LIFECYCLE_DEED SET visited = FALSE`)

		if err != nil {
			return nil, utils.ErrorHandler(err, " first run errored on db output")
		}
	}

	rows, err := s.db.Query(`SELECT job_id FROM JOB_LIFECYCLE_DEED where job_state LIKE 'False' and visited = FALSE ORDER by last_seen_listed_at ASC limit 100;`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "query on payload messed up ")
	}
	defer rows.Close()

	var output []string
	for rows.Next() {

		var res string

		rows.Scan(&res)

		output = append(output, res)
	}

	fmt.Println(len(output))

	if len(output) == 0 {

		_, err := s.db.Exec(
			`UPDATE JOB_LIFECYCLE_DEED SET visited = FALSE`)

		if err != nil {
			return nil, utils.ErrorHandler(err, " last run update error on visited")
		}
	}

	for index := range output {

		job_id := output[index]

		_, err := s.db.Exec(`UPDATE JOB_LIFECYCLE_DEED SET visited = TRUE where job_id = $1`, job_id)

		if err != nil {
			return nil, utils.ErrorHandler(err, "update error on lifecycle visited workflow")
		}
	}

	return output, nil
}

func (s *PostgresStore) RedirectAshLead() ([]models.JobRedirect_DEED, error) {

	rows, err := s.db.Query(`
	SELECT job_id , job_url FROM REDIRECT_DEED where job_url LIKE '%ashby%' and visited =FALSE;
	`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "error on upload")
	}

	defer rows.Close()

	var output []models.JobRedirect_DEED
	for rows.Next() {

		var res models.JobRedirect_DEED

		err = rows.Scan(&res.JobId, &res.JobUrl)

		if err != nil {
			return nil, utils.ErrorHandler(err, "Scann load error on redirectlink")
		}

		output = append(output, res)
	}
	return output, nil

}

// TODO ADD TIMESTAMP ON TABLE FOR LASTSCAN
func (s *PostgresStore) SeekAshCompany() ([]models.AshCompany, error) {

	rows, err := s.db.Query(`
	SELECT location_name , company_url , COMPANY_ASH.company_id , last_scanned_at from COMPANY_ASH , JOBS_ASH where JOBS_ASH.company_id = COMPANY_ASH.company_id order by last_scanned_at ASC limit 75;
	`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "error on upload")
	}

	defer rows.Close()

	var output []models.AshCompany
	for rows.Next() {

		var res models.AshCompany

		err = rows.Scan(&res.LocationName, &res.CompanyUrl, &res.CompanyId, &res.Lastscannedat)

		if err != nil {
			return nil, utils.ErrorHandler(err, "Scann load error on redirectlink")
		}

		output = append(output, res)
	}
	return output, nil
}

func (s *PostgresStore) SeekGreenCompany() ([]models.GreenbyCompany, error) {

	rows, err := s.db.Query(`
		SELECT DISTINCT ON (JOBS_GREEN.company_id) job_id, JOBS_GREEN.company_id, company_name, job_url,
       		COMPANY_GREEN.job_board_public_url, last_scanned_at
			FROM JOBS_GREEN, COMPANY_GREEN
			WHERE COMPANY_GREEN.company_id = JOBS_GREEN.company_id
			ORDER BY JOBS_GREEN.company_id, last_scanned_at ASC
			LIMIT 75;
	`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "error on upload")
	}
	defer rows.Close()

	var output []models.GreenbyCompany

	for rows.Next() {

		var res models.GreenbyCompany

		err = rows.Scan(&res.JobId, &res.CompanyId, &res.CompanyName, &res.JobUrl, &res.CompanyUrl, &res.Lastscannedat)

		if err != nil {
			return nil, utils.ErrorHandler(err, "scaan load error on greenby")
		}
		output = append(output, res)

	}
	return output, nil
}

func (s *PostgresStore) RedirectGreenLead() ([]models.JobRedirect_DEED, error) {

	rows, err := s.db.Query(`
	SELECT job_id , job_url FROM REDIRECT_DEED where job_url LIKE '%greenhouse%' and visited =FALSE;
	`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "error on upload")
	}

	defer rows.Close()

	var output []models.JobRedirect_DEED
	for rows.Next() {

		var res models.JobRedirect_DEED
		res.Origin = "deed"
		err = rows.Scan(&res.JobId, &res.JobUrl)

		if err != nil {
			return nil, utils.ErrorHandler(err, "Scann load error on redirectlink")
		}

		output = append(output, res)
	}
	return output, nil

}
func (s *PostgresStore) RedirectLinkGreen() ([]models.JobRedirect_LinkGreen, error) {

	rows, err := s.db.Query(`
		SELECT rl.job_id, jg.company_apply_url
		FROM redirect_link rl
		JOIN JOB_METADATA jg ON jg.job_id = rl.job_id
		JOIN JOB_LIFECYCLE jl ON jg.job_id = jl.job_id
		WHERE jg.company_apply_url LIKE '%greenhouse%'
		AND jl.job_state = 'LISTED'
		AND rl.status_green = 'pending'
		LIMIT 100;
	`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "error on upload")
	}
	defer rows.Close()

	var output []models.JobRedirect_LinkGreen
	for rows.Next() {

		var res models.JobRedirect_LinkGreen

		err = rows.Scan(&res.JobId, &res.JobUrl)
		res.Origin = "link"
		if err != nil {
			return nil, utils.ErrorHandler(err, "Scann load error on redirectlink")
		}

		output = append(output, res)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.ErrorHandler(err, "rows iteration error")
	}
	rows.Close() // done reading, free the connection before issuing updates

	for _, job := range output {
		_, err := s.db.Exec(`
			UPDATE redirect_link
			SET status_green = 'in_progress'
			WHERE job_id = $1;
		`, job.JobId)

		if err != nil {
			return nil, utils.ErrorHandler(err, "error updating status to in_progress")
		}
	}

	return output, nil
}

func (s *PostgresStore) RedirectLinkAsh() ([]models.JobRedirect_LinkAsh, error) {

	rows, err := s.db.Query(`
		SELECT rl.job_id, jg.company_apply_url
		FROM redirect_link rl
		JOIN JOB_METADATA jg ON jg.job_id = rl.job_id
		JOIN JOB_LIFECYCLE jl ON jg.job_id = jl.job_id
		WHERE jg.company_apply_url LIKE '%ashby%'
		AND jl.job_state = 'LISTED'
		AND rl.status_ash = 'pending'
		LIMIT 100;
	`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "error on upload")
	}
	defer rows.Close()

	var output []models.JobRedirect_LinkAsh
	for rows.Next() {

		var res models.JobRedirect_LinkAsh

		err = rows.Scan(&res.JobId, &res.JobUrl)
		res.Origin = "link"
		if err != nil {
			return nil, utils.ErrorHandler(err, "Scann load error on redirectlink")
		}

		output = append(output, res)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.ErrorHandler(err, "rows iteration error")
	}
	rows.Close()
	for _, job := range output {
		_, err := s.db.Exec(`
			UPDATE redirect_link
			SET status_ash = 'in_progress'
			WHERE job_id = $1;
		`, job.JobId)

		if err != nil {
			return nil, utils.ErrorHandler(err, "error updating status to in_progress")
		}
	}

	return output, nil
}

func (s *PostgresStore) SeekGreenJdChecker() ([]models.JobRedirect_LinkGreen, error) {

	rows, err := s.db.Query(`SELECT job_id , job_url FROM JOBS_GREEN where VISITED = FALSE and not exists (SELECT * FROM JOB_DESCRIPTIONS_GREEN WHERE jobs_GREEN.job_id = JOB_DESCRIPTIONS_GREEN.job_id) limit 75;`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "yep yep but no")
	}
	defer rows.Close()

	var results []models.JobRedirect_LinkGreen

	for rows.Next() {

		var res models.JobRedirect_LinkGreen

		rows.Scan(&res.JobId, &res.JobUrl)

		results = append(results, res)
	}
	rows.Close()

	for _, job := range results {
		_, err := s.db.Exec(`
		UPDATE JOBS_GREEN SET visited = TRUE where job_id = $1;
		`, job.JobId)

		if err != nil {
			return nil, utils.ErrorHandler(err, "update error on db")
		}

	}

	return results, nil
}

func (s *PostgresStore) SeekAshJdChecker() ([]models.JobRedirect_LinkAsh, error) {

	rows, err := s.db.Query(`SELECT job_id , job_url FROM JOBS_ASH where VISITED = FALSE and not exists (SELECT * FROM JOB_DESCRIPTIONS_ASH WHERE JOBS_ASH.job_id = JOB_DESCRIPTIONS_ASH.job_id) limit 75;`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "yep yep but no")
	}
	defer rows.Close()

	var results []models.JobRedirect_LinkAsh

	for rows.Next() {

		var res models.JobRedirect_LinkAsh

		rows.Scan(&res.JobId, &res.JobUrl)

		results = append(results, res)
	}
	rows.Close()

	for _, job := range results {
		_, err := s.db.Exec(`
		UPDATE JOBS_ASH SET visited = TRUE where job_id = $1;
		`, job.JobId)

		if err != nil {
			return nil, utils.ErrorHandler(err, "update error on db")
		}

	}

	return results, nil
}

func (s *PostgresStore) SeekDeedJdChecker() ([]models.JobRedirect_LinkAsh, error) {

	rows, err := s.db.Query(`SELECT job_id  FROM JOBS_DEED where VISITED = FALSE and not exists (SELECT * FROM JOB_DESCRIPTION_DEED WHERE JOBS_DEED.job_id = JOB_DESCRIPTION_DEED.job_id) order by date_advertised desc limit 75;`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "yep yep but no")
	}
	defer rows.Close()

	var results []models.JobRedirect_LinkAsh

	for rows.Next() {

		var res models.JobRedirect_LinkAsh

		rows.Scan(&res.JobId)

		results = append(results, res)
	}
	rows.Close()

	for _, job := range results {
		_, err := s.db.Exec(`
		UPDATE JOBS_DEED SET visited = TRUE where job_id = $1;
		`, job.JobId)

		if err != nil {
			return nil, utils.ErrorHandler(err, "update error on db")
		}

	}

	return results, nil
}

func (s *PostgresStore) SendWorkweek() ([]models.JOB_SEARCH_TERM_WORKWEEK, error) {

	rows, err := s.db.Query(`SELECT JOBS.job_id , company_apply_url from JOBS  JOIN JOB_METADATA on JOBS.job_id = JOB_METADATA.job_id JOIN JOB_LIFECYCLE on JOBS.job_id = JOB_LIFECYCLE.job_id WHERE JOB_LIFECYCLE.job_state LIKE 'LISTED' and company_apply_url LIKE '%workday%' limit 75;`)

	if err != nil {
		return nil, utils.ErrorHandler(err, "yep yep but no")
	}
	defer rows.Close()

	var results []models.JOB_SEARCH_TERM_WORKWEEK

	for rows.Next() {

		var res models.JOB_SEARCH_TERM_WORKWEEK

		rows.Scan(&res.Job_id)

		results = append(results, res)
	}
	rows.Close()
	//TODO THINK OF SOME WAY TO AVOID REPEATED WORK
	// for _, job := range results {
	// 	_, err := s.db.Exec(`
	// 	UPDATE JOBS_DEED SET visited = TRUE where job_id = $1;
	// 	`, job.)

	// 	if err != nil {
	// 		return nil, utils.ErrorHandler(err, "update error on db")
	// 	}

	// }

	return results, nil
}

// SELECT count(*) FROM JOB_metadata, JOB_LIFECYCLE where JOB_METADATA.job_id = JOB_LIFECYCLE.job_id and company_apply_url LIKE '%ashby%' and JOB_LIFECYCLE.job_state LIKE 'LISTED';

// SELECT count(*) FROM JOB_metadata, JOB_LIFECYCLE where JOB_METADATA.job_id = JOB_LIFECYCLE.job_id and company_apply_url LIKE '%greenhouse%' and JOB_LIFECYCLE.job_state LIKE 'LISTED';

func (s *PostgresStore) QueryToJson(query string, args ...interface{}) ([]map[string]interface{}, []ColumnMeta, error) {

	db, err := ConnectDb(1)

	if err != nil {
		return nil, nil, utils.ErrorHandler(err, "db conn error")
	}

	defer db.Close()

	var currentUser string
	db.QueryRow("SELECT current_user").Scan(&currentUser)
	fmt.Println("CURRENT USER:", currentUser)

	var currentDb string
	db.QueryRow("SELECT current_database()").Scan(&currentDb)
	fmt.Println("CURRENT DB:", currentDb)

	fmt.Println("QUERY:", query)
	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println("QUERY ERROR:", err)
		return nil, nil, err
	}

	defer rows.Close()

	//todo dep injection as param and rows.close instead of db.close
	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, err
	}

	scanTypes := make([]reflect.Type, len(colTypes))
	colMeta := make([]ColumnMeta, len(colTypes))

	for i, ct := range colTypes {
		scanTypes[i] = ct.ScanType()
		colMeta[i] = ColumnMeta{
			Name: ct.Name(),
			Type: ct.DatabaseTypeName(),
		}
	}

	var results []map[string]interface{}

	for rows.Next() {

		scanTargets := make([]interface{}, len(cols))

		for i, t := range scanTypes {
			scanTargets[i] = reflect.New(t).Interface()
		}

		if err := rows.Scan(scanTargets...); err != nil {
			return nil, nil, err
		}

		row := make(map[string]interface{})

		for i, col := range cols {
			row[col] = reflect.ValueOf(scanTargets[i]).Elem().Interface()
		}
		results = append(results, row)
	}
	return results, colMeta, rows.Err()
}
