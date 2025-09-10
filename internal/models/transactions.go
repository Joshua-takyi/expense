package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id       uuid.UUID `db:"id" json:"id"`
	Amount   float64   `db:"amount" json:"amount"`
	Type     string    `db:"type" json:"type"` // "income" or "expense"
	Category string    `db:"category" json:"category"`
	UserId   uuid.UUID `db:"user_id" json:"user_id"`
	Created  time.Time `db:"created" json:"created"`
	Updated  time.Time `db:"updated" json:"updated"`
}

type TransactionService interface {
	AddTransaction(ctx context.Context, tx *Transaction, userID uuid.UUID) error
	UpdateTransaction(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	RemoveTransaction(ctx context.Context, id uuid.UUID) error
	GetTransactionDetails(ctx context.Context, id uuid.UUID) (*Transaction, error)
	ListUserTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Transaction, error)
}

func (r *Repository) AddTransaction(ctx context.Context, tx *Transaction, userID uuid.UUID) error {
	query := "INSERT INTO transactions (id, amount, type, category, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	now := time.Now()
	tx.Created = now
	tx.Updated = now
	result, err := r.DB.ExecContext(ctx, query, tx.Id, tx.Amount, tx.Type, tx.Category, userID, tx.Created, tx.Updated)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, transaction not added")
	}
	return nil
}
func (r *Repository) UpdateTransaction(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
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
	query := fmt.Sprintf("UPDATE transactions SET %s WHERE id=$%d", strings.Join(setClauses, ", "), i)
	args = append(args, id)
	_, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) RemoveTransaction(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM transactions WHERE id=$1"
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, transaction not found")
	}
	return nil
}
func (r *Repository) GetTransactionDetails(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	query := `SELECT id, amount, type, category, user_id, created_at, updated_at FROM transactions WHERE id=$1`
	row := r.DB.QueryRowContext(ctx, query, id)
	var tx Transaction
	err := row.Scan(&tx.Id, &tx.Amount, &tx.Type, &tx.Category, &tx.UserId, &tx.Created, &tx.Updated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}
	return &tx, nil
}
func (r *Repository) ListUserTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Transaction, error) {
	query := `SELECT id, amount, type, category, user_id, created_at, updated_at FROM transactions WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	transactions := []Transaction{}
	rows, err := r.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.Id, &tx.Amount, &tx.Type, &tx.Category, &tx.UserId, &tx.Created, &tx.Updated); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}
