package sqlconnect

import (
	"api/utils"
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

func SearchSimilarJobs(query string) error {

	db, err := ConnectDb()

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

	db, err := ConnectDb()

	if err != nil {
		return nil, utils.ErrorHandler(err, "db conn error")
	}

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

	db, err := ConnectDb()

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
		SELECT job_id FROM JOB_LIFECYCLE WHERE job_state LIKE 'LISTED'
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

func SeekReopened() ([]string, error) {

	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "db conn error")
	}
	rows, err := db.Query(`
		SELECT job_id FROM JOB_LIFECYCLE WHERE job_state LIKE 'SUSPENDED'
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
