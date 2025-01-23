package models

type User struct {
	ID               int    `json:"id"`
	UserName         string `json:"name"`
	Password         string `json:"password"`
	Email            string `json:"email"`
	VerificationCode string
	IsVerified       bool   `json:"is_verified"`
	Token            string `json:"token"`
	Role             string `json:"role"`
}
