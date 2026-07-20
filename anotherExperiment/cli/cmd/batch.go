package cmd

import (
	"cli/models"
	"database/sql"
	"encoding/json"
	"log"

	"bufio"
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

func extractTimestampStr(filepath string) (string, error) {
	markers := []string{"redirectlinksInd-", "live-roles-deed-", "expired-roles-deed-"}

	var rest string
	found := false
	for _, m := range markers {
		if idx := strings.Index(filepath, m); idx != -1 {
			rest = filepath[idx+len(m):]
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("filepath does not contain a recognised marker: %s", filepath)
	}

	rest = strings.Split(rest, ".csv")[0]
	parts := strings.Split(rest, "_")

	if len(parts) < 2 {
		return "", fmt.Errorf("filepath missing workflowid_timestamp structure: %s", filepath)
	}

	return parts[1], nil
}
func CsvFile(filepath string, tablename string) error {
	// func CsvFile(filepath string, item chan<- models.COMPANY) {
	//jobModel := jobs.model models.JOBS{}

	if tablename == "FILE_KEYS" {

		InsertNewKeys(filepath, tablename)
		return nil
	}

	if tablename == "FILE_KEYS_DEED" {
		InsertNewKeys(filepath, tablename)
		return nil
	}

	if tablename == "FILE_KEYS_ASH" {
		InsertNewKeys(filepath, tablename)
		return nil
	}

	if tablename == "FILE_KEYS_GREEN" {
		InsertNewKeys(filepath, tablename)
		return nil
	}

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

	if tablename == "JOBS_DEED" && len(records) > 0 {

		Job_And_search_loader(records, tablename, filepath)
		return nil
	}

	if tablename == "JOBS_ASH" && len(records) > 0 {
		Job_and_search_loader_ash(records, tablename, filepath)
		return nil
	}

	if tablename == "JOBS_GREEN" && len(records) > 0 {
		Job_and_search_loader_green(records, tablename, filepath)
		return nil
	}

	if tablename == "JOB_LIFECYCLE" && len(records) > 0 {
		Jobs_LifecycleLoader(records, tablename, filepath)
		return nil
	}

	if tablename == "JOB_LIFECYCLE_DEED" && len(records) > 0 {
		Jobs_LifecycleDeedLoader(records, tablename, filepath)
		return nil
	}
	if tablename == "JOB_LIFECYCLE_ASH" && len(records) > 0 {
		Jobs_LifecycleAshLoader(records, tablename, filepath)
		return nil
	}
	if tablename == "JOB_LIFECYCLE_GREEN" && len(records) > 0 {
		Jobs_LifecycleGreenLoader(records, tablename, filepath)
		return nil
	}
	if tablename == "JOB_LIFECYCLE_UPDATE" && len(records) > 0 {
		Jobs_LifeCycleLiveRolesUpdater(records, filepath)
		return nil
	}
	if tablename == "TEAM_ASH" && len(records) > 0 {
		Team_loader_ash(records, tablename, filepath)
		return nil
	}
	fmt.Println("records length:", len(records))
	if tablename == "JOB_LIFECYCLE_UPDATE_SUSPENDED" && len(records) > 0 {
		Jobs_LifeCycleSuspendedRolesUpdater(records, filepath)
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

func InsertNewKeys(filepath string, tablename string) error {

	file, err := os.Open(filepath)
	if err != nil {
		//return ErrorHandler(err, ""), jobModel
		return ErrorHandler(err, "filepat open issue")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		key := models.File_key{
			File_name: scanner.Text(),
		}

		AddNewRow(key, tablename)
	}

	return nil

}

func Job_And_search_loader(records []map[string]string, tablename string, filepath string) {
	search_term := models.Search_term{Search_term: records[0]["search_term"]}
	if tablename != "JOBS" {
		AddNewRow(search_term, "SEARCH_TERM_DEED")
	} else {
		AddNewRow(search_term, "SEARCH_TERM")
	}
	search_term_id, err := getSearchTermIdHelper(records[0]["search_term"], tablename)
	if err != nil {
		fmt.Println("err observed in search term retrieval", ErrorHandler(err, "you brought this on yourself"))
	}

	var meta_data []string

	if tablename != "JOBS" {
		meta_data = strings.Split(strings.Split(strings.Split(filepath, "processedJobsInd-")[1], ".csv")[0], "_")

	} else {
		//workflowid := strings.Split(strings.Split(filepath, "processedJobs-")[1], ".csv")[0]
		meta_data = strings.Split(strings.Split(strings.Split(filepath, "processedJobs-")[1], ".csv")[0], "_")
	}
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
	if tablename != "JOBS" {
		_, err := AddNewRow(SearchWorkflow, "SEARCH_WORKFLOW_DEED")
		if err != nil {
			fmt.Println("workflow insert failed:", err)
			return
		}
	} else {
		AddNewRow(SearchWorkflow, "SEARCH_WORKFLOW")

	}

	if tablename != "JOBS" {
		for index, record := range records {

			value, err := JobLoaderDeed(record)
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

			JobSearchWorkflow := models.JOB_SEARCH_TERM_DEED{
				Job_id:      value.Job_id,
				Workflow_id: workflowid,
				Is_new_job:  skipped == 0,
			}
			AddNewRow(JobSearchWorkflow, "JOB_SEARCH_TERM_DEED")
		}
	} else {
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
				Is_new_job:  skipped == 0,
			}
			AddNewRow(JobSearchWorkflow, "JOB_SEARCH_TERM")
		}
	}

	// if tablename != "JOBS" {
	// 	DuplicateCount = AddJobToDb(records, tablename, DuplicateCount, workflowid, JobLoaderDeed)
	// }

	//UpdateSearchWorkflowCounts(workflowid, len(records), len(records)-DuplicateCount, tablename)
	if err := UpdateSearchWorkflowCounts(workflowid, len(records), len(records)-DuplicateCount, tablename); err != nil {
		fmt.Println("failed to update workflow counts:", err)
	}
}
func Job_and_search_loader_ash(records []map[string]string, tablename string, filepath string) {

	db, err := ConnectDb() //awful
	var meta_data []string
	if strings.HasPrefix(filepath, "processedJobsAsh") {
		meta_data = strings.Split(strings.Split(strings.Split(filepath, "processedJobsAsh-")[1], ".csv")[0], "_")
	} else {
		meta_data = strings.Split(strings.Split(strings.Split(filepath, "AshJobsByCompany-")[1], ".csv")[0], "_")
	}
	workflowid := meta_data[0]
	timestamp, err := parseTimestamp(meta_data[1])
	if err != nil {
		fmt.Println("workflowid extraction or timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
	}
	//InsertTime := time.Now()

	InsertTime := timestamp
	DuplicateCount := 0
	SearchWorkflow := models.SearchWorkflowAsh{
		Workflow_id:      workflowid,
		Run_at:           InsertTime,
		Total_jobs_found: 0,
		Net_new_found:    0,
	}
	_, err = AddNewRow(SearchWorkflow, "SEARCH_WORKFLOW_ASH")
	if err != nil {
		fmt.Println("workflow insert failed:", err)
		return
	}

	for index, record := range records {
		value, err := JobLoaderAsh(record)
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

		// job row is persisted at this point — safe to mark redirect status done
		if value.OriginLinkID != nil {
			if _, err := db.Exec(
				`UPDATE redirect_link SET status_ash = 'done' WHERE job_id = $1 AND status_ash = 'in_progress'`,
				*value.OriginLinkID,
			); err != nil {
				fmt.Println("failed to update redirect_link for job_id", *value.OriginLinkID, ErrorHandler(err, "redirect update failed"))
			}
		} else if value.Origin_deed_id != nil {
			if _, err := db.Exec(
				`UPDATE redirect_deed SET status_ash = 'done' WHERE job_id = $1 AND status_ash = 'in_progress'`,
				*value.Origin_deed_id,
			); err != nil {
				fmt.Println("failed to update redirect_deed for job_id", *value.Origin_deed_id, ErrorHandler(err, "redirect update failed"))
			}
		}

		JobSearchWorkflow := models.JOB_SEARCH_TERM_ASH{
			Job_id:      value.JobID,
			Workflow_id: workflowid,
			Is_new_job:  skipped == 0,
		}
		_, err = AddNewRow(JobSearchWorkflow, "JOB_SEARCH_TERM_ASH")
		if err != nil {
			fmt.Println("job search term insert failed for job_id", value.JobID, ErrorHandler(err, "yep"))
		}
	}

	UpdateSearchWorkflowCounts(workflowid, len(records), len(records)-DuplicateCount, tablename)
}
func Job_and_search_loader_green(records []map[string]string, tablename string, filepath string) {

	db, err := ConnectDb() //awful
	var meta_data []string
	if strings.HasPrefix(filepath, "processedJobsGreen") {
		meta_data = strings.Split(strings.Split(strings.Split(filepath, "processedJobsGreen-")[1], ".csv")[0], "_")
	} // else {
	// 	meta_data = strings.Split(strings.Split(strings.Split(filepath, "AshJobsByCompany-")[1], ".csv")[0], "_")
	// }

	workflowid := meta_data[0]
	timestamp, err := parseTimestamp(meta_data[1])
	if err != nil {
		fmt.Println("workflowid extraction or timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
	}
	InsertTime := timestamp
	DuplicateCount := 0
	SearchWorkflow := models.SearchWorkflowGreen{
		Workflow_id:      workflowid,
		Run_at:           InsertTime,
		Total_jobs_found: 0,
		Net_new_found:    0,
	}
	_, err = AddNewRow(SearchWorkflow, "SEARCH_WORKFLOW_GREEN")
	if err != nil {
		fmt.Println("workflow insert failed:", err)
		return
	}

	for index, record := range records {
		value, err := JobLoaderGreen(record)
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
		JobSearchWorkflow := models.JOB_SEARCH_TERM_GREEN{
			Job_id:      value.JobID,
			Workflow_id: workflowid,
			Is_new_job:  skipped == 0,
		}
		_, err = AddNewRow(JobSearchWorkflow, "JOB_SEARCH_TERM_GREEN")

		if err == nil && value.OriginLinkID != nil {

			_, err = db.Exec("UPDATE REDIRECT_DEED SET visited = TRUE WHERE job_id = $1", value.OriginLinkID)

		}
		UpdateSearchWorkflowCounts(workflowid, len(records), len(records)-DuplicateCount, tablename)

	}
}

func ModelLoader(tablename string, record map[string]string) (interface{}, error) {

	switch tablename {
	case "COMPANY":
		return CompanyLoader(record)
	case "COMPANY_DEED":
		return CompanyDeedLoader(record)
	case "COMPANY_GREEN":
		return Company_GreenLoader(record)
	case "COMPANY_ASH":
		return Company_AshLoader(record)

	case "COMPANY_METADATA":
		return Company_MetadataLoader(record)
	case "COMPANY_METADATA_DEED":
		return Company_MetadataDeedLoader(record)

	// case "JOBS_ASH":
	// 	return JobLoaderAsh(record)
	// case "JOBS":
	// 	return JobLoader(record)
	case "JOB_METADATA":
		return Jobs_MetadataLoader(record)
	case "JOB_DESCRIPTION":
		return Jobs_DescriptionLoader(record)
	case "JOB_DESCRIPTION_DEED":
		return Jobs_DescriptionDeedLoader(record)

	case "JOB_DESCRIPTIONS_ASH":
		return Jobs_DescriptionAshLoader(record)

	case "JOB_DESCRIPTIONS_GREEN":
		return Jobs_DescriptionGreenLoader(record)

	case "COMPANY_DETAIL":
		return CompanyDetailLoader(record)

	case "REDIRECT_DEED":
		return RedirectDeedLoader(record)

	case "DEPARTMENT_ASH":
		return DepartmentAshLoader(record)
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

func DepartmentAshLoader(record map[string]string) (models.COMPANY_ASH_DEPARTMENT, error) {

	AshDepartment := models.COMPANY_ASH_DEPARTMENT{
		DepartmentId:    record["departmentId"],
		Department_name: record["DepartmentName"],
		CompanyId:       record["organizationId"],
	}
	return AshDepartment, nil

}

func TeamAshLoader(record map[string]string) (models.COMPANY_ASH_TEAM, error) {

	Ashteam := models.COMPANY_ASH_TEAM{
		TeamId:           record["teamId"],
		TeamName:         record["teamName"],
		DepartmentId:     record["departmentId"],
		ParentTeamId:     NullableString(record["parentTeamId"]),
		TeamExternalName: NullableString(record["teamExternalName"]),
	}
	return Ashteam, nil
}
func Jobs_DescriptionDeedLoader(record map[string]string) (models.JobDescription_DEED, error) {

	job_id := record["job_id"]

	Job_Description := models.JobDescription_DEED{
		JobId:          job_id,
		JobDescription: record["job_description"],
	}
	return Job_Description, nil
}
func Jobs_DescriptionAshLoader(record map[string]string) (models.JobDescription_Ash, error) {

	job_id := record["job_id"]

	Job_Description := models.JobDescription_Ash{
		JobId:          job_id,
		JobDescription: record["descriptionHtml"],
	}
	return Job_Description, nil
}

func Jobs_DescriptionGreenLoader(record map[string]string) (models.JobDescription_Green, error) {

	job_id, err := strconv.Atoi(record["job_id"])

	if err != nil {
		return models.JobDescription_Green{}, ErrorHandler(err, "whoops")
	}

	Job_Description := models.JobDescription_Green{
		JobId:                    job_id,
		JobDescription:           record["job_description"],
		JobDescriptionConclusion: NullableString(record["job_description_conclusion"]),
		Company_about:            NullableString(record["company_about"]),
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

func NullableString(s string) *string {

	if s == "" {
		return nil
	}

	return &s
}

func Company_MetadataDeedLoader(record map[string]string) (models.CompanyDeed_Metadata, error) {

	Company_metadata := models.CompanyDeed_Metadata{
		CompanyId:            record["company_id"],
		Name:                 record["company_name"],
		Employee_count_range: NullableString(record["company_size"]),
		Industry:             NullableString(record["industry"]),
		Revenue:              NullableString(record["revenue"]),
		Description:          NullableString(record["description"]),
		Url:                  NullableString(record["company_url"]),
	}

	return Company_metadata, nil
}

func Company_AshLoader(record map[string]string) (models.COMPANY_ASH, error) {

	for k, v := range record {
		fmt.Printf("key=%q value=%q\n", k, v)
	}
	Company_Ash := models.COMPANY_ASH{
		CompanyId:         record["organizationId"],
		Name:              record["companyName"],
		Public_website:    NullableString(record["publicWebsite"]),
		CustomJobsPageUrl: NullableString(record["customJobsPageUrl"]),
		Timezone:          NullableString(record["timezone"]),
		Company_url:       NullableString(record["company_url"]),
	}

	return Company_Ash, nil
}

func Company_GreenLoader(record map[string]string) (models.COMPANY_GREEN, error) {

	company_id, err := strconv.Atoi(record["company_id"])

	if err != nil {
		fmt.Println("timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
		return models.COMPANY_GREEN{}, ErrorHandler(err, "whoops")
	}
	Company_Green := models.COMPANY_GREEN{
		CompanyId:         company_id,
		CompanyName:       record["company_name"],
		JobBoardPublicUrl: NullableString(record["job_board_public_url"]),
		DomainUrl:         NullableString(record["domainurl"]),
		Company_url:       NullableString(record["company_url"]),
		CompanyAbout:      NullableString(record["company_about"]),
	}

	return Company_Green, nil
}

func toJSONB(s string) *json.RawMessage {
	if s == "" {
		return nil
	}
	if !json.Valid([]byte(s)) {
		log.Printf("invalid JSON, skipping: %q", s)
		return nil
	}
	raw := json.RawMessage(s)
	return &raw
}
func CompanyDetailLoader(record map[string]string) (models.CompanyDetail, error) {

	staffCount, err := strconv.Atoi(record["staff_count"])
	if err != nil {
		return models.CompanyDetail{}, ErrorHandler(err, "uh oh error parsing employee count ")
	}

	companyId, err := strconv.Atoi(record["company_id"])
	if err != nil {
		return models.CompanyDetail{}, ErrorHandler(err, "uh oh error parsing company id")
	}

	ms, err := strconv.ParseInt(record["created_at"], 10, 64)
	if err != nil {
		return models.CompanyDetail{}, ErrorHandler(err, "uh oh error parsing created_at")
	}
	createdAt := time.Unix(ms/1000, 0).UTC()

	companyDetail := models.CompanyDetail{
		CompanyId:   companyId,
		Name:        record["company_name"],
		CompanySlug: record["company_slug"],
		CompanyUrl:  record["company_url"],
		// Specialties:         json.RawMessage(record["specialties"]),
		// Locations:           json.RawMessage(record["locations"]),
		Specialties:         toJSONB(record["specialties"]),
		Locations:           toJSONB(record["locations"]),
		ExtendedDescription: record["extended_description"],
		StaffCount:          staffCount,
		HeadquarterCity:     record["headquarter_city"],
		Created_at:          createdAt,
	}

	return companyDetail, nil
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

func CompanyDeedLoader(record map[string]string) (models.COMPANY_DEED, error) {

	company_id := record["company_id"]

	company := models.COMPANY_DEED{
		CompanyId:    company_id,
		Name:         record["company_name"],
		Employer_url: record["employer_url"],
	}
	return company, nil
}

func RedirectDeedLoader(record map[string]string) (models.RedirectDeed, error) {

	redirect := models.RedirectDeed{
		JobId:   record["job_id"],
		Job_url: record["job_url"],
	}
	return redirect, nil
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

func JobLoaderDeed(record map[string]string) (models.JOBS_DEED, error) {

	//companyid, _ := GetCompanyByIdFromName(record["company"])
	job_id := record["job_id"]

	Company_id := record["company_id"]

	var organic_apply int

	_, err3 := strconv.Atoi(record["organic_apply"])
	if err3 != nil {
		//could be a solo person posting the job
		organic_apply = -1
	} else {
		organic_apply, err3 = strconv.Atoi(record["organic_apply"])
		if err3 != nil {
			return models.JOBS_DEED{}, ErrorHandler(err3, "uh oh Company id fail parse")
		}
	}

	is_repost, err := strconv.ParseBool(record["is_repost"])

	if err != nil {
		return models.JOBS_DEED{}, ErrorHandler(err, "uh oh is repost fail bool parse")
	}

	is_latest, err := strconv.ParseBool(record["is_latest"])

	if err != nil {
		return models.JOBS_DEED{}, ErrorHandler(err, "uh oh is latest fail bool parse")
	}

	date_pub, err := strconv.ParseInt(record["date_published"], 10, 64)
	if err != nil {
		return models.JOBS_DEED{}, ErrorHandler(err, "uh oh date published fail bool parse")
	}
	date_published := time.Unix(date_pub/1000, 0)

	date_ad, err := strconv.ParseInt(record["date_advertised"], 10, 64)
	if err != nil {
		return models.JOBS_DEED{}, ErrorHandler(err, "uh oh date published fail bool parse")
	}
	date_advertised := time.Unix(date_ad/1000, 0)

	job := models.JOBS_DEED{
		Job_id:          job_id,
		Title:           record["title"],
		Date_published:  date_published,
		Date_advertised: date_advertised,
		Job_url:         record["job_url"],
		Search_term:     record["search_term"],
		Organic_apply:   organic_apply,
		Is_repost:       is_repost,
		Is_latest:       is_latest,
		Company_id:      Company_id,
	}

	return job, nil

}

// add boolean to determine what to pass
func JobLoaderAsh(record map[string]string) (models.JobAsh, error) {

	is_listed, err := strconv.ParseBool(record["isListed"])
	if err != nil {
		return models.JobAsh{}, ErrorHandler(err, "uh oh is_listed fail bool parse")
	}

	// "2026-04-16" -> YYYY-MM-DD
	//datePub, err := time.Parse("2006-01-02", record["publishedDate"])

	// if err != nil {
	// 	return models.JobAsh{}, ErrorHandler(err, "uh oh published_date fail time parse")
	// }

	// "2026-05-19T15:40:48.707Z" -> RFC3339 with milliseconds
	//updatedAt, err := time.Parse(time.RFC3339Nano, record["updatedAt"])
	// if err != nil {
	// 	return models.JobAsh{}, ErrorHandler(err, "uh oh updated_at fail time parse")
	// }
	var originLinkID *int64
	var originDeedID *string

	origin := record["origin"]

	switch origin {
	case "link":
		raw := strings.TrimSpace(record["origin_link_id"])
		if raw != "" {
			v, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				log.Printf("invalid origin_link_id %q: %v", raw, err)
			} else {
				originLinkID = &v
			}
		}
	default:
		raw := strings.TrimSpace(record["origin_deed_id"])
		if raw != "" {
			originDeedID = &raw
		}
	}
	job := models.JobAsh{
		JobID:    record["job_id"],
		Title:    record["title"],
		IsListed: &is_listed,
		// PublishedDate:  datePub,
		// UpdatedAt:      updatedAt,
		JobURL:         record["job_url"],
		TeamName:       NullableString(record["teamName"]),
		DepartmentName: NullableString(record["DepartmentName"]),
		LocationName:   NullableString(record["locationName"]),
		EmploymentType: NullableString(record["employmentType"]),
		WorkplaceType:  NullableString(record["workplaceType"]),
		Salary:         NullableString(record["salary"]),
		CompanyID:      record["organizationId"],
		OriginLinkID:   originLinkID,
		Origin_deed_id: originDeedID,
	}
	return job, nil
}

func JobLoaderGreen(record map[string]string) (models.JobGreen, error) {

	job_id, err := strconv.Atoi(record["job_id"])

	if err != nil {
		fmt.Println("jobid conversion went wrong", ErrorHandler(err, "you brought this on yourself"))
		return models.JobGreen{}, ErrorHandler(err, "whoops")
	}

	company_id, err := strconv.Atoi(record["company_id"])

	if err != nil {
		fmt.Println("companyid extraction went wrong", ErrorHandler(err, "you brought this on yourself"))
		return models.JobGreen{}, ErrorHandler(err, "whoops")
	}

	// "2026-05-19T15:40:48.707Z" -> RFC3339 with milliseconds
	datePub, err := time.Parse(time.RFC3339Nano, record["publishedDate"])
	if err != nil {
		return models.JobGreen{}, ErrorHandler(err, "uh oh updated_at fail time parse")
	}

	job := models.JobGreen{
		JobID:         job_id,
		Title:         record["title"],
		LocationName:  NullableString(record["locationName"]),
		CompanyID:     company_id,
		Company_name:  record["company_name"],
		JobURL:        record["job_url"],
		RedirectURL:   NullableString(record["redirect_url"]),
		PublishedDate: datePub,
		Salary:        NullableString(record["salary"]),
		OriginLinkID:  NullableString(record["origin_link_id"]),
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

func Jobs_LifecycleAshLoader(records []map[string]string, tablename string, filepath string) {

	if strings.HasPrefix(filepath, "processedJobsAsh") {
		meta_data := strings.Split(strings.Split(strings.Split(filepath, "processedJobsAsh-")[1], ".csv")[0], "_")

		timestamp, err := parseTimestamp(meta_data[1])

		if err != nil {
			fmt.Println("workflowid extraction or timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
		}

		for index, record := range records {

			value, err := Jobs_LifecycleAshmodel(record, timestamp)
			if err != nil {
				fmt.Println("record at index of job metadata for lifecycle: has not been saved", index, ErrorHandler(err, "you brought this on yourself"))
				continue
			}
			AddNewRow(value, "JOB_LIFECYCLE_ASH")

		}
	} else if strings.HasPrefix(filepath, "AshJobsByCompany") {
		meta_data := strings.Split(strings.Split(strings.Split(filepath, "AshJobsByCompany-")[1], ".csv")[0], "_")

		timestamp, err := parseTimestamp(meta_data[1])

		if err != nil {
			fmt.Println("workflowid extraction or timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
		}
		for index, record := range records {

			value, err := Jobs_LifecycleAshmodel(record, timestamp)
			if err != nil {
				fmt.Println("record at index of job metadata for lifecycle: has not been saved", index, ErrorHandler(err, "you brought this on yourself"))
				continue
			}
			AddNewRow(value, "JOB_LIFECYCLE_ASH")

		}
	} else {
		fmt.Println("please implement")

		//meta_data := strings.Split(strings.Split(strings.Split(filepath, "redirectlinksInd-")[1], ".csv")[0], "_")
		// timestampStr, err := extractTimestampStr(filepath)
		// if err != nil {
		// 	fmt.Println("could not extract timestamp from filepath", ErrorHandler(err, "you brought this on yourself"))
		// 	return // or continue
		// }
		// timestamp, err := parseTimestamp(timestampStr)
		// if err != nil {
		// 	fmt.Println("workflowid extraction or timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
		// }
		// db, err := ConnectDb()
		// if err != nil {
		// 	fmt.Println("db conn gone wrong", ErrorHandler(err, "you brought this on yourself"))
		// }
		// defer db.Close()

		// for index, record := range records {

		// 	if err != nil {
		// 		fmt.Println("record at index of job metadata for lifecycle: has not been saved", index, ErrorHandler(err, "you brought this on yourself"))
		// 		continue
		// 	}
		// 	fmt.Println("this is job state", record["job_state"])
		// 	if strings.EqualFold(strings.TrimSpace(record["job_state"]), "true") {

		// 		_, err = db.Exec("UPDATE JOB_LIFECYCLE_DEED SET first_seen_closed_at = $1, job_state = $2 WHERE job_id = $3", timestamp, record["job_state"], record["job_id"])

		// 		if err != nil {
		// 			//http.Error(w, " error updating Student ", http.StatusInternalServerError)
		// 			fmt.Println("record at index ", index, " for expired job_lifecycle not saved", ErrorHandler(err, "Db query JobLifecycle update error"))
		// 		}

		// 	} else {

		// 		_, err = db.Exec("UPDATE JOB_LIFECYCLE_DEED SET last_seen_listed_at = $1 WHERE job_id = $2", timestamp, record["job_id"])

		// 		if err != nil {
		// 			//http.Error(w, " error updating Student ", http.StatusInternalServerError)
		// 			fmt.Println("record at index ", index, " for still live job_lifecycle not saved", ErrorHandler(err, "Db query JobLifecycle update error"))
		// 		}
		// 	}
		// }

	}

}
func Jobs_LifecycleGreenLoader(records []map[string]string, tablename string, filepath string) {
	if strings.HasPrefix(filepath, "processedJobsGreen") {
		meta_data := strings.Split(strings.Split(strings.Split(filepath, "processedJobsGreen-")[1], ".csv")[0], "_")

		timestamp, err := parseTimestamp(meta_data[1])

		if err != nil {
			fmt.Println("workflowid extraction or timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
		}
		//TODO HANDLE DUPALICATES TO UPDATE LAST SEEN AND ALSO JOBSTATE
		for index, record := range records {

			value, err := Jobs_LifecycleGreenmodel(record, timestamp)
			if err != nil {
				fmt.Println("record at index of job metadata for lifecycle: has not been saved", index, ErrorHandler(err, "you brought this on yourself"))
				continue
			}
			AddNewRow(value, "JOB_LIFECYCLE_GREEN")

		}
	} else {
		fmt.Println("please implement")
	}
	// } else if strings.HasPrefix(filepath, "AshJobsByCompany") {
	// 	meta_data := strings.Split(strings.Split(strings.Split(filepath, "AshJobsByCompany-")[1], ".csv")[0], "_")

	// 	timestamp, err := parseTimestamp(meta_data[1])

	// 	if err != nil {
	// 		fmt.Println("workflowid extraction or timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
	// 	}
	// 	for index, record := range records {

	// 		value, err := Jobs_LifecycleAshmodel(record, timestamp)
	// 		if err != nil {
	// 			fmt.Println("record at index of job metadata for lifecycle: has not been saved", index, ErrorHandler(err, "you brought this on yourself"))
	// 			continue
	// 		}
	// 		AddNewRow(value, "JOB_LIFECYCLE_ASH")

	// 	}

}

func Jobs_LifecycleDeedLoader(records []map[string]string, tablename string, filepath string) {

	if strings.HasPrefix(filepath, "processedJobsInd") {
		meta_data := strings.Split(strings.Split(strings.Split(filepath, "processedJobsInd-")[1], ".csv")[0], "_")

		timestamp, err := parseTimestamp(meta_data[1])

		if err != nil {
			fmt.Println("workflowid extraction or timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
		}

		for index, record := range records {

			value, err := Jobs_LifecycleDeedmodel(record, timestamp)

			if err != nil {
				fmt.Println("record at index of job metadata for lifecycle: has not been saved", index, ErrorHandler(err, "you brought this on yourself"))
				continue
			}
			AddNewRow(value, "JOB_LIFECYCLE_DEED")

		}
	} else {
		//meta_data := strings.Split(strings.Split(strings.Split(filepath, "redirectlinksInd-")[1], ".csv")[0], "_")
		timestampStr, err := extractTimestampStr(filepath)
		if err != nil {
			fmt.Println("could not extract timestamp from filepath", ErrorHandler(err, "you brought this on yourself"))
			return // or continue
		}
		timestamp, err := parseTimestamp(timestampStr)
		if err != nil {
			fmt.Println("workflowid extraction or timestamp extraction wrong", ErrorHandler(err, "you brought this on yourself"))
		}
		db, err := ConnectDb()
		if err != nil {
			fmt.Println("db conn gone wrong", ErrorHandler(err, "you brought this on yourself"))
		}
		defer db.Close()

		for index, record := range records {

			if err != nil {
				fmt.Println("record at index of job metadata for lifecycle: has not been saved", index, ErrorHandler(err, "you brought this on yourself"))
				continue
			}
			fmt.Println("this is job state", record["job_state"])
			if strings.EqualFold(strings.TrimSpace(record["job_state"]), "true") {

				_, err = db.Exec("UPDATE JOB_LIFECYCLE_DEED SET first_seen_closed_at = $1, job_state = $2 WHERE job_id = $3", timestamp, record["job_state"], record["job_id"])

				if err != nil {
					//http.Error(w, " error updating Student ", http.StatusInternalServerError)
					fmt.Println("record at index ", index, " for expired job_lifecycle not saved", ErrorHandler(err, "Db query JobLifecycle update error"))
				}

			} else {

				_, err = db.Exec("UPDATE JOB_LIFECYCLE_DEED SET last_seen_listed_at = $1 WHERE job_id = $2", timestamp, record["job_id"])

				if err != nil {
					//http.Error(w, " error updating Student ", http.StatusInternalServerError)
					fmt.Println("record at index ", index, " for still live job_lifecycle not saved", ErrorHandler(err, "Db query JobLifecycle update error"))
				}
			}
		}

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
		SuspendedCount:   0,
	}

	return Job_lifeCycle, nil

}

func Jobs_LifecycleDeedmodel(record map[string]string, timestamp time.Time) (models.JobLifeCycleDeed, error) {

	nextScan := timestamp.AddDate(0, 0, 7)
	Job_lifeCycle := models.JobLifeCycleDeed{
		JobId:            record["job_id"],
		JobState:         record["job_state"],
		FirstSeenAt:      timestamp,
		LastSeenListedAt: timestamp,
		NextScanAt:       &nextScan,
	}

	return Job_lifeCycle, nil

}
func Jobs_LifecycleAshmodel(record map[string]string, timestamp time.Time) (models.JobLifeCycleAsh, error) {

	nextScan := timestamp.AddDate(0, 0, 7)
	Job_lifeCycle := models.JobLifeCycleAsh{
		JobId:            record["job_id"],
		Job_state:        record["isListed"],
		FirstSeenAt:      timestamp,
		LastSeenListedAt: timestamp,
		NextScanAt:       &nextScan,
	}

	return Job_lifeCycle, nil

}

func Jobs_LifecycleGreenmodel(record map[string]string, timestamp time.Time) (models.JobLifeCycleGreen, error) {

	nextScan := timestamp.AddDate(0, 0, 7)
	job_id, err := strconv.Atoi(record["job_id"])
	if err != nil {
		fmt.Println("something wrong", ErrorHandler(err, "you brought this on yourself"))
		return models.JobLifeCycleGreen{}, ErrorHandler(err, "whoops")
	}

	//TODO  jobstate is default true will have to update on consequrn riuns
	Job_lifeCycle := models.JobLifeCycleGreen{
		JobId:            job_id,
		FirstSeenAt:      timestamp,
		LastSeenListedAt: timestamp,
		NextScanAt:       &nextScan,
	}

	return Job_lifeCycle, nil

}

func Jobs_LifeCycleLiveRolesUpdater(records []map[string]string, filepath string) error {
	// following this format live-roles-38418f82-1d77-4760-8ce4-e40d43917d75_2026-04-29-19-23-53+00-00.csv
	//meta_data := strings.Split(strings.Split(strings.Split(filepath, "job_metadata-")[1], ".csv")[0], "_")
	meta_data := strings.Split(strings.Split(filepath, ".csv")[0], "_")
	timestamp, err := parseTimestamp(meta_data[1])

	if err != nil {
		return ErrorHandler(err, "failed to parse timestamp")
	}

	file_strings := strings.SplitN(meta_data[0], "-", 3)

	file_type := file_strings[0]
	uuid := file_strings[2]
	//TODO store uuid's in db sep table as done work
	fmt.Println("Executing this file now", uuid)
	db, err := ConnectDb()
	if err != nil {
		fmt.Println("db conn gone wrong", ErrorHandler(err, "you brought this on yourself"))
		return ErrorHandler(err, "whoops")
	}
	defer db.Close()
	// we want to scan live roles to capture the last time they were seen being live.
	if file_type == "live" {

		for index, record := range records {

			job_id, err := strconv.Atoi(record["job_id"])

			if err != nil {
				fmt.Println("something wrong", ErrorHandler(err, "you brought this on yourself"))
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
				fmt.Println("record extraction wrong", ErrorHandler(err, "you brought this on yourself"))
				return ErrorHandler(err, "whoops")
			}

			var suspended_count int
			err = db.QueryRow("SELECT suspended_count FROM JOB_LIFECYCLE WHERE job_id = $1", job_id).Scan(&suspended_count)

			if err != nil {
				if err == sql.ErrNoRows {
					fmt.Println("fetching record at index: has not been saved", index, ErrorHandler(err, "you brought this on yourself"))
					continue
				}
			}

			if job_state == "CLOSED" {

				_, err = db.Exec("UPDATE JOB_LIFECYCLE SET first_seen_closed_at = $1 , job_state = $2 WHERE job_id = $3", timestamp, job_state, job_id)

				if err != nil {
					//http.Error(w, " error updating Student ", http.StatusInternalServerError)
					fmt.Println("record at index ", index, " for live roles as not been saved", ErrorHandler(err, "Db query JobLifecycle update error"))
				}

			}

			if job_state == "SUSPENDED" {
				suspended_count += 1
				_, err = db.Exec("UPDATE JOB_LIFECYCLE SET last_seen_listed_at = $1 , suspended_count = $2 ,job_state = $3  WHERE job_id = $4", timestamp, suspended_count, job_state, job_id)

				if err != nil {
					//http.Error(w, " error updating Student ", http.StatusInternalServerError)
					fmt.Println("record at index ", index, " for live roles as not been saved", ErrorHandler(err, "Db query JobLifecycle update error"))
				}
			}
		}

	}

	return nil
}

func Jobs_LifeCycleSuspendedRolesUpdater(records []map[string]string, filepath string) error {

	//meta_data := strings.Split(strings.Split(strings.Split(filepath, "job_metadata-")[1], ".csv")[0], "_")
	meta_data := strings.Split(strings.Split(filepath, ".csv")[0], "_")
	timestamp, err := parseTimestamp(meta_data[1])

	db, err := ConnectDb()
	if err != nil {
		fmt.Println("db conn gone wrong", ErrorHandler(err, "you brought this on yourself"))
		return ErrorHandler(err, "whoops")
	}
	defer db.Close()

	for index, record := range records {
		job_id, err := strconv.Atoi(record["job_id"])

		job_state := record["job_state"]

		if err != nil {
			fmt.Println("record extraction wrong", ErrorHandler(err, "you brought this on yourself"))
			return ErrorHandler(err, "whoops")
		}

		var suspended_count int
		err = db.QueryRow("SELECT suspended_count FROM JOB_LIFECYCLE WHERE job_id = $1", job_id).Scan(suspended_count)

		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("fetching record at index: has not been saved", index, ErrorHandler(err, "you brought this on yourself"))
				continue
			}
		}

		if job_state == "CLOSED" {

			_, err = db.Exec("UPDATE JOB_LIFECYCLE SET first_seen_closed_at = $1 , job_state = $2 WHERE job_id = $3", timestamp, job_state, job_id)

			if err != nil {
				//http.Error(w, " error updating Student ", http.StatusInternalServerError)
				fmt.Println("record at index ", index, " for live roles as not been saved", ErrorHandler(err, "Db query JobLifecycle update error"))
			}

		}

		if job_state == "LISTED" {
			_, err = db.Exec("UPDATE JOB_LIFECYCLE SET last_seen_listed_at = $1 ,job_state = $2  WHERE job_id = $3", timestamp, job_state, job_id)
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

func getSearchTermIdHelper(searchTerm string, tablename string) (int, error) {
	db, err := ConnectDb()
	if err != nil {
		return -1, ErrorHandler(err, "db conn error")
	}
	defer db.Close()
	var searchtermtable string
	if tablename != "JOBS" {
		searchtermtable = "SEARCH_TERM_DEED"
	} else {
		searchtermtable = "SEARCH_TERM"
	}
	query := fmt.Sprintf("SELECT search_term_id FROM %s WHERE term = $1", searchtermtable)
	row := db.QueryRow(query, searchTerm)
	// row := db.QueryRow(`
	//     SELECT search_term_id FROM SEARCH_TERM
	//     WHERE term = $1
	// `, searchTerm)
	var search_term_id int
	if err := row.Scan(&search_term_id); err != nil {
		return -1, ErrorHandler(err, "row scan error")
	}
	return search_term_id, nil
}

func UpdateSearchWorkflowCounts(workflowid string, totalJobs int, netNew int, tablename string) error {
	db, err := ConnectDb()
	if err != nil {
		return ErrorHandler(err, "db conn error")
	}
	defer db.Close()
	switch tablename {
	case "JOBS":
		_, err = db.Exec(`UPDATE SEARCH_WORKFLOW SET total_jobs_found = $1, net_new_jobs = $2 WHERE workflow_id = $3`, totalJobs, netNew, workflowid)
	case "JOBS_ASH":
		_, err = db.Exec(`UPDATE SEARCH_WORKFLOW_ASH SET total_jobs_found = $1, net_new_jobs = $2 WHERE workflow_id = $3`, totalJobs, netNew, workflowid)
	case "JOBS_GREEN":
		_, err = db.Exec(`UPDATE SEARCH_WORKFLOW_GREEN SET total_jobs_found = $1, net_new_jobs = $2 WHERE workflow_id = $3`, totalJobs, netNew, workflowid)
	default:
		_, err = db.Exec(`UPDATE SEARCH_WORKFLOW_DEED SET total_jobs_found = $1, net_new_jobs = $2 WHERE workflow_id = $3`, totalJobs, netNew, workflowid)
	}

	if err != nil {
		return ErrorHandler(err, "update search workflow counts error")
	}

	return nil
}

// two pass logic on every row needed for self referencing parent id that might not exist yet . set null then set if exists.
func Team_loader_ash(records []map[string]string, tablename string, filepath string) error {

	db, err := ConnectDb()
	if err != nil {
		return ErrorHandler(err, "team loader db connection failed")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return ErrorHandler(err, "team loader begin tx failed")
	}
	defer tx.Rollback()

	seen := make(map[string]bool)
	var teams []models.COMPANY_ASH_TEAM

	for _, record := range records {
		// child team (the team this job posting is directly tagged with)
		team, err := TeamAshLoader(record)
		if err != nil {
			return ErrorHandler(err, "team parse failed")
		}
		if team.TeamId != "" && !seen[team.TeamId] {
			seen[team.TeamId] = true
			teams = append(teams, team)
		}

		// parent team — same row gives us its real id + name, add it too
		parentID := record["parentTeamId"]
		if parentID != "" && !seen[parentID] {
			seen[parentID] = true
			teams = append(teams, models.COMPANY_ASH_TEAM{
				TeamId:       parentID,
				TeamName:     record["parentTeam"],
				DepartmentId: record["departmentId"],
				ParentTeamId: nil,
			})
		}
	}

	for _, t := range teams {
		_, err := tx.Exec(`
            INSERT INTO TEAM_ASH (team_id, team_name, department_id, team_external_name)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (team_id) DO NOTHING
        `, t.TeamId, t.TeamName, t.DepartmentId, t.TeamExternalName)
		if err != nil {
			return ErrorHandler(err, "team insert (pass 1) failed")
		}
	}

	for _, t := range teams {
		if t.ParentTeamId == nil {
			continue
		}
		_, err := tx.Exec(`
            UPDATE TEAM_ASH SET parent_team_id = $1 WHERE team_id = $2
        `, t.ParentTeamId, t.TeamId)
		if err != nil {
			return ErrorHandler(err, "team parent backfill (pass 2) failed")
		}
	}

	return tx.Commit()
}
