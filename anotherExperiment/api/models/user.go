package models

type User struct {
	UserID       int    `json:"user_id" db:"user_id"`
	FirstName    string `json:"first_name" db:"first_name"`
	LastName     string `json:"last_name" db:"last_name"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	Age          int    `json:"age" db:"age"` //use age as Poc for zero knowledge proof
}
