package models

//todo handle uplaoding logo to db
type COMPANY struct {
	//Job_id      int       `json:"id,omitempty" db:"job_id"`
	CompanyId int    `json:"company_id,omitempty" db:"company_id"`
	Name      string `json:"name,omitempty" db:"name"`
	Logo      string `json:"logo,omitempty" db:"logo"`
}
