package models

type CompanyCorrelation struct {
	Company_deed_id string `json:"company_id,omitempty" db:"company_deed_id"`
	Company_link_id int    `json:"company_id,omitempty" db:"company_link_id"`
}
