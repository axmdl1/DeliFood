package repo

import (
	"DeliFood/backend/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	DB *mongo.Collection
}

func NewUserRepo(db *mongo.Collection) *UserRepo {
	return &UserRepo{DB: db}
}

func (ur *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.DB.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no user found with email: %s", email)
		}
		return nil, fmt.Errorf("error fetching user by email: %w", err)
	}
	return &user, nil
}

func (ur *UserRepo) Register(user models.User) error {
	count, err := ur.DB.CountDocuments(context.Background(), bson.M{"$or": []bson.M{
		{"email": user.Email},
		{"username": user.UserName},
	}})
	if err != nil {
		return fmt.Errorf("failed to check for existing user: %w", err)
	}
	if count > 0 {
		return errors.New("email or username already exists")
	}
	_, err = ur.DB.InsertOne(context.Background(), user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (ur *UserRepo) CheckEmailOrUsernameExists(email, username string) (bool, error) {
	count, err := ur.DB.CountDocuments(context.Background(), bson.M{"$or": []bson.M{
		{"email": email},
		{"username": username},
	}})
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}

func (ur *UserRepo) VerifyEmail(email, code string) error {
	var user models.User
	err := ur.DB.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	if user.VerificationCode != code {
		return errors.New("invalid verification code")
	}
	_, err = ur.DB.UpdateOne(
		context.Background(),
		bson.M{"email": email},
		bson.M{"$set": bson.M{"is_verified": true}},
	)
	if err != nil {
		return fmt.Errorf("failed to update user verification status: %w", err)
	}
	return nil
}

func (ur *UserRepo) Authenticate(email, password string) (models.User, error) {
	var user models.User
	err := ur.DB.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, errors.New("invalid email or password")
		}
		return models.User{}, fmt.Errorf("error fetching user: %w", err)
	}
	if !user.IsVerified {
		return models.User{}, errors.New("email not verified")
	}
	// Password comparison logic goes here
	return user, nil
}
