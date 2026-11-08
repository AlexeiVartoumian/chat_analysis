package models

import "time"

type Search_term struct {
	//search term id generated
	Search_term string `json:"search_term,omitempty" db:"term"`
}

type SearchWorkflow struct {
	Workflow_id      string    `json:"workflow_id,omitempty" db:"workflow_id"`
	Search_term_id   int       `json:"search_term,omitempty" db:"search_term_id"`
	Run_at           time.Time `json:"run_at,omitempty" db:"run_at"`
	Total_jobs_found int       `json:"total_jobs_found,omitempty" db:"total_jobs_found"`
	Net_new_found    int       `json:"net_new_jobs,omitempty" db:"net_new_jobs"`
}

type SearchWorkflowAsh struct {
	Workflow_id      string    `json:"workflow_id,omitempty" db:"workflow_id"`
	Run_at           time.Time `json:"run_at,omitempty" db:"run_at"`
	Total_jobs_found int       `json:"total_jobs_found,omitempty" db:"total_jobs_found"`
	Net_new_found    int       `json:"net_new_jobs,omitempty" db:"net_new_jobs"`
}

type SearchWorkflowGreen struct {
	Workflow_id      string    `json:"workflow_id,omitempty" db:"workflow_id"`
	Run_at           time.Time `json:"run_at,omitempty" db:"run_at"`
	Total_jobs_found int       `json:"total_jobs_found,omitempty" db:"total_jobs_found"`
	Net_new_found    int       `json:"net_new_jobs,omitempty" db:"net_new_jobs"`
}

type SearchWorkflowWork struct {
	Workflow_id      string    `json:"workflow_id,omitempty" db:"workflow_id"`
	Run_at           time.Time `json:"run_at,omitempty" db:"run_at"`
	Total_jobs_found int       `json:"total_jobs_found,omitempty" db:"total_jobs_found"`
	Net_new_found    int       `json:"net_new_jobs,omitempty" db:"net_new_jobs"`
}

type SearchWorkflowLev struct {
	Workflow_id      string    `json:"workflow_id,omitempty" db:"workflow_id"`
	Run_at           time.Time `json:"run_at,omitempty" db:"run_at"`
	Total_jobs_found int       `json:"total_jobs_found,omitempty" db:"total_jobs_found"`
	Net_new_found    int       `json:"net_new_jobs,omitempty" db:"net_new_jobs"`
}

type SearchWorkflowOrac struct {
	Workflow_id      string    `json:"workflow_id,omitempty" db:"workflow_id"`
	Run_at           time.Time `json:"run_at,omitempty" db:"run_at"`
	Total_jobs_found int       `json:"total_jobs_found,omitempty" db:"total_jobs_found"`
	Net_new_found    int       `json:"net_new_jobs,omitempty" db:"net_new_jobs"`
}

type SearchWorkflowSmart struct {
	Workflow_id      string    `json:"workflow_id,omitempty" db:"workflow_id"`
	Run_at           time.Time `json:"run_at,omitempty" db:"run_at"`
	Total_jobs_found int       `json:"total_jobs_found,omitempty" db:"total_jobs_found"`
	Net_new_found    int       `json:"net_new_jobs,omitempty" db:"net_new_jobs"`
}

type SearchWorkflowBambi struct {
	Workflow_id      string    `json:"workflow_id,omitempty" db:"workflow_id"`
	Run_at           time.Time `json:"run_at,omitempty" db:"run_at"`
	Total_jobs_found int       `json:"total_jobs_found,omitempty" db:"total_jobs_found"`
	Net_new_found    int       `json:"net_new_jobs,omitempty" db:"net_new_jobs"`
}

type SearchWorkflowWorkab struct {
	Workflow_id      string    `json:"workflow_id,omitempty" db:"workflow_id"`
	Run_at           time.Time `json:"run_at,omitempty" db:"run_at"`
	Total_jobs_found int       `json:"total_jobs_found,omitempty" db:"total_jobs_found"`
	Net_new_found    int       `json:"net_new_jobs,omitempty" db:"net_new_jobs"`
}

