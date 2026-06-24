package models

import "database/sql"

type CompanyCorrelation struct {
	Company_link_id   int            `json:"company_id,omitempty" db:"company_link_id"`
	Company_deed_id   string         `json:"companydeed_id,omitempty" db:"company_deed_id"`
	Company_name      string         `json:"company_name,omitempty" db:"company_name"`
	Company_link_url  sql.NullString `json:"company_link_url,omitempty" db:"company_link_url"`
	Company_deed_url  sql.NullString `json:"company_deed_url,omitempty" db:"company_deed_url"`
	Company_apply_url sql.NullString `json:"company_apply_url,omitempty" db:"company_apply_url"`
	JobDeed_url       sql.NullString `json:"jobdeed_url,omitempty" db:"job_deed_url"`
}
