package models

// CartItem represents an item in the cart
type CartItem struct {
	ID         int
	UserID     int
	FoodID     int
	FoodName   string
	FoodPrice  float64
	Quantity   int
	TotalPrice float64
}
