package models

type User struct {
	ID               int    `json:"id"`
	UserName         string `json:"name"`
	Password         string `json:"password"`
	Email            string `json:"email"`
	VerificationCode string `json:"verificationCode"`
	IsVerified       bool   `json:"is_verified"`
	Role             string `json:"role"`
	Token            string `json:"token"`
}
