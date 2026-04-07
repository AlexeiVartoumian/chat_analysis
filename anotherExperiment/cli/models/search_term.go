package models

import "time"

type Search_term struct {
	//search term id generated
	Search_term string `json:"search_term,omitempty" db:"search_term"`
}

type SearchWorkflow struct {
	Workflow_id      string    `json:"workflow_id,omitempty" db:"workflow_id"`
	Search_term_id   int       `json:"search_term,omitempty" db:"search_term_id"`
	Run_at           time.Time `json:"run_at,omitempty" db:"run_at"`
	Total_jobs_found int       `json:"total_jobs_found,omitempty" db:"total_jobs_found"`
	Net_new_found    int       `json:"net_new_jobs,omitempty" db:"net_new_jobs_found"`
}

type JOB_SEARCH_TERM struct {
	Job_id      int    `json:"id,omitempty" db:"job_id"`
	Workflow_id string `json:"workflow_id,omitempty" db:"workflow_id"`
}
