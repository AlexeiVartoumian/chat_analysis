package cmd

import (
	"cli/models"

	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

//*
// need to do a few things here . find the file use csv to parse the contents then parse the date to pass ino loader
// maybe make a model loader function too
// */

func CsvFile(filepath string, tablename string) error {
	// func CsvFile(filepath string, item chan<- models.COMPANY) {
	//jobModel := jobs.model models.JOBS{}
	file, err := os.Open(filepath)
	if err != nil {
		//return ErrorHandler(err, ""), jobModel
		return ErrorHandler(err, "")
	}
	defer file.Close()

	records, err := gocsv.CSVToMaps(file)
	if err != nil {

		return ErrorHandler(err, "uh oh")
	}
	//if tablename == "JOBS" && len(records) > 0 {
	if tablename == "JOBS" && len(records) > 0 {

		Job_And_search_loader(records, tablename, filepath)
		return nil
	}

	for index, record := range records {
		fmt.Println(record)
		fmt.Println(" ")
		value, err := ModelLoader(tablename, record)

		if err != nil {
			//return ErrorHandler(err, "you brought this on yourself")
			//dont return process other records
			fmt.Println("record at index: has not been saved", index, ErrorHandler(err, "you brought this on yourself"))
			continue
		}

		AddNewRow(value, tablename)

	}

	return nil
}

func Job_And_search_loader(records []map[string]string, tablename string, filepath string) {
	search_term := models.Search_term{Search_term: records[0]["search_term"]}

	AddNewRow(search_term, tablename)
	search_term_id, err := getSearchTermIdHelper(records[0]["search_term"])
	workflowid := strings.Split(strings.Split(filepath, "processedJobs-")[1], ".csv")[0]
	if err != nil {
		fmt.Println("err observed in search term retriavel", ErrorHandler(err, "you brought this on yourself"))
	}
	InsertTime := time.Now()
	DuplicateCount := 0
	for index, record := range records {

		value, err := JobLoader(record)
		if err != nil {

			fmt.Println("record at index: has not been saved", index, ErrorHandler(err, "you brought this on yourself"))
			continue
		}
		skipped, _ := AddNewRow(value, tablename)

		DuplicateCount += skipped

		JobSearchWorkflow := models.JOB_SEARCH_TERM{
			Job_id:      value.Job_id,
			Workflow_id: workflowid,
		}
		AddNewRow(JobSearchWorkflow, "JOB_SEARCH_TERM")
	}
	SearchWorkflow := models.SearchWorkflow{
		Workflow_id:      workflowid,
		Search_term_id:   search_term_id,
		Run_at:           InsertTime,
		Total_jobs_found: len(records),
		Net_new_found:    len(records) - DuplicateCount,
	}
	AddNewRow(SearchWorkflow, "SEARCH_WORKFLOW")
}

func ModelLoader(tablename string, record map[string]string) (interface{}, error) {

	switch tablename {
	case "COMPANY":
		return CompanyLoader(record)
	case "COMPANY_METADATA":
		return Company_MetadataLoader(record)
	// case "JOBS":
	// 	return JobLoader(record)
	case "JOB_METADATA":
		return Jobs_MetadataLoader(record)
	case "JOB_DESCRIPTION":
		return Jobs_DescriptionLoader(record)
	default:
		return nil, nil
	}
}

func Jobs_DescriptionLoader(record map[string]string) (models.JOB_DESCRIPTION, error) {

	job_id, err := strconv.Atoi(record["job_id"])

	if err != nil {
		return models.JOB_DESCRIPTION{}, ErrorHandler(err, "whoops")
	}
	Job_Description := models.JOB_DESCRIPTION{
		JobId:          job_id,
		JobDescription: record["job_description"],
	}
	return Job_Description, nil
}

func Jobs_MetadataLoader(record map[string]string) (models.Jobs_metadata, error) {

	job_id, err := strconv.Atoi(record["job_id"])

	if err != nil {
		return models.Jobs_metadata{}, ErrorHandler(err, "uh oh")
	}

	Jobs_metadata := models.Jobs_metadata{
		JobId:           job_id,
		ApplicantsCount: record["applicants_count"],
		CompanyApplyUrl: record["company_apply_url"],
		JobState:        record["job_state"],
	}

	return Jobs_metadata, nil
}

func Company_MetadataLoader(record map[string]string) (models.Company_Metadata, error) {

	employeeCount, err := strconv.Atoi(record["employee_count"])

	if err != nil {
		return models.Company_Metadata{}, ErrorHandler(err, "uh oh error parsing employee count ")
	}

	Company_id, err := strconv.Atoi(record["company_id"])

	if err != nil {
		return models.Company_Metadata{}, ErrorHandler(err, "uh oh error parsing company id")
	}

	Company_metadata := models.Company_Metadata{
		CompanyId:          Company_id,
		Industry:           record["industry"],
		Name:               record["company_name"],
		Description:        record["company_about"],
		EmployeeCount:      employeeCount,
		EmployeeCountRange: record["employee_count_range"],
	}

	return Company_metadata, nil
}

func CompanyLoader(record map[string]string) (models.COMPANY, error) {

	if record["company_id"] == "N/A" {
		return models.COMPANY{}, ErrorHandler(nil, "nil value")
	}
	company_id, err := strconv.Atoi(record["company_id"])
	if err != nil {
		return models.COMPANY{}, ErrorHandler(err, "uh oh Company id fail parse")

	}

	company := models.COMPANY{
		CompanyId: company_id,
		Name:      record["company"],
		Logo:      record["logo"],
	}

	return company, nil
}

func JobLoader(record map[string]string) (models.JOBS, error) {

	//companyid, _ := GetCompanyByIdFromName(record["company"])
	job_id, err := strconv.Atoi(urlHelper(record["job_url"]))
	if err != nil {
		return models.JOBS{}, ErrorHandler(err, "uh oh jobid id fail parse")
	}

	var Company_id int

	_, err2 := strconv.Atoi(record["company_id"])
	if err2 != nil {
		//could be a solo person posting the job

	} else {
		Company_id, err = strconv.Atoi(record["company_id"])
		if err != nil {
			return models.JOBS{}, ErrorHandler(err, "uh oh Company id fail parse")

		}
	}

	easy_apply, err := strconv.ParseBool(record["easy_apply"])

	if err != nil {
		return models.JOBS{}, ErrorHandler(err, "uh oh easy apply fail bool parse")
	}

	promoted, err := strconv.ParseBool(record["promoted"])

	if err != nil {
		return models.JOBS{}, ErrorHandler(err, "uh oh prmoted fail bool parse")
	}

	datePosted := record["posted_date"]
	fmt.Println(datePosted, "value")
	time1, err := time.Parse("2006-01-02", datePosted)

	if err != nil {
		return models.JOBS{}, ErrorHandler(err, "something happened with the date")
	}

	job := models.JOBS{
		Job_id:      job_id,
		Title:       record["title"],
		Location:    record["location"],
		Salary:      record["salary"],
		Date_Posted: time1,
		Job_url:     record["job_url"],
		Search_term: record["search_term"],
		Easy_apply:  easy_apply,
		Promoted:    promoted,
		Expiry_Date: time.Now(),
		Company_id:  Company_id,
	}

	return job, nil

}

func urlHelper(url string) string {

	parts := strings.Split(url, `/`)

	job_id := parts[len(parts)-1]
	return job_id
}

func getSearchTermIdHelper(searchTerm string) (int, error) {

	db, err := ConnectDb()

	if err != nil {
		return -1, ErrorHandler(err, "db conn error")
	}
	defer db.Close()

	row := db.QueryRow(`
		Select search_term_id FROM SEARCH_TERM
		WHERE term = %s
	`, searchTerm)

	var search_term_id int
	if err := row.Scan(&searchTerm); err != nil {
		return -1, ErrorHandler(err, "row scan error")
	}

	return search_term_id, nil
}
