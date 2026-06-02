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

type JOBS_DEED struct {
	Job_id          string    `json:"id,omitempty" db:"job_id"`
	Title           string    `json:"title,omitempty" db:"title"`
	Date_published  time.Time `json:"date_published,omitempty" db:"date_published"`
	Date_advertised time.Time `json:"date_advertised,omitempty" db:"date_advertised"`
	Job_url         string    `json:"job_url,omitempty" db:"job_url"`
	Search_term     string    `json:"search_term,omitempty" db:"search_term"`
	Organic_apply   int       `json:"organic_apply,omitempty" db:"organic_apply"`
	Is_repost       bool      `json:"is_repost,omitempty" db:"is_repost"`
	Is_latest       bool      `json:"is_latest,omitempty" db:"is_latest"`
	Company_id      string    `json:"company_id,omitempty" db:"company_id"`
}

type Everything struct {
	JOBS
	Jobs_metadata
	Company_Metadata
	COMPANY
	JobDescription
}