type SearchWorkflowTailor struct {
	Workflow_id      string    `db:"workflow_id"`
	Run_at           time.Time `db:"run_at"`
	Total_jobs_found int       `db:"total_jobs_found"`
	Net_new_found    int       `db:"net_new_jobs"`
}

type JOB_SEARCH_TERM struct {
	Job_id      int    `json:"id,omitempty" db:"job_id"`
	Workflow_id string `json:"workflow_id,omitempty" db:"workflow_id"`
	Is_new_job  bool   `json:"is_new_job,omitempty" db:"is_new_job"`
}

type JOB_SEARCH_TERM_DEED struct {
	Job_id      string `json:"id,omitempty" db:"job_id"`
	Workflow_id string `json:"workflow_id,omitempty" db:"workflow_id"`
	Is_new_job  bool   `json:"is_new_job,omitempty" db:"is_new_job"`
}

type JOB_SEARCH_TERM_ASH struct {
	Job_id      string `json:"id,omitempty" db:"job_id"`
	Workflow_id string `json:"workflow_id,omitempty" db:"workflow_id"`
	Is_new_job  bool   `json:"is_new_job,omitempty" db:"is_new_job"`
}

type JOB_SEARCH_TERM_GREEN struct {
	Job_id      int    `json:"id,omitempty" db:"job_id"`
	Workflow_id string `json:"workflow_id,omitempty" db:"workflow_id"`
	Is_new_job  bool   `json:"is_new_job,omitempty" db:"is_new_job"`
}
type JOB_SEARCH_TERM_WORK struct {
	Job_id      string `json:"id,omitempty" db:"job_id"`
	Workflow_id string `json:"workflow_id,omitempty" db:"workflow_id"`
	Is_new_job  bool   `json:"is_new_job,omitempty" db:"is_new_job"`
}

type JOB_SEARCH_TERM_BAMBI struct {
	Job_id      string `json:"id,omitempty" db:"job_id"`
	Workflow_id string `json:"workflow_id,omitempty" db:"workflow_id"`
	Is_new_job  bool   `json:"is_new_job,omitempty" db:"is_new_job"`
}

type JOB_SEARCH_TERM_WORKAB struct {
	Job_id      string `json:"id,omitempty" db:"job_id"`
	Workflow_id string `json:"workflow_id,omitempty" db:"workflow_id"`
	Is_new_job  bool   `json:"is_new_job,omitempty" db:"is_new_job"`
}

type JOB_SEARCH_TERM_LEV struct {
	Job_id      string `json:"id,omitempty" db:"job_id"`
	Workflow_id string `json:"workflow_id,omitempty" db:"workflow_id"`
	Is_new_job  bool   `json:"is_new_job,omitempty" db:"is_new_job"`
}

type JOB_SEARCH_TERM_ORAC struct {
	Job_id      string `json:"id,omitempty" db:"job_id"`
	Workflow_id string `json:"workflow_id,omitempty" db:"workflow_id"`
	Is_new_job  bool   `json:"is_new_job,omitempty" db:"is_new_job"`
}

type JOB_SEARCH_TERM_SMART struct {
	Job_id      string `json:"id,omitempty" db:"job_id"`
	Workflow_id string `json:"workflow_id,omitempty" db:"workflow_id"`
	Is_new_job  bool   `json:"is_new_job,omitempty" db:"is_new_job"`
}

type JOB_SEARCH_TERM_RIP struct {
	Job_id      string `db:"job_id"`
	Workflow_id string `db:"workflow_id"`
	Is_new_job  bool   `db:"is_new_job"`
}

type File_key struct {
	//search term id generated
	File_name string `json:"file_name,omitempty" db:"file_name"`
}

type SearchWorkflowRip struct {
	Workflow_id      string    `db:"workflow_id"`
	Run_at           time.Time `db:"run_at"`
	Total_jobs_found int       `db:"total_jobs_found"`
	Net_new_found    int       `db:"net_new_jobs"`
}

type JOB_SEARCH_TERM_TAILOR struct {
	Job_id      string `db:"job_id"`
	Workflow_id string `db:"workflow_id"`
	Is_new_job  bool   `db:"is_new_job"`
}
