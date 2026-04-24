package cmd

import (
	"cli/models"
	"encoding/json"

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

func parseTimestamp(ts string) (time.Time, error) {
	// Convert "2026-04-13-23-18-27+00-00" → "2026-04-13T23:18:27+00:00"
	if len(ts) < 25 {
		return time.Time{}, fmt.Errorf("unexpected timestamp format: %s", ts)
	}
	iso := ts[:10] + "T" + ts[11:13] + ":" + ts[14:16] + ":" + ts[17:19] + ts[19:22] + ":" + ts[23:]
	return time.Parse(time.RFC3339, iso)
}
func CsvFile(filepath string, tablename string) error {
	// func CsvFile(filepath string, item chan<- models.COMPANY) {
	//jobModel := jobs.model models.JOBS{}
	file, err := os.Open(filepath)
	if err != nil {
		//return ErrorHandler(err, ""), jobModel
		return ErrorHandler(err, "")
	}
	defer file.Close()

	bom := make([]byte, 3)
	file.Read(bom)
	if string(bom) != "\xef\xbb\xbf" {
		file.Seek(0, 0)
	}

	records, err := gocsv.CSVToMaps(file)
	if err != nil {

		return ErrorHandler(err, "uh oh")
	}
	//if tablename == "JOBS" && len(records) > 0 {
	if tablename == "JOBS" && len(records) > 0 {

		Job_And_search_loader(records, tablename, filepath)
		return nil
	}

	if tablename == "JOB_LIFECYCLE" && len(records) > 0 {
		Jobs_LifecycleLoader(records, tablename, filepath)
		return nil
	}

	if tablename == "JOB_LIFECYCLE_UPDATE" && len(records) > 0 {
		Jobs_LifeCycleLiveRolesUpdater(records, filepath)
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

	AddNewRow(search_term, "SEARCH_TERM")
	search_term_id, err := getSearchTermIdHelper(records[0]["search_term"])
	if err != nil {
		fmt.Println("err observed in search term retrieval", ErrorHandler(err, "you brought this on yourself"))
	}
	//workflowid := strings.Split(strings.Split(filepath, "processedJobs-")[1], ".csv")[0]
	meta_data := strings.Split(strings.Split(strings.Split(filepath, "processedJobs-")[1], ".csv")[0], "_")

	workflowid := meta_data[0]
	timestamp, err := parseTimestamp(meta_data[1])

	if err != nil {
		fmt.Println("workflowid extraction or timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
	}
	//InsertTime := time.Now()
	InsertTime := timestamp
	DuplicateCount := 0

	SearchWorkflow := models.SearchWorkflow{
		Workflow_id:      workflowid,
		Search_term_id:   search_term_id,
		Run_at:           InsertTime,
		Total_jobs_found: 0,
		Net_new_found:    0,
	}
	AddNewRow(SearchWorkflow, "SEARCH_WORKFLOW")

	for index, record := range records {
		value, err := JobLoader(record)
		if err != nil {
			fmt.Println("record at index: has not been saved", index, ErrorHandler(err, "you brought this on yourself"))
			continue
		}
		//skipped, _ := AddNewRow(value, tablename)
		skipped, err := AddNewRow(value, tablename)

		if err != nil {
			fmt.Println("Error occured ", ErrorHandler(err, "yep"))
			break
		}

		DuplicateCount += skipped

		JobSearchWorkflow := models.JOB_SEARCH_TERM{
			Job_id:      value.Job_id,
			Workflow_id: workflowid,
		}
		AddNewRow(JobSearchWorkflow, "JOB_SEARCH_TERM")
	}

	UpdateSearchWorkflowCounts(workflowid, len(records), len(records)-DuplicateCount)
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
	// case "JOB_LIFECYCLE":
	// 	return Jobs_LifecycleLoader(record)
	default:
		return nil, nil
	}
}

func Jobs_DescriptionLoader(record map[string]string) (models.JobDescription, error) {

	job_id, err := strconv.Atoi(record["job_id"])

	if err != nil {
		return models.JobDescription{}, ErrorHandler(err, "whoops")
	}

	Job_Description := models.JobDescription{
		JobId:          job_id,
		JobDescription: record["job_description"],
		Encodings:      json.RawMessage(record["encodings"]),
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

	var company_id int

	if record["company_id"] == "N/A" {
		company_id = -1

		record["company"] = "Unknown / individual"
		//return models.COMPANY{}, ErrorHandler(nil, "nil value")
	} else {

		var err error
		company_id, err = strconv.Atoi(record["company_id"])
		if err != nil {
			return models.COMPANY{}, ErrorHandler(err, "uh oh Company id fail parse")

		}
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
		Company_id = -1
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

func Jobs_LifecycleLoader(records []map[string]string, tablename string, filepath string) {

	meta_data := strings.Split(strings.Split(strings.Split(filepath, "job_metadata-")[1], ".csv")[0], "_")

	timestamp, err := parseTimestamp(meta_data[1])

	if err != nil {
		fmt.Println("workflowid extraction or timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
	}

	for index, record := range records {

		value, err := Jobs_Lifecyclemodel(record, timestamp)

		if err != nil {
			fmt.Println("record at index of job metadata for lifecycle: has not been saved", index, ErrorHandler(err, "you brought this on yourself"))
			continue
		}
		AddNewRow(value, "JOB_LIFECYCLE")

	}

}
func Jobs_Lifecyclemodel(record map[string]string, timestamp time.Time) (models.JobLifeCycle, error) {

	job_id, err := strconv.Atoi(record["job_id"])

	if err != nil {
		fmt.Println("timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
		return models.JobLifeCycle{}, ErrorHandler(err, "whoops")
	}
	nextScan := timestamp.AddDate(0, 0, 7)
	Job_lifeCycle := models.JobLifeCycle{
		JobId:            job_id,
		JobState:         record["job_state"],
		FirstSeenAt:      timestamp,
		LastSeenListedAt: timestamp,
		NextScanAt:       &nextScan,
	}

	return Job_lifeCycle, nil

}

func Jobs_LifeCycleLiveRolesUpdater(records []map[string]string, filepath string) error {

	//meta_data := strings.Split(strings.Split(strings.Split(filepath, "job_metadata-")[1], ".csv")[0], "_")
	meta_data := strings.Split(strings.Split(filepath, ".csv")[0], "_")
	timestamp, err := parseTimestamp(meta_data[1])

	db, err := ConnectDb()
	if err != nil {
		fmt.Println("db conn gone wrong", ErrorHandler(err, "you brought this on yourself"))
		return ErrorHandler(err, "whoops")
	}
	defer db.Close()

	if meta_data[0] == "live-roles" {

		for index, record := range records {
			job_id, err := strconv.Atoi(record["job_id"])

			if err != nil {
				fmt.Println("timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
				return ErrorHandler(err, "whoops")
			}

			_, err = db.Exec("UPDATE JOB_LIFECYCLE SET last_seen_listed_at = $1 WHERE job_id = $2", timestamp, job_id)

			if err != nil {
				//http.Error(w, " error updating Student ", http.StatusInternalServerError)
				fmt.Println("record at index ", index, " for live roles as not been saved", ErrorHandler(err, "Db query JobLifecycle update error"))
			}
		}
	} else {
		for index, record := range records {
			job_id, err := strconv.Atoi(record["job_id"])

			job_state := record["job_state"]

			if err != nil {
				fmt.Println("timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
				return ErrorHandler(err, "whoops")
			}

			_, err = db.Exec("UPDATE JOB_LIFECYCLE SET first_seen_closed_at = $1 , job_state = $2 WHERE job_id = $3", timestamp, job_state, job_id)

			if err != nil {
				//http.Error(w, " error updating Student ", http.StatusInternalServerError)
				fmt.Println("record at index ", index, " for live roles as not been saved", ErrorHandler(err, "Db query JobLifecycle update error"))
			}
		}

	}

	return nil
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
        SELECT search_term_id FROM SEARCH_TERM
        WHERE term = $1
    `, searchTerm)
	var search_term_id int
	if err := row.Scan(&search_term_id); err != nil {
		return -1, ErrorHandler(err, "row scan error")
	}
	return search_term_id, nil
}

func UpdateSearchWorkflowCounts(workflowid string, totalJobs int, netNew int) error {
	db, err := ConnectDb()
	if err != nil {
		return ErrorHandler(err, "db conn error")
	}
	defer db.Close()

	_, err = db.Exec(`
        UPDATE SEARCH_WORKFLOW 
        SET total_jobs_found = $1, net_new_jobs = $2
        WHERE workflow_id = $3
    `, totalJobs, netNew, workflowid)

	if err != nil {
		return ErrorHandler(err, "update search workflow counts error")
	}

	return nil
}
