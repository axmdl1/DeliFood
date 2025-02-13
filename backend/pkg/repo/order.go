package repo

import (
	"database/sql"
	"fmt"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) CreateOrder(userID int, totalPrice float64) (int, error) {
	var orderID int
	err := r.db.QueryRow(`
        INSERT INTO orders (user_id, total_price)
        VALUES ($1, $2)
        RETURNING id
    `, userID, totalPrice).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert order: %w", err)
	}
	return orderID, nil
}

func (r *OrderRepo) UpdateOrderStatus(orderID int, newStatus string) error {
	_, err := r.db.Exec(`
        UPDATE orders
        SET status = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
    `, newStatus, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}
