package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Joshua-takyi/expense/server/internal/helpers"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Password  string             `bson:"password" json:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

var validate = validator.New()

type UserService interface {
	RegisterUser(ctx context.Context, user *User) error
	AuthenticateUser(ctx context.Context, email, password string) (*User, error)
	UpdateUserProfile(ctx context.Context, id primitive.ObjectID, user *User) error
	DeleteUserAccount(ctx context.Context, id primitive.ObjectID) error
	GetUserProfile(ctx context.Context, id primitive.ObjectID) (*User, error)
}

type Repository struct {
	DB *mongo.Client
}

type Service interface {
	UserService
	TransactionService
}

func (r *Repository) checkUserExists(ctx context.Context, email string) (bool, error) {
	var exists bool

	filter := bson.M{"email": email}

	count, err := r.DB.Database("expensetracker").Collection("users").CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %w", err)
	}
	exists = count > 0
	return exists, nil
}

func (r *Repository) RegisterUser(ctx context.Context, user *User) error {
	if r.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	user.Password = strings.TrimSpace(user.Password)

	// ok := helpers.IsStrongPassword(user.Password)
	// if !ok {
	// 	return fmt.Errorf("password is not strong enough")
	// }

	exists, err := r.checkUserExists(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	if err := validate.Struct(user); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	now := time.Now()
	user.Id = primitive.NewObjectID()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.Password = hashedPassword

	_, err = r.DB.Database("expensetracker").Collection("users").InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}

	return nil
}

func (r *Repository) AuthenticateUser(ctx context.Context, email, password string) (*User, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	password = strings.TrimSpace(password)

	user := User{}
	filter := bson.M{"email": email}
	err := r.DB.Database("expensetracker").Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	ok := helpers.CheckPasswordHash(password, user.Password)
	if !ok {
		return nil, fmt.Errorf("invalid password")
	}

	// Remove password before returning user
	user.Password = ""
	return &user, nil

}

func (r *Repository) UpdateUserProfile(ctx context.Context, id primitive.ObjectID, user *User) error {

	if r.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": user}
	update["$set"].(bson.M)["updated_at"] = time.Now()

	result, err := r.DB.Database("expensetracker").Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating user profile: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no user found with id %s", id)
	}
	return nil
}

func (r *Repository) DeleteUserAccount(ctx context.Context, id primitive.ObjectID) error {
	if r.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	filter := bson.M{"_id": id}
	result, err := r.DB.Database("expensetracker").Collection("users").DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no user found with id %s", id)
	}
	return nil
}

func (r *Repository) GetUserProfile(ctx context.Context, id primitive.ObjectID) (*User, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	user := User{}
	filter := bson.M{"_id": id}
	err := r.DB.Database("expensetracker").Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user with id %s not found", id)
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	// Remove password before returning user
	user.Password = ""
	return &user, nil
}
