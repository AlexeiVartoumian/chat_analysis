package models

type CompanyCorrelation struct {
	Company_link_id  int    `json:"company_id,omitempty" db:"company_link_id"`
	Company_deed_id  string `json:"company_id,omitempty" db:"company_deed_id"`
	Company_name     string `json:"company_name,omitempty" db:"company_name"`
	Company_link_url string `json:"company_link_url,omitempty" db:"company_link_url"`
	Company_deed_url string `json:"company_deed_url,omitempty" db:"company_deed_url"`
}
