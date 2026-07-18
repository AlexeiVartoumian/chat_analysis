package models

import "encoding/json"

type JobDescription struct {
	JobId          int             `json:"job_id,omitempty" db:"job_id"`
	JobDescription string          `json:"job_description,omitempty" db:"job_description"`
	Encodings      json.RawMessage `json:"encodings,omitempty" db:"encodings"`
	Embedding      []float32       // post add after first batch insert
}

type JobDescription_DEED struct {
	JobId          string `json:"job_id,omitempty" db:"job_id"`
	JobDescription string `json:"job_description,omitempty" db:"job_description"`
}

type JobRedirect_DEED struct {
	JobId  string `json:"job_id,omitempty" db:"job_id"`
	JobUrl string `json:"job_url,omitempty" db:"job_url"`
	Origin string `json:"origin,omitempty" db:"origin"`
}

type JobRedirect_LinkGreen struct {
	JobId  string `json:"job_id,omitempty" db:"job_id"`
	JobUrl string `json:"job_url,omitempty" db:"job_url"`
	Origin string `json:"origin,omitempty" db:"origin"`
}

type JobRedirect_LinkAsh struct {
	JobId  string `json:"job_id,omitempty" db:"job_id"`
	JobUrl string `json:"job_url,omitempty" db:"job_url"`
	Origin string `json:"origin,omitempty" db:"origin"`
}
