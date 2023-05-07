package postgresql

import (
	"context"
	"fmt"

	"github.com/buivuanh/elotusteam-hackathon/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) GetByID(ctx context.Context, db *pgxpool.Pool, userID int) (*domain.User, error) {
	query := `
		SELECT id, username, hashed_password, created_at, deleted_at
		FROM users
		WHERE id = $1
	`
	var user domain.User
	err := db.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.UserName,
		&user.HashedPassword,
		&user.CreatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (r *UserRepo) GetByUserName(ctx context.Context, db *pgxpool.Pool, userName string) (*domain.User, error) {
	query := `
		SELECT id, username, hashed_password, created_at, deleted_at
		FROM users
		WHERE username = $1
	`
	var user domain.User
	err := db.QueryRow(ctx, query, userName).Scan(
		&user.UserID,
		&user.UserName,
		&user.HashedPassword,
		&user.CreatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

func (r *UserRepo) Insert(ctx context.Context, db *pgxpool.Pool, user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (username, hashed_password, created_at)
		VALUES ($1, $2, now())
		RETURNING id
	`
	err := db.QueryRow(ctx, query, user.UserName, user.HashedPassword).Scan(&user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	return user, nil
}
