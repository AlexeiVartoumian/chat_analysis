package models

import "encoding/json"

type JobDescription struct {
	JobId          int             `json:"job_id,omitempty" db:"job_id"`
	JobDescription string          `json:"job_description,omitempty" db:"job_description"`
	Encodings      json.RawMessage `json:"encodings,omitempty" db:"encodings"`
	Embedding      []float32       // post add after first batch insert
}
