package repo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"DeliFood/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepo struct {
	Foods *mongo.Collection
	Users *mongo.Collection
}

func NewAdminRepo(foods, users *mongo.Collection) *AdminRepo {
	return &AdminRepo{
		Foods: foods,
		Users: users,
	}
}

func (ar *AdminRepo) GetUserRole(userID string) (string, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", err
	}

	var user models.User
	err = ar.Users.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return "", err
	}

	return user.Role, nil
}

func (ar *AdminRepo) AddFood(food *models.Food) error {
	food.CreatedAt = time.Now()
	food.UpdatedAt = time.Now()
	_, err := ar.Foods.InsertOne(context.Background(), food)
	if err != nil {
		return fmt.Errorf("Error inserting food: %w", err)
	}
	return nil
}

func (ar *AdminRepo) UpdateFood(food models.Food) error {
	filter := bson.M{"_id": food.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        food.Name,
			"category":    food.Category,
			"image":       food.Image,
			"description": food.Description,
			"price":       food.Price,
			"updated_at":  time.Now(),
		},
	}
	_, err := ar.Foods.UpdateOne(context.Background(), filter, update)
	return err
}

func (ar *AdminRepo) UpdateUserRole(userID int, role string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"role": role}}
	_, err := ar.Users.UpdateOne(context.Background(), filter, update)
	return err
}

func (ar *AdminRepo) GetAllFoods() ([]models.Food, error) {
	cursor, err := ar.Foods.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error fetching foods: %w", err)
	}
	defer cursor.Close(context.Background())

	var foods []models.Food
	for cursor.Next(context.Background()) {
		var food models.Food
		if err := cursor.Decode(&food); err != nil {
			return nil, fmt.Errorf("error decoding food: %w", err)
		}
		foods = append(foods, food)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error during cursor iteration: %w", err)
	}

	return foods, nil
}
