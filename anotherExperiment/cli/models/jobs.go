package models

import (
	"encoding/json"
	"time"
)

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
	// Is_repost       bool             `json:"is_repost,omitempty" db:"is_repost"`
	// Is_latest       bool             `json:"is_latest,omitempty" db:"is_latest"`
	Company_id      string           `json:"company_id,omitempty" db:"company_id"`
	Salary          string           `json:"salary,omitempty" db:"salary"`
	Location        string           `json:"location,omitempty" db:"location"`
	IndeedApplyable bool             `json:"indeedApplyable,omitempty" db:"indeed_applyable"`
	Taxonomy        *json.RawMessage `json:"taxonomy,omitempty" db:"taxonomy"`
}

type JobAsh struct {
	JobID            string           `json:"job_id" db:"job_id"`
	Title            string           `json:"title" db:"title"`
	IsListed         *bool            `json:"is_listed,omitempty" db:"is_listed"`
	DepartmentName   *string          `json:"department_name,omitempty" db:"department_name"`
	TeamNames        *json.RawMessage `json:"team_names,omitempty" db:"team_names"`
	TeamName         *string          `json:"teamName,omitempty" db:"team_name"`
	Team_id          *string          `json:"teamId,omitempty" db:"team_id"`
	ParentTeam_id    *string          `json:"parentTeamId,omitempty" db:"parent_team_id"`
	TeamExternalName *string          `json:"teamExternalName,omitempty" db:"team_external_name"`
	LocationName     *string          `json:"location_name,omitempty" db:"location_name"`
	EmploymentType   *string          `json:"employment_type,omitempty" db:"employment_type"`
	WorkplaceType    *string          `json:"workplace_type,omitempty" db:"workplace_type"`
	Salary           *string          `json:"salary,omitempty" db:"salary"`
	JobURL           string           `json:"job_url" db:"job_url"`
	CompanyID        string           `json:"company_id" db:"company_id"`
	OriginLinkID     *int64           `json:"origin_link_id,omitempty" db:"origin_link_id"`
	Origin_deed_id   *string          `json:"origin_deed_id,omitempty" db:"origin_deed_id"`
	Origin_ash       *string          `json:"origin_ash,omitempty" db:"origin_ash"`
}

type JobGreen struct {
	JobID          int              `json:"job_id" db:"job_id"`
	Title          string           `json:"title" db:"title"`
	LocationName   *string          `json:"locationName,omitempty" db:"location_name"`
	CompanyID      int              `json:"company_id" db:"company_id"`
	Company_name   string           `json:"company_name,omitempty" db:"company_name"`
	JobURL         string           `json:"job_url" db:"job_url"`
	RedirectURL    *string          `json:"redirect_url" db:"redirect_url"`
	PublishedDate  time.Time        `json:"publishedDate" db:"published_date"`
	DepartmentPath *json.RawMessage `json:"department_path,omitempty" db:"department_path"`
	Salary         *string          `json:"salary,omitempty" db:"salary"`
	OriginLinkID   *int64           `json:"origin_link_id,omitempty" db:"origin_link_id"`
	Origin_deed_id *string          `json:"origin_deed_id,omitempty" db:"origin_deed_id"`
	Origin_green   *string          `json:"origin,omitempty" db:"origin_green"`
	DepartmentId   *int64           `json:"department_id,omitempty" db:"department_id"`
}
type Everything struct {
	JOBS
	Jobs_metadata
	Company_Metadata
	COMPANY
	JobDescription
}
