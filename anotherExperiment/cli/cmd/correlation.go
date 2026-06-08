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

	rows, err := db.Query(`
	SELECT cd.company_id, c.company_id FROM COMPANY_DEED cd JOIN COMPANY c ON c.name LIKE '%' || cd.name || '%' OR cd.name LIKE '%' || c.name || '%'
	`)

	if err != nil {
		return ErrorHandler(err, "db scan error")
	}
	defer rows.Close()

	//var correlations []models.CompanyCorrelation

	for rows.Next() {

		var res models.CompanyCorrelation

		rows.Scan(&res.Company_deed_id, &res.Company_link_id)

		AddNewRow(res, "COMPANY_CORRELATION")
	}

	return nil

}
