package models

type JOB_DESCRIPTION struct {
	JobId          int    `json:"job_id,omitempty" db:"job_id"`
	JobDescription string `json:"job_description,omitempty" db:"job_description"`
}
