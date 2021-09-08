package model

import "time"

type User struct {
	Id        int
	Username  string
	Email     string
	CreatedAt *time.Time
}

type UserLogin struct {
	Username string
	Password string
}

type UserRegistration struct {
	Username string
	Password string
	Email    string
}
