package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Transaction struct {
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	Amount   float64            `bson:"amount" json:"amount"`
	Type     string             `bson:"type" json:"type"` // "income" or "expense"
	Category string             `bson:"category" json:"category"`
	UserId   primitive.ObjectID `bson:"user_id" json:"user_id"`
	Created  time.Time          `bson:"created" json:"created"`
	Updated  time.Time          `bson:"updated" json:"updated"`
}

type TransactionService interface {
	AddTransaction(ctx context.Context, tx *Transaction, userID primitive.ObjectID) error
	UpdateTransaction(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	RemoveTransaction(ctx context.Context, id primitive.ObjectID) error
	GetTransactionDetails(ctx context.Context, id primitive.ObjectID) (*Transaction, error)
	ListUserTransactions(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]Transaction, error)
}

func (r *Repository) AddTransaction(ctx context.Context, tx *Transaction, userID primitive.ObjectID) error {
	if r.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	tx.Id = primitive.NewObjectID()
	tx.UserId = userID
	tx.Created = time.Now()
	tx.Updated = time.Now()

	collection := r.DB.Database("expensetracker").Collection("transactions")
	_, err := collection.InsertOne(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to add transaction: %v", err)
	}
	return nil
}
func (r *Repository) UpdateTransaction(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	setClauses := []string{}
	args := []interface{}{}
	i := 1
	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s=$%d", key, i))
		args = append(args, value)
		i++
	}
	if len(setClauses) == 0 {
		return fmt.Errorf("no valid updates provided")
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at=$%d", i))
	collection := r.DB.Database("expensetracker").Collection("transactions")
	args = append(args, time.Now())
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}
	update["$set"].(bson.M)["updated_at"] = time.Now()
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %v", err)
	}
	return nil
}

func (r *Repository) RemoveTransaction(ctx context.Context, id primitive.ObjectID) error {
	if r.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	filter := bson.M{"_id": id}
	collection := r.DB.Database("expensetracker").Collection("transactions")
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to remove transaction: %v", err)
	}
	return nil
}

func (r *Repository) GetTransactionDetails(ctx context.Context, id primitive.ObjectID) (*Transaction, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	collection := r.DB.Database("expensetracker").Collection("transactions")
	var tx Transaction
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&tx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}
	return &tx, nil
}

func (r *Repository) ListUserTransactions(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]Transaction, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	filter := bson.M{"user_id": userID}
	options := options.Find()
	if limit > 0 {
		options.SetLimit(int64(limit))
	}
	if offset > 0 {
		options.SetSkip(int64(offset))
	}
	options.SetSort(bson.D{{Key: "created", Value: -1}}) // Sort by created date descending
	collection := r.DB.Database("expensetracker").Collection("transactions")
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %v", err)
	}
	defer cursor.Close(ctx)

	var transactions []Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, fmt.Errorf("failed to decode transactions: %v", err)
	}

	return transactions, nil
}
