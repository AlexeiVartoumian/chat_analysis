package sqlconnect

import (
	"api/models"
	"api/utils"
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
		//return utils.ErrorHandler(err, ""), jobModel
		return utils.ErrorHandler(err, "")
	}
	defer file.Close()

	records, err := gocsv.CSVToMaps(file)
	if err != nil {

		return utils.ErrorHandler(err, "uh oh")
	}

	for index, record := range records {
		fmt.Println(record)
		fmt.Println(" ")
		value, err := ModelLoader(tablename, record)

		if err != nil {
			//return utils.ErrorHandler(err, "you brought this on yourself")
			//dont return process other records
			fmt.Println("record at index: has not been saved", index, utils.ErrorHandler(err, "you brought this on yourself"))
			continue
		}

		AddNewRow(value, tablename)

	}

	return nil
}

func ModelLoader(tablename string, record map[string]string) (interface{}, error) {

	switch tablename {
	case "COMPANY":
		return CompanyLoader(record)
	case "COMPANY_METADATA":
		return Company_MetadataLoader(record)
	case "JOBS":
		return JobLoader(record)
	case "JOBS_METADATA":
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
		return models.JOB_DESCRIPTION{}, utils.ErrorHandler(err, "whoops")
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
		return models.Jobs_metadata{}, utils.ErrorHandler(err, "uh oh")
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
		return models.Company_Metadata{}, utils.ErrorHandler(err, "uh oh error parsing employee count ")
	}

	Company_id, err := strconv.Atoi(record["company_id"])

	if err != nil {
		return models.Company_Metadata{}, utils.ErrorHandler(err, "uh oh error parsing company id")
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
		return models.COMPANY{}, utils.ErrorHandler(nil, "nil value")
	}
	company_id, err := strconv.Atoi(record["company_id"])
	if err != nil {
		return models.COMPANY{}, utils.ErrorHandler(err, "uh oh Company id fail parse")

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
		return models.JOBS{}, utils.ErrorHandler(err, "uh oh jobid id fail parse")
	}

	var Company_id int

	_, err2 := strconv.Atoi(record["company_id"])
	if err2 != nil {
		//could be a solo person posting the job

	} else {
		Company_id, err = strconv.Atoi(record["company_id"])
		if err != nil {
			return models.JOBS{}, utils.ErrorHandler(err, "uh oh Company id fail parse")

		}
	}

	easy_apply, err := strconv.ParseBool(record["easy_apply"])

	if err != nil {
		return models.JOBS{}, utils.ErrorHandler(err, "uh oh easy apply fail bool parse")
	}

	promoted, err := strconv.ParseBool(record["promoted"])

	if err != nil {
		return models.JOBS{}, utils.ErrorHandler(err, "uh oh prmoted fail bool parse")
	}

	datePosted := record["posted_date"]
	fmt.Println(datePosted, "value")
	time1, err := time.Parse("2006-01-02", datePosted)

	if err != nil {
		return models.JOBS{}, utils.ErrorHandler(err, "something happened with the date")
	}

	job := models.JOBS{
		Job_id:      job_id,
		Title:       record["title"],
		Location:    record["location"],
		Salary:      record["salary"],
		Date_Posted: time1,
		Job_url:     record["job_url"],
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
