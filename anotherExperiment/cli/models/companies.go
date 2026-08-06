package models

//todo handle uplaoding logo to db
type COMPANY struct {
	//Job_id      int       `json:"id,omitempty" db:"job_id"`
	CompanyId int    `json:"company_id,omitempty" db:"company_id"`
	Name      string `json:"name,omitempty" db:"name"`
	Logo      string `json:"logo,omitempty" db:"logo"`
}

type COMPANY_DEED struct {
	//Job_id      int       `json:"id,omitempty" db:"job_id"`
	CompanyId    string `json:"company_id,omitempty" db:"company_id"`
	Name         string `json:"company_name,omitempty" db:"name"`
	Employer_url string `json:"employer_url,omitempty" db:"employer_url"`
}

type COMPANY_ASH struct {
	//Job_id      int       `json:"id,omitempty" db:"job_id"`
	CompanyId         string  `json:"organizationId,omitempty" db:"company_id"`
	Name              string  `json:"companyName,omitempty" db:"name"`
	Public_website    *string `json:"publicWebsite,omitempty" db:"public_website"`
	CustomJobsPageUrl *string `json:"customJobsPageUrl,omitempty" db:"job_page_url"`
	Timezone          *string `json:"timezone,omitempty" db:"timezone"`
	Company_url       *string `json:"company_url,omitempty" db:"company_url"`
}

type COMPANY_GREEN struct {
	//Job_id      int       `json:"id,omitempty" db:"job_id"`
	CompanyId         int     `json:"company_id,omitempty" db:"company_id"`
	CompanyName       string  `json:"companyName,omitempty" db:"name"`
	JobBoardPublicUrl *string `json:"job_board_public_url,omitempty" db:"job_board_public_url"`
	DomainUrl         *string `json:"customJobsPageUrl,omitempty" db:"domainurl"`
	Company_url       *string `json:"company_url,omitempty" db:"company_url"`
	CompanyAbout      *string `json:"companyAbout,omitempty" db:"company_about"`
}

type COMPANY_ASH_TEAM struct {
	TeamId           string  `json:"teamId,omitempty" db:"team_id"`
	TeamName         string  `json:"teamName,omitempty" db:"team_name"`
	DepartmentId     string  `json:"departmentId,omitempty" db:"department_id"`
	ParentTeamId     *string `json:"parentTeamId,omitempty" db:"parent_team_id"`
	TeamExternalName *string `json:"teamExternalName,omitempty" db:"team_external_name"`
	CompanyId        string  `json:"organizationId,omitempty" db:"company_id"`
}

type CompanyGreenDepartment struct {
	DepartmentId    int    `json:"departmentId,omitempty" db:"department_id"`
	Department_name string `json:"DepartmentName,omitempty" db:"department_name"`
	ParentTeamId    *int64 `json:"parentTeamId,omitempty" db:"parent_id"`
	CompanyId       int    `json:"organizationId,omitempty" db:"company_id"`
}

type DeptNodeAsh struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	ExternalName *string `json:"externalName"`
	ParentTeamId *string `json:"parentTeamId"`
}

type DeptNode struct {
	ID       int64      `json:"id"`
	Value    int64      `json:"value"`
	Name     string     `json:"name"`
	Label    string     `json:"label"`
	Children []DeptNode `json:"children"`
}
