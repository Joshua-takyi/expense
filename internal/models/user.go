package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/joshua/expensetracker/internal/helpers"
)

type User struct {
	Id           uuid.UUID     `db:"id" json:"id"`
	Name         string        `db:"name" json:"name" validate:"required"`
	Email        string        `db:"email" json:"email" validate:"required,email"`
	Password     string        `db:"password" json:"-"`
	CreatedAt    time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at" json:"updated_at"`
	Transactions []Transaction `json:"transactions,omitempty"`
}

var validate = validator.New()

type UserService interface {
	RegisterUser(ctx context.Context, user *User) error
	AuthenticateUser(ctx context.Context, email, password string) (*User, error)
	UpdateUserProfile(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteUserAccount(ctx context.Context, id uuid.UUID) error
	GetUserProfile(ctx context.Context, id uuid.UUID) (*User, error)
}

type Repository struct {
	DB *sql.DB
}

type Service interface {
	UserService
	TransactionService
}

func (r *Repository) checkUserExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)"
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) RegisterUser(ctx context.Context, user *User) error {
	// check if user email already exists
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
	ok := helpers.IsStrongPassword(user.Password)
	if !ok {
		return fmt.Errorf("password is not strong enough")
	}

	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	query := "INSERT INTO users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err = r.DB.Exec(query, user.Name, user.Email, hashedPassword, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *Repository) AuthenticateUser(ctx context.Context, email, password string) (*User, error) {
	user := &User{}
	query := "SELECT id, name, email, password, created_at, updated_at FROM users WHERE email=$1"
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, err
	}

	if !helpers.CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("invalid password")
	}
	// Remove password before returning user
	user.Password = ""
	return user, nil
}

func (r *Repository) UpdateUserProfile(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {

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
	query := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d", strings.Join(setClauses, ", "), i)
	args = append(args, id)
	_, err := r.DB.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) DeleteUserAccount(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM users WHERE id=$1"
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error fetching rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %s", id)
	}
	return nil
}

func (r *Repository) GetUserProfile(ctx context.Context, id uuid.UUID) (*User, error) {
	user := &User{}
	query := "SELECT id, name, email, transactions, created_at, updated_at FROM users WHERE id=$1"
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Transactions, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %s not found", id)
		}
		return nil, err
	}
	return user, nil
}
