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

type JobDescription_Work struct {
	JobId          string `json:"job_id,omitempty" db:"job_id"`
	JobDescription string `json:"jobDescription,omitempty" db:"jobDescription"`
}

type JobDescription_Bambi struct {
	JobId             string `json:"job_id,omitempty" db:"job_id"`
	Description       string `json:"description,omitempty" db:"jobdescription"`
	EmploymentStatus  string `json:"employment_status,omitempty" db:"employment_status"`
	MinimumExperience string `json:"minimum_experience,omitempty" db:"minimum_experience"`
}

type JobDescription_Workab struct {
	JobId        string `json:"job_id,omitempty" db:"job_id"`
	Description  string `json:"description,omitempty" db:"description"`
	Requirements string `json:"requirements,omitempty" db:"requirements"`
	Benefits     string `json:"benefits,omitempty" db:"benefits"`
}

type JobDescription_Orac struct {
	JobId               string `json:"job_id,omitempty" db:"job_id"`
	ShortJobDescription string `json:"shortjobDescription,omitempty" db:"shortjobdescription"`
	JobDescription      string `json:"jobDescription,omitempty" db:"jobdescription"`
}

type JobDescription_Smart struct {
	JobId                 string `json:"job_id,omitempty" db:"job_id"`
	JobDescription        string `json:"jobDescription,omitempty" db:"jobdescription"`
	Qualifications        string `json:"qualifications,omitempty" db:"qualifications"`
	Additionalinformation string `json:"additional_information,omitempty" db:"additionalinformation"`
}

type JobDescription_Ash struct {
	JobId          string `json:"job_id,omitempty" db:"job_id"`
	JobDescription string `json:"descriptionHtml,omitempty" db:"description_html"`
}

type JobDescription_Green struct {
	JobId                    int     `json:"job_id,omitempty" db:"job_id"`
	JobDescription           string  `json:"job_description,omitempty" db:"job_description"`
	JobDescriptionConclusion *string `json:"job_description_conclusion,omitempty" db:"job_description_conclusion"`
	Company_about            *string `json:"company_about,omitempty" db:"company_about"`
}

type JobDescription_Lev struct {
	JobId           string          `json:"job_id,omitempty" db:"job_id"`
	JobDescription  string          `json:"jobDescription,omitempty" db:"jobDescription"`
	Requirements    json.RawMessage `json:"requirements,omitempty" db:"requirements"`
	AdditionalPlain string          `json:"additionalPlain,omitempty" db:"additionalPlain"`
}

type JobDescription_Rip struct {
	JobId           string `json:"job_id" db:"job_id"`
	DescriptionHtml string `json:"description_html" db:"description_html"`
}

type JobDescription_Tailor struct {
	JobId       string `json:"job_id" db:"job_id"`
	ContentHtml string `json:"content_html" db:"content_html"`
}
