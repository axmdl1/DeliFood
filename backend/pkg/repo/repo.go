package repo

import (
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

func AddUser(db *sql.DB, user User) error {
	//Hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	stmt, err := db.Prepare("INSERT INTO users(name, surname, username, email, password) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Surname, user.Username, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("Failed to add user: %w", err)
	}

	return nil
}
