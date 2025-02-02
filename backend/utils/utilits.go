package utils

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/golang-jwt/jwt/v5"
)

var SendVerificationEmailFunc = SendVerificationEmail

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

var jwtKey = []byte("your_secret_key")

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (c *Claims) Valid() error {
	if c.ExpiresAt.Time.Before(time.Now()) {
		return jwt.ErrTokenExpired
	}
	return nil
}

// jwtSecret is used to sign and validate JWT tokens
var jwtSecret string

func InitJWTSecret() {
	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret_key" // Replace with a strong key in production
	}
}

func GetJWTSecret() string {
	return jwtSecret
}

// GenerateJWT generates a JWT token
func GenerateJWT(userID int, email, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateJWT validates a JWT and returns claims
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func SendVerificationEmail(email, code string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "mr.akhmedali@bk.ru")
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "Email Verification Code")
	mail.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s", code))

	dialer := gomail.NewDialer("smtp.mail.ru", 587, "mr.akhmedali@bk.ru", "LVWZUunmUvMW8giSXLe0")
	return dialer.DialAndSend(mail)
}
