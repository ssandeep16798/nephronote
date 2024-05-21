package models

import "time"

type UserRegistrationForm struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
	UserName  string    `json:"username"`
	Password  string    `json:"password"`
	Gender    string    `json:"gender"`
	Mobile    string    `json:"mobile"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRegisterResponse struct {
	UserName string `json:"email"`
	Message  string `json:"message"`
}

type AuthenticatedUser struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}
