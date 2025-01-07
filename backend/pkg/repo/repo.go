package repo

import (
	"DeliFood/backend/pkg/logger"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AddUser(db *sql.DB, user User, log *logger.Logger) error {
	log.Info("Starting to add a new user", map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	})

	// Hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	stmt, err := db.Prepare("INSERT INTO users(name, surname, username, email, password) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		log.Error("Failed to prepare SQL statement", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Surname, user.Username, user.Email, user.Password)
	if err != nil {
		log.Error("Failed to execute SQL statement", map[string]interface{}{
			"error":    err.Error(),
			"username": user.Username,
			"email":    user.Email,
		})
		return fmt.Errorf("failed to add user: %w", err)
	}

	log.Info("User added successfully", map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	})
	return nil
}
