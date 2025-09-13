package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Transaction struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Amount      float64            `bson:"amount" json:"amount"`
	Type        string             `bson:"type" json:"type"`
	Description string             `bson:"description" json:"description"`
	Note        string             `bson:"note" json:"note"`
	Category    string             `bson:"category" json:"category"`
	UserId      primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type TransactionService interface {
	AddTransaction(ctx context.Context, tx *Transaction, userID primitive.ObjectID) error
	UpdateTransaction(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	RemoveTransaction(ctx context.Context, id primitive.ObjectID) error
	GetTransactionDetails(ctx context.Context, id primitive.ObjectID) (*Transaction, error)
	ListUserTransactions(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]Transaction, error)
	GetTransactionByQuery(ctx context.Context, userID primitive.ObjectID, query string, category []string, order string, limit, offset int) ([]Transaction, error)
}

func (r *Repository) AddTransaction(ctx context.Context, tx *Transaction, userID primitive.ObjectID) error {
	if r.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	if err := validate.Struct(tx); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	tx.UserId = userID
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()
	tx.Id = primitive.NewObjectID()

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

	// Add updated_at timestamp to the updates
	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	collection := r.DB.Database("expensetracker").Collection("transactions")
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
	options.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by created date descending
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

func (r *Repository) GetTransactionByQuery(ctx context.Context, userID primitive.ObjectID, query string, category []string, order string, limit, offset int) ([]Transaction, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	order = strings.ToLower(order)

	// Start with base filter for user
	filters := []bson.M{{"user_id": userID}}

	// Add search query filter if provided
	if query != "" {
		queryFilter := bson.M{
			"$or": []bson.M{
				{"description": bson.M{"$regex": query, "$options": "i"}},
				{"note": bson.M{"$regex": query, "$options": "i"}},
				{"type": bson.M{"$regex": query, "$options": "i"}},
				{"category": bson.M{"$regex": query, "$options": "i"}},
			},
		}
		filters = append(filters, queryFilter)
	}

	// Add category filter
	if len(category) > 0 {
		// Convert categories to lowercase
		var lowerCategories []string
		for _, cat := range category {
			lowerCat := strings.ToLower(strings.TrimSpace(cat))
			if lowerCat != "" && lowerCat != "all" {
				lowerCategories = append(lowerCategories, lowerCat)
			}
		}

		if len(lowerCategories) > 0 {
			if len(lowerCategories) == 1 {
				filters = append(filters, bson.M{"category": lowerCategories[0]})
			} else {
				filters = append(filters, bson.M{"category": bson.M{"$in": lowerCategories}})
			}
		}
	}

	// Build the amount filter
	// if amount != "" {
	// 	if strings.Contains(amount, "-") {
	// 		parts := strings.Split(amount, "-")
	// 		if len(parts) == 2 {
	// 			minAmountStr := strings.TrimSpace(parts[0])
	// 			maxAmountStr := strings.TrimSpace(parts[1])

	// 			minAmount, err1 := strconv.ParseFloat(minAmountStr, 64)
	// 			maxAmount, err2 := strconv.ParseFloat(maxAmountStr, 64)

	// 			if err1 == nil && err2 == nil {
	// 				filters = append(filters, bson.M{
	// 					"amount": bson.M{
	// 						"$gte": minAmount,
	// 						"$lte": maxAmount,
	// 					},
	// 				})
	// 			}
	// 		}
	// 	} else {
	// 		// Try to parse as exact amount
	// 		if exactAmount, err := strconv.ParseFloat(amount, 64); err == nil {
	// 			filters = append(filters, bson.M{"amount": exactAmount})
	// 		}
	// 	}
	// }

	// Build the sort option
	var sortOption bson.D
	switch order {
	case "new":
		sortOption = bson.D{{Key: "created_at", Value: -1}} // New to old (descending)
	case "old":
		sortOption = bson.D{{Key: "created_at", Value: 1}} // Old to new (ascending)
	case "asc":
		sortOption = bson.D{{Key: "amount", Value: 1}} // Amount: Low to High
	case "desc":
		sortOption = bson.D{{Key: "amount", Value: -1}} // Amount: High to Low
	default:
		sortOption = bson.D{{Key: "created_at", Value: -1}} // Default to new to old
	}

	// Combine all filters using $and
	var filter bson.M
	if len(filters) == 1 {
		filter = filters[0]
	} else {
		filter = bson.M{"$and": filters}
	}

	// Set up find options
	options := options.Find()
	if limit > 0 {
		options.SetLimit(int64(limit))
	}
	if offset > 0 {
		options.SetSkip(int64(offset))
	}
	options.SetSort(sortOption) // Use the calculated sort option

	collection := r.DB.Database("expensetracker").Collection("transactions")
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, fmt.Errorf("failed to search transactions: %v", err)
	}
	defer cursor.Close(ctx)

	var transactions []Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, fmt.Errorf("failed to decode transactions: %v", err)
	}

	return transactions, nil
}
