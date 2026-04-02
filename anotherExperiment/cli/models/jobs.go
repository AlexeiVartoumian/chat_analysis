package models

import "time"

type JOBS struct {
	Job_id      int       `json:"id,omitempty" db:"job_id"`
	Title       string    `json:"title,omitempty" db:"title"`
	Location    string    `json:"location,omitempty" db:"location"`
	Salary      string    `json:"salary,omitempty" db:"salary"`
	Date_Posted time.Time `json:"date_posted,omitempty" db:"date_posted"`
	Job_url     string    `json:"job_url,omitempty" db:"job_url"`
	Search_term string    `json:"search_term,omitempty" db:"search_term"`
	Easy_apply  bool      `json:"easy_apply,omitempty" db:"easy_apply"`
	Promoted    bool      `json:"promoted,omitempty" db:"promoted"`
	Expiry_Date time.Time `json:"expiry_date,omitempty" db:"expiry_date"`
	Company_id  int       `json:"company_id,omitempty" db:"company_id"`
}

type Everything struct {
	JOBS
	Jobs_metadata
	Company_Metadata
	COMPANY
	JOB_DESCRIPTION
}
