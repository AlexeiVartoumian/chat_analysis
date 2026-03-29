package models

type Company_Metadata struct {
	CompanyId          int    `json:"company_id,omitempty" db:"company_id"`
	Name               string `json:"name,omitempty" db:"name"`
	Industry           string `json:"industry,omitempty" db:"industry"`
	Description        string `json:"description,omitempty" db:"description"`
	EmployeeCount      int    `json:"employee_count,omitempty" db:"employee_count"`
	EmployeeCountRange string `json:"employee_count_range,omitempty" db:"employee_count_range"`
}
