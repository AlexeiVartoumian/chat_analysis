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

type JobAsh struct {
	JobID          string  `json:"job_id" db:"job_id"`
	Title          string  `json:"title" db:"title"`
	IsListed       *bool   `json:"is_listed,omitempty" db:"is_listed"`
	DepartmentName *string `json:"department_name,omitempty" db:"department_name"`
	TeamName       *string `json:"team_name,omitempty" db:"team_name"`
	LocationName   *string `json:"location_name,omitempty" db:"location_name"`
	EmploymentType *string `json:"employment_type,omitempty" db:"employment_type"`
	WorkplaceType  *string `json:"workplace_type,omitempty" db:"workplace_type"`
	// PublishedDate  time.Time `json:"published_date" db:"published_date"`
	// UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	Salary         *string `json:"salary,omitempty" db:"salary"`
	JobURL         string  `json:"job_url" db:"job_url"`
	CompanyID      string  `json:"company_id" db:"company_id"`
	OriginLinkID   *int64  `json:"origin_link_id,omitempty" db:"origin_link_id"`
	Origin_deed_id *string `json:"origin_deed_id,omitempty" db:"origin_deed_id"`
}

type JobGreen struct {
	JobID          int       `json:"job_id" db:"job_id"`
	Title          string    `json:"title" db:"title"`
	LocationName   *string   `json:"locationName,omitempty" db:"location_name"`
	CompanyID      int       `json:"company_id" db:"company_id"`
	Company_name   string    `json:"company_name,omitempty" db:"company_name"`
	JobURL         string    `json:"job_url" db:"job_url"`
	RedirectURL    *string   `json:"redirect_url" db:"redirect_url"`
	PublishedDate  time.Time `json:"publishedDate" db:"published_date"`
	Salary         *string   `json:"salary,omitempty" db:"salary"`
	OriginLinkID   *int64    `json:"origin_link_id,omitempty" db:"origin_link_id"`
	Origin_deed_id *string   `json:"origin_deed_id,omitempty" db:"origin_deed_id"`
}

type Everything struct {
	JOBS
	Jobs_metadata
	Company_Metadata
	COMPANY
	JobDescription
}
