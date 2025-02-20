package utils

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"
	"html/template"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))
var delifoodMail = os.Getenv("DELIFOOD_MAIL")
var delifoodPassword = os.Getenv("DELIFOOD_PASSWORD")

type Claims struct {
	UserID primitive.ObjectID `json:"user_id"`
	Email  string             `json:"email"`
	Role   string             `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret string

func InitJWTSecret() {
	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret_key"
	}
}

func GetJWTSecret() string {
	return jwtSecret
}

func GenerateJWT(userID primitive.ObjectID, email string, role string) (string, error) {
	expirationTime := time.Now().Add(12 * time.Hour)
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
	mail.SetHeader("From", os.Getenv("DELIFOOD_MAIL"))
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "Email Verification Code")
	mail.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s", code))

	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("DELIFOOD_MAIL"), os.Getenv("DELIFOOD_PASSWORD"))
	return dialer.DialAndSend(mail)
}

var TmplFuncs = template.FuncMap{
	"add": func(x, y int) int { return x + y },
	"sub": func(x, y int) int { return x - y },
	"iter": func(n int) []int {
		result := make([]int, n)
		for i := 0; i < n; i++ {
			result[i] = i + 1
		}
		return result
	},
}
