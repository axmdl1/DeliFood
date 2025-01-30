package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestGenerateJWT(t *testing.T) {
	jwtKey = []byte("test_secret_key")

	userID := 1
	email := "test@example.com"
	role := "user"
	//generate token
	tokenString, err := GenerateJWT(userID, email, role)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		t.Fatalf("Failed to parse JWT: %v", err)
	}

	if !token.Valid {
		t.Fatal("Generated JWT is not valid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		t.Fatal("Failed to parse claims")
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}

	if time.Until(claims.ExpiresAt.Time) > 24*time.Hour || time.Until(claims.ExpiresAt.Time) < 23*time.Hour {
		t.Errorf("Token expiration time is incorrect: %v", claims.ExpiresAt.Time)
	}
}
