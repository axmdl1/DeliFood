package handlers

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// Проверка хеширования пароля
func TestPasswordHashing(t *testing.T) {
	password := "securepassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	// Проверка, что хеш не совпадает с оригиналом
	if string(hashedPassword) == password {
		t.Errorf("Hashed password should not be the same as plain password")
	}

	// Проверка верификации пароля
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		t.Errorf("Password рщverification failed")
	}
}
