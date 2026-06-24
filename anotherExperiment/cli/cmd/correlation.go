package cmd

import "cli/models"

//TODO add url from other table
func Corellation() error {

	db, err := ConnectDb()

	if err != nil {
		return ErrorHandler(err, "db conn err")
	}
	defer db.Close()

	// queryString := "SELECT cd.name, cd.company_id, c.company_id FROM COMPANY_DEED cd JOIN COMPANY c ON c.name LIKE '%' || cd.name || '%' OR cd.name LIKE '%' || c.name || '%';' "
	// db.Exec(queryString)

	// 	SELECT
	//     cd.company_id as deed_id,
	//     c.company_id as company_id,
	//     cd.name as deed_name,
	//     c.name as company_name,
	//     c.url as company_url  -- already have this
	//     -- cd.url as deed_url  -- just uncomment when it's available
	// FROM COMPANY_DEED cd
	// JOIN COMPANY c ON c.name LIKE '%' || cd.name || '%'
	//    OR cd.name LIKE '%' || c.name || '%'
	// ORDER BY cd.company_id

	//this query sucks
	// 	rows, err := db.Query(`
	// 	SELECT
	// 	c.company_id as company_id,
	//     cd.company_id as deed_id,
	//     c.name as company_name,
	//     cm.company_url as company_url,
	//     cdm.url as deed_url
	// FROM COMPANY_DEED cd
	// JOIN COMPANY c ON c.name LIKE '%' || cd.name || '%'
	//    OR cd.name LIKE '%' || c.name || '%'
	// LEFT JOIN COMPANY_DETAIL cm ON c.company_id = cm.company_id
	// LEFT JOIN COMPANY_METADATA_DEED cdm ON cd.company_id = cdm.company_id
	// ORDER BY cd.company_id;
	// 	`)

	rows, err := db.Query(`
	SELECT
	COMPANY.company_id as company_link_id,
	COMPANY_DEED.company_id as company_deed_id,
    COMPANY.name as company_name,
    COMPANY_DETAIL.company_url as company_url,
    COMPANY_METADATA_DEED.url as deed_url , JOB_METADATA.company_apply_url , JOBS_DEED.job_url
FROM COMPANY_DEED
JOIN COMPANY  ON COMPANY.name = COMPANY_DEED.name
JOIN JOBS on JOBS.company_id = COMPANY.company_id
LEFT JOIN COMPANY_DETAIL  ON COMPANY.company_id = COMPANY_DETAIL.company_id
LEFT JOIN COMPANY_METADATA_DEED ON COMPANY_DEED.company_id = COMPANY_METADATA_DEED.company_id JOIN JOB_METADATA on JOBS.job_id = JOB_METADATA.job_id  join JOBS_DEED on JOBS_DEED.company_id = COMPANY_DEED.company_id;
	`)

	if err != nil {
		return ErrorHandler(err, "db scan error")
	}
	defer rows.Close()

	//var correlations []models.CompanyCorrelation

	for rows.Next() {

		var res models.CompanyCorrelation

		err := rows.Scan(&res.Company_link_id, &res.Company_deed_id, &res.Company_name, &res.Company_link_url, &res.Company_deed_url, &res.Company_apply_url, &res.JobDeed_url)

		if err != nil {
			ErrorHandler(err, "db scan error")
			continue
		}
		AddNewRow(res, "COMPANY_CORRELATION")
	}

	return nil

}

// SELECT
//     COMPANY.name as company_name,
//     COMPANY_DETAIL.company_url as company_url,
//     COMPANY_METADATA_DEED.url as deed_url , JOB_METADATA.company_apply_url , JOBS_DEED.job_url
// FROM COMPANY_DEED
// JOIN COMPANY  ON COMPANY.name = COMPANY_DEED.name
// JOIN JOBS on JOBS.company_id = COMPANY.company_id
// LEFT JOIN COMPANY_DETAIL  ON COMPANY.company_id = COMPANY_DETAIL.company_id
// LEFT JOIN COMPANY_METADATA_DEED ON COMPANY_DEED.company_id = COMPANY_METADATA_DEED.company_id JOIN JOB_METADATA on JOBS.job_id = JOB_METADATA.job_id  join JOBS_DEED on JOBS_DEED.company_id = COMPANY_DEED.company_id ;
