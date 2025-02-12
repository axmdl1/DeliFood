package repo

import (
	"DeliFood/backend/models"
	"database/sql"
	"fmt"
)

type AdminRepo struct {
	DB *sql.DB
}

func NewAdminRepo(db *sql.DB) *AdminRepo {
	return &AdminRepo{DB: db}
}

func (ar *AdminRepo) AddFood(food *models.Food) error {
	query := `INSERT INTO foods (name, category, image, description, price) VALUES ($1, $2, $3, $4, $5)`

	_, err := ar.DB.Exec(query, food.Name, food.Category, food.Image, food.Description, food.Price)
	if err != nil {
		return fmt.Errorf("Error inserting food: %s", err)
	}

	return nil
}

func (ar *AdminRepo) UpdateFood(food models.Food) error {
	_, err := ar.DB.Exec(`
		UPDATE foods SET name = $1, category = $2, image = $3, description = $4, price = $5 
		WHERE id = $6`,
		food.Name, food.Category, food.Image, food.Description, food.Price, food.ID,
	)
	return err
}

// UpdateUserRole updates the role of a user in the database
func (ar *AdminRepo) UpdateUserRole(userID int, role string) error {
	_, err := ar.DB.Exec("UPDATE users SET role = $1 WHERE id = $2", role, userID)
	return err
}

// GetAllFoods fetches all food items from the database
func (ar *AdminRepo) GetAllFoods() ([]models.Food, error) {
	var foods []models.Food

	// Query the database to get all food details
	rows, err := ar.DB.Query(`
		SELECT id, name, category, image, description, price 
		FROM foods`)
	if err != nil {
		return nil, fmt.Errorf("error fetching foods: %w", err)
	}
	defer rows.Close()

	// Iterate through the rows and scan each food item into the foods slice
	for rows.Next() {
		var food models.Food
		if err := rows.Scan(&food.ID, &food.Name, &food.Category, &food.Image, &food.Description, &food.Price); err != nil {
			return nil, fmt.Errorf("error scanning food row: %w", err)
		}
		foods = append(foods, food)
	}

	// Check for any error that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return foods, nil
}
