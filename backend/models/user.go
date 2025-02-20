package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserName         string             `bson:"username" json:"name"`
	Password         string             `bson:"password" json:"password"`
	Email            string             `bson:"email" json:"email"`
	VerificationCode string             `bson:"verification_code" json:"verificationCode"`
	IsVerified       bool               `bson:"is_verified" json:"is_verified"`
	Role             string             `bson:"role" json:"role"`
	Token            string             `bson:"token" json:"token"`
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)
