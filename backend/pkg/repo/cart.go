package repo

import (
	"DeliFood/backend/models"
	"database/sql"
	"fmt"
)

// CartRepo is the repository that interacts with the cart_items table
type CartRepo struct {
	DB *sql.DB
}

// NewCartRepo creates a new CartRepo instance
func NewCartRepo(db *sql.DB) *CartRepo {
	return &CartRepo{DB: db}
}

// AddItemToCart adds an item to the user's cart
func (repo *CartRepo) AddItemToCart(userID int, foodID int, quantity int, foodName string, foodPrice float64) error {
	// Store the cart item in the database
	_, err := repo.DB.Exec(`
		INSERT INTO cart_items (user_id, food_id, quantity, food_name, food_price)
		VALUES ($1, $2, $3, $4, $5)`,
		userID, foodID, quantity, foodName, foodPrice)
	if err != nil {
		return fmt.Errorf("failed to add item to cart: %w", err)
	}
	return nil
}

// UpdateItemQuantity updates the quantity of an item in the cart
func (cr *CartRepo) UpdateItemQuantity(userID, foodID, quantity int) error {
	_, err := cr.DB.Exec("UPDATE cart_items SET quantity = $1, updated_at = CURRENT_TIMESTAMP WHERE user_id = $2 AND food_id = $3", quantity, userID, foodID)
	return err
}

// RemoveItemFromCart removes an item from the cart
func (cr *CartRepo) RemoveItemFromCart(userID int, itemID int) error {
	_, err := cr.DB.Exec(`
		DELETE FROM cart_items
		WHERE user_id = $1 AND id = $2
	`, userID, itemID)

	if err != nil {
		return fmt.Errorf("failed to remove item from cart: %w", err)
	}
	return nil
}

// GetCartItems retrieves items in the user's cart
func (cr *CartRepo) GetCartItems(userID int) ([]models.CartItem, error) {
	rows, err := cr.DB.Query(`
		SELECT ci.id, ci.quantity, f.name, f.price
		FROM cart_items ci
		JOIN foods f ON ci.food_id = f.id
		WHERE ci.user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve cart items: %w", err)
	}
	defer rows.Close()

	var cartItems []models.CartItem
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(&item.ID, &item.Quantity, &item.FoodName, &item.FoodPrice)
		if err != nil {
			return nil, fmt.Errorf("error scanning cart item row: %w", err)
		}
		cartItems = append(cartItems, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return cartItems, nil
}
