package models

import "time"

type User struct {
	ID           int
	Login        string
	PasswordHash string
	CreatedAt    time.Time
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
