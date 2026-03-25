package models

type Application struct {
	ApplicationID int    `json:"application_id" db:"application_id"`
	Status        string `json:"status" db:"status"`
	UserID        int    `json:"user_id" db:"user_id"`
	JobId         int    `json:"job_id,omitempty" db:"job_id"`
	Age           int    `json:"age" db:"age"` //use age as Poc for zero knowledge proof
}

type ApplicationStatus string

const (
	StatusLive    = "live"
	StatusExpired = "expired"
)
