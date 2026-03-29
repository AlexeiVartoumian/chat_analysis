package models

type Jobs_metadata struct {
	JobId           int    `json:"job_id,omitempty" db:"job_id"`
	ApplicantsCount string `json:"applicants_count,omitempty" db:"applicants_count"`
	CompanyApplyUrl string `json:"company_apply_url,omitempty" db:"company_apply_url"`
	JobState        string `json:"job_state,omitempty" db:"job_state"`
}
