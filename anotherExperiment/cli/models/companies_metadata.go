package models

type Company_Metadata struct {
	CompanyId          int    `json:"company_id,omitempty" db:"company_id"`
	Name               string `json:"name,omitempty" db:"name"`
	Industry           string `json:"industry,omitempty" db:"industry"`
	Description        string `json:"description,omitempty" db:"description"`
	EmployeeCount      int    `json:"employee_count,omitempty" db:"employee_count"`
	EmployeeCountRange string `json:"employee_count_range,omitempty" db:"employee_count_range"`
}

type CompanyDeed_Metadata struct {
	CompanyId            string  `json:"company_id,omitempty" db:"company_id"`
	Name                 string  `json:"company_name,omitempty" db:"name"`
	Employee_count_range *string `json:"company_size,omitempty" db:"employee_count_range"`
	Industry             *string `json:"industry,omitempty" db:"industry"`
	Revenue              *string `json:"revenue,omitempty" db:"revenue"`
	Description          *string `json:"description,omitempty" db:"description"`
	Url                  *string `json:"company_url,omitempty" db:"url"`
}
