package repo

import (
	"DeliFood/backend/models"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (ur *UserRepo) SignUp(user models.User) error {
	_, err := ur.DB.Exec("INSERT INTO users (username, email, password, verificationcode, isverified) VALUES ($1, $2, $3, $4, $5)",
		user.UserName, user.Email, user.Password, user.VerificationCode, false)
	if err != nil {
		return fmt.Errorf("failed to insert user %w", err)
	}
	return nil
}

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (ur *UserRepo) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	var token sql.NullString // Use sql.NullString to handle NULL values in the database

	err := ur.DB.QueryRow(`
		SELECT id, username, email, password, token 
		FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &token)

	if err == sql.ErrNoRows {
		return models.User{}, fmt.Errorf("no user found with email: %s", email)
	}

	if err != nil {
		return models.User{}, fmt.Errorf("database query error: %w", err)
	}

	// Convert sql.NullString to string, handle NULL case
	if token.Valid {
		user.Token = token.String
	} else {
		user.Token = "" // Default value for NULL token
	}

	return user, nil
}

func (ur *UserRepo) UpdateUserToken(userID int, token string) error {
	_, err := ur.DB.Exec(`UPDATE users SET token = $1 WHERE id = $2`, token, userID)
	return err
}

func (ur *UserRepo) GetUserByToken(token string) (*models.User, error) {
	var user models.User
	err := ur.DB.QueryRow(`SELECT id, username, email, password, token FROM users WHERE token = $1`, token).
		Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.Token)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
