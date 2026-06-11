package models

import (
	"encoding/json"
	"time"
)

type Company_Metadata struct {
	CompanyId          int    `json:"company_id,omitempty" db:"company_id"`
	Name               string `json:"name,omitempty" db:"name"`
	Industry           string `json:"industry,omitempty" db:"industry"`
	Description        string `json:"description,omitempty" db:"description"`
	EmployeeCount      int    `json:"employee_count,omitempty" db:"employee_count"`
	EmployeeCountRange string `json:"employee_count_range,omitempty" db:"employee_count_range"`
}

type CompanyDetail struct {
	CompanyId           int              `json:"company_id,omitempty" db:"company_id"`
	Name                string           `json:"name,omitempty" db:"name"`
	CompanySlug         string           `json:"company_slug,omitempty" db:"company_slug"`
	CompanyUrl          string           `json:"company_url,omitempty" db:"company_url"`
	Specialties         *json.RawMessage `json:"specialties,omitempty" db:"specialties"`
	Locations           *json.RawMessage `json:"locations,omitempty" db:"locations"`
	ExtendedDescription string           `json:"extended_description,omitempty" db:"extended_description"`
	StaffCount          int              `json:"staff_count,omitempty" db:"staff_count"`
	HeadquarterCity     string           `json:"headquarter_city,omitempty" db:"headquarter_city"`
	Created_at          time.Time        `json:"created_at,omitempty" db:"created_at"`
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
