package users

import (
	"time"
)

type User struct {
	ID			string		`json:"id"`
	Username	string		`json:"username"`
	Password	string		`json:"password"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

type CreateUserPayload struct {
	Username	string		`json:"username"`
	Password	string		`json:"password"`
}

type LoginPayload struct {
	Username	string		`json:"username"`
	Password	string		`json:"password"`
}
