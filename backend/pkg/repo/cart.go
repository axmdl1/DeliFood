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

// AddItemToCart adds an item to the user's cart.
func (cr *CartRepo) AddItemToCart(userID, foodID, quantity int) error {
	// Check if item already exists in the cart
	var existingQuantity int
	err := cr.DB.QueryRow("SELECT quantity FROM cart_items WHERE user_id = $1 AND food_id = $2", userID, foodID).Scan(&existingQuantity)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking cart: %w", err)
	}

	if existingQuantity > 0 {
		// If item exists, update the quantity
		_, err := cr.DB.Exec("UPDATE cart_items SET quantity = quantity + $1, updated_at = CURRENT_TIMESTAMP WHERE user_id = $2 AND food_id = $3", quantity, userID, foodID)
		return err
	}

	// If item doesn't exist, insert it into the cart
	_, err = cr.DB.Exec("INSERT INTO cart_items (user_id, food_id, quantity) VALUES ($1, $2, $3)", userID, foodID, quantity)
	return err
}

// UpdateItemQuantity updates the quantity of an item in the cart
func (cr *CartRepo) UpdateItemQuantity(userID, foodID, quantity int) error {
	_, err := cr.DB.Exec("UPDATE cart_items SET quantity = $1, updated_at = CURRENT_TIMESTAMP WHERE user_id = $2 AND food_id = $3", quantity, userID, foodID)
	return err
}

// RemoveItemFromCart removes an item from the user's cart
func (cr *CartRepo) RemoveItemFromCart(userID, foodID int) error {
	_, err := cr.DB.Exec("DELETE FROM cart_items WHERE user_id = $1 AND food_id = $2", userID, foodID)
	return err
}

// GetCartItems retrieves all items in the user's cart
func (cr *CartRepo) GetCartItems(userID int) ([]models.CartItem, error) {
	rows, err := cr.DB.Query(`
		SELECT ci.id, f.name, f.price, ci.quantity, f.image
		FROM cart_items ci
		JOIN foods f ON ci.food_id = f.id
		WHERE ci.user_id = $1
	`, userID)

	if err != nil {
		return nil, fmt.Errorf("error retrieving cart items: %w", err)
	}
	defer rows.Close()

	var cartItems []models.CartItem
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Quantity, &item.Image)
		if err != nil {
			return nil, fmt.Errorf("error scanning cart item: %w", err)
		}
		cartItems = append(cartItems, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return cartItems, nil
}
