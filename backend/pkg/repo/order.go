package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepo struct {
	Col *mongo.Collection
}

func NewOrderRepo(col *mongo.Collection) *OrderRepo {
	return &OrderRepo{Col: col}
}

func (r *OrderRepo) CreateOrder(userID int, totalPrice float64) (interface{}, error) {
	order := bson.M{
		"user_id":     userID,
		"total_price": totalPrice,
		"status":      "pending",
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}
	res, err := r.Col.InsertOne(context.Background(), order)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}
	return res.InsertedID, nil
}

func (r *OrderRepo) UpdateOrderStatus(orderID interface{}, newStatus string) error {
	filter := bson.M{"_id": orderID}
	update := bson.M{
		"$set": bson.M{
			"status":     newStatus,
			"updated_at": time.Now(),
		},
	}
	res, err := r.Col.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	if res.MatchedCount == 0 {
		return errors.New("order not found")
	}
	return nil
}
