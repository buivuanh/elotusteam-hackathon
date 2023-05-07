package infrastructure

import (
	"context"

	"github.com/buivuanh/elotusteam-hackathon/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo interface {
	GetByID(ctx context.Context, db *pgxpool.Pool, userID int) (*domain.User, error)
	GetByUserName(ctx context.Context, db *pgxpool.Pool, userName string) (*domain.User, error)
	Insert(ctx context.Context, db *pgxpool.Pool, user *domain.User) (*domain.User, error)
}

type FileInfo interface {
	InsertImage(ctx context.Context, db *pgxpool.Pool, file *domain.Image) (*domain.Image, error)
}
