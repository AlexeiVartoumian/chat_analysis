package models

import "time"

//todo handle uplaoding logo to db
type COMPANY struct {
	//Job_id      int       `json:"id,omitempty" db:"job_id"`
	CompanyId int    `json:"company_id,omitempty" db:"company_id"`
	Name      string `json:"name,omitempty" db:"name"`
	Logo      string `json:"logo,omitempty" db:"logo"`
}

type COMPANY_DEED struct {
	//Job_id      int       `json:"id,omitempty" db:"job_id"`
	CompanyId    string `json:"company_id,omitempty" db:"company_id"`
	Name         string `json:"company_name,omitempty" db:"name"`
	Employer_url string `json:"employer_url,omitempty" db:"employer_url"`
}

type CompanyDetail struct {
	Company_id   int    `json:"company_id"`
	Company_name string `json:"name"`
}

type AshCompany struct {
	Teamname       string
	DepartmentName string
	LocationName   string
	CompanyUrl     string
	CompanyId      string
	Lastscannedat  time.Time
}

type GreenbyCompany struct {
	JobId         int    `json:"job_id,omitempty" db:"job_id"`
	CompanyId     int    `json:"company_id,omitempty" db:"company_id"`
	CompanyName   string `json:"company_name,omitempty" db:"company_name"`
	JobUrl        string `json:"job_url,omitempty" db:"job_url"`
	CompanyUrl    string `json:"job_board_public_url,omitempty" db:"job_board_public_url"`
	Lastscannedat time.Time
}
