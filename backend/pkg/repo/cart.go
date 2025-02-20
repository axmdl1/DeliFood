package repo

import (
	"DeliFood/backend/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartRepo struct {
	Col *mongo.Collection
}

func NewCartRepo(col *mongo.Collection) *CartRepo {
	return &CartRepo{Col: col}
}

func (repo *CartRepo) AddItemToCart(userID int, foodID int, quantity int, foodName string, foodPrice float64) error {
	cartItem := bson.M{
		"user_id":    userID,
		"food_id":    foodID,
		"quantity":   quantity,
		"food_name":  foodName,
		"food_price": foodPrice,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	_, err := repo.Col.InsertOne(context.Background(), cartItem)
	if err != nil {
		return fmt.Errorf("failed to add item to cart: %w", err)
	}
	return nil
}

func (cr *CartRepo) UpdateItemQuantity(userID, itemID, quantity int) error {
	filter := bson.M{"user_id": userID, "_id": itemID}
	update := bson.M{
		"$set": bson.M{
			"quantity":   quantity,
			"updated_at": time.Now(),
		},
	}
	_, err := cr.Col.UpdateOne(context.Background(), filter, update)
	return err
}

func (cr *CartRepo) RemoveItemFromCart(userID int, itemID int) error {
	filter := bson.M{"user_id": userID, "_id": itemID}
	_, err := cr.Col.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to remove item from cart: %w", err)
	}
	return nil
}

func (cr *CartRepo) GetCartItems(userID int) ([]models.CartItem, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := cr.Col.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve cart items: %w", err)
	}
	defer cursor.Close(context.Background())

	var cartItems []models.CartItem
	for cursor.Next(context.Background()) {
		var item models.CartItem
		if err := cursor.Decode(&item); err != nil {
			return nil, fmt.Errorf("error decoding cart item: %w", err)
		}
		cartItems = append(cartItems, item)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error during cursor iteration: %w", err)
	}
	return cartItems, nil
}
