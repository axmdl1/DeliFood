package repo

import (
	"DeliFood/backend/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FoodRepo struct {
	Col *mongo.Collection
}

func NewFoodRepo(col *mongo.Collection) *FoodRepo {
	return &FoodRepo{Col: col}
}

func (fr *FoodRepo) DeleteFood(id primitive.ObjectID) error {
	_, err := fr.Col.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete food: %w", err)
	}
	return nil
}

func (fr *FoodRepo) GetFood(category, sortParam string) ([]models.Food, error) {
	var foods []models.Food
	filter := bson.M{}
	if category != "" {
		filter["category"] = category
	}

	opts := options.Find()
	if sortParam != "" {
		switch sortParam {
		case "price-asc":
			opts.SetSort(bson.D{{Key: "price", Value: 1}})
		case "price-desc":
			opts.SetSort(bson.D{{Key: "price", Value: -1}})
		case "name":
			opts.SetSort(bson.D{{Key: "name", Value: 1}})
		}
	}

	cursor, err := fr.Col.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve food items: %w", err)
	}
	defer cursor.Close(context.Background())

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

func (fr *FoodRepo) GetFoodByID(foodID primitive.ObjectID) (*models.Food, error) {
	var food models.Food
	err := fr.Col.FindOne(context.Background(), bson.M{"_id": foodID}).Decode(&food)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no food found with ID: %v", foodID)
		}
		return nil, fmt.Errorf("error fetching food by ID: %w", err)
	}
	return &food, nil
}

func (fr *FoodRepo) GetFoodCategory() ([]string, error) {
	categories, err := fr.Col.Distinct(context.Background(), "category", bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve food categories: %w", err)
	}

	var cats []string
	for _, cat := range categories {
		if s, ok := cat.(string); ok {
			cats = append(cats, s)
		}
	}

	return cats, nil
}
