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
	err := ur.DB.QueryRow(`
		SELECT id, username, email, password, token, role 
		FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.Token, &user.Role)
	if err == sql.ErrNoRows {
		return models.User{}, fmt.Errorf("no user found with email: %s", email)
	}
	if err != nil {
		return models.User{}, fmt.Errorf("database query error: %w", err)
	}
	return user, nil
}

func (ur *UserRepo) UpdateUserToken(userID int, token string) error {
	_, err := ur.DB.Exec(`UPDATE users SET token = $1 WHERE id = $2`, token, userID)
	return err
}

func (ur *UserRepo) UpdateUserRole(email string, newRole string) error {
	// Update the user's role based on their email
	_, err := ur.DB.Exec(`
        UPDATE users 
        SET role = $1 
        WHERE email = $2`,
		newRole, email)

	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	return nil
}

func (ur *UserRepo) GetUserByToken(token string) (*models.User, error) {
	var user models.User
	var dbToken sql.NullString

	err := ur.DB.QueryRow(`
		SELECT id, username, email, password, token, role
		FROM users WHERE token = $1
	`, token).Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &dbToken, &user.Role)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no user found with token: %s", token)
	}

	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}

	// Handle NULL token
	if dbToken.Valid {
		user.Token = dbToken.String
	} else {
		user.Token = ""
	}

	return &user, nil
}

func (ur *UserRepo) AddFood(food models.Food) error {
	_, err := ur.DB.Exec(`
		INSERT INTO foods (name, category, image, description, price) 
		VALUES ($1, $2, $3, $4, $5)`,
		food.Name, food.Category, food.Image, food.Description, food.Price,
	)
	return err
}

func (ur *UserRepo) UpdateFood(food models.Food) error {
	_, err := ur.DB.Exec(`
		UPDATE foods SET name = $1, category = $2, image = $3, description = $4, price = $5 
		WHERE id = $6`,
		food.Name, food.Category, food.Image, food.Description, food.Price, food.ID,
	)
	return err
}

func (ur *UserRepo) DeleteFood(id int) error {
	_, err := ur.DB.Exec(`DELETE FROM foods WHERE id = $1`, id)
	return err
}

func (ur *UserRepo) GetFood(category, sortParam string) ([]models.Food, error) {
	// Base query
	query := "SELECT name, category, image, description, price FROM foods"

	// Add filtering by category
	var args []interface{}
	if category != "" {
		query += " WHERE category = $1"
		args = append(args, category)
	}

	// Add sorting
	switch sortParam {
	case "price-asc":
		query += " ORDER BY price ASC"
	case "price-desc":
		query += " ORDER BY price DESC"
	case "name":
		query += " ORDER BY name ASC"
	}

	// Execute query
	rows, err := ur.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve food items: %w", err)
	}
	defer rows.Close()

	// Parse rows into food items
	var foods []models.Food
	for rows.Next() {
		var food models.Food
		err := rows.Scan(&food.Name, &food.Category, &food.Image, &food.Description, &food.Price)
		if err != nil {
			return nil, fmt.Errorf("error scanning food row: %w", err)
		}
		foods = append(foods, food)
	}

	// Check for errors after iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return foods, nil
}
