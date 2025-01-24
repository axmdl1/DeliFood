package repo

import (
	"DeliFood/backend/models"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

// GetUserByEmail retrieves a user from the database by their email.
func (ur *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.DB.QueryRow(`
		SELECT id, username, email, password, verificationcode, isverified, role 
		FROM users WHERE email = $1`, email).Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.VerificationCode, &user.IsVerified, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with email: %s", email)
		}
		return nil, fmt.Errorf("error fetching user by email: %w", err)
	}

	return &user, nil
}

// Register user
func (ur *UserRepo) Register(user models.User) error {
	// Check for duplicate email or username
	var exists bool
	err := ur.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR username = $2)", user.Email, user.UserName).Scan(&exists)
	if err != nil || exists {
		return errors.New("email or username already exists")
	}

	// Hash the password
	/*hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}*/

	// Insert user into database
	_, err = ur.DB.Exec("INSERT INTO users (username, email, password, verificationcode, isverified, role) VALUES ($1, $2, $3, $4, $5, $6)",
		user.UserName, user.Email, user.Password, user.VerificationCode, false, user.Role)
	return err
}

func (ur *UserRepo) UpdateVerificationStatus(userID int, isVerified bool) error {
	_, err := ur.DB.Exec(`UPDATE users SET isverified = $1 WHERE id = $2`, isVerified, userID)
	return err
}

func (ur *UserRepo) CheckEmailOrUsernameExists(email, username string) (bool, error) {
	var count int
	err := ur.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM users 
		WHERE email = $1 OR username = $2
	`, email, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Verify Email
func (ur *UserRepo) VerifyEmail(email, code string) error {
	var dbCode string
	err := ur.DB.QueryRow("SELECT verificationcode FROM users WHERE email = $1", email).Scan(&dbCode)
	if err != nil || dbCode != code {
		return errors.New("invalid verification code")
	}

	_, err = ur.DB.Exec("UPDATE users SET isverified = TRUE WHERE email = $1", email)
	return err
}

// Authenticate user
func (ur *UserRepo) Authenticate(email, password string) (models.User, error) {
	var user models.User

	err := ur.DB.QueryRow("SELECT id, username, email, password, role, isverified FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.Role, &user.IsVerified)

	if err == sql.ErrNoRows {
		return models.User{}, errors.New("invalid email or password")
	}

	if !user.IsVerified {
		return models.User{}, errors.New("email not verified")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Password mismatch for email:", email)
		return models.User{}, errors.New("invalid email or password")
	}

	return user, nil
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
