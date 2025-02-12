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
