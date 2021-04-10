package model

import "time"

type User struct {
	Id            int        `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	Created_at    *time.Time `json:"created_at"`
	Last_login_at *time.Time `json:"last_login_at"`
}

type UserRegistration struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
