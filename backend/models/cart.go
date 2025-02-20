package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// CartItem represents an item in the cart
type CartItem struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	FoodID     primitive.ObjectID `bson:"food_id" json:"food_id"`
	FoodName   string             `bson:"food_name" json:"food_name"`
	FoodPrice  float64            `bson:"food_price" json:"food_price"`
	Quantity   int                `bson:"quantity" json:"quantity"`
	TotalPrice float64            `bson:"total_price" json:"total_price"`
}
