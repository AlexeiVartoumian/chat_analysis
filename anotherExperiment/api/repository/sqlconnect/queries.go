package sqlconnect

import (
	"api/utils"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

func SearchSimilarJobs(query string) error {

	db, err := ConnectDb()

	if err != nil {
		return utils.ErrorHandler(err, "db conn error")
	}
	defer db.Close()

	embedding, err := GetEmbedding(query)

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
