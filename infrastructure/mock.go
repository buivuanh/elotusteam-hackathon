package infrastructure

import (
	"context"

	"github.com/buivuanh/elotusteam-hackathon/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

// Mock UserRepo implementation
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByID(ctx context.Context, db *pgxpool.Pool, userID int) (*domain.User, error) {
	args := m.Called(ctx, db, userID)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) GetByUserName(ctx context.Context, db *pgxpool.Pool, userName string) (*domain.User, error) {
	args := m.Called(ctx, db, userName)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) Insert(ctx context.Context, db *pgxpool.Pool, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, db, user)
	return args.Get(0).(*domain.User), args.Error(1)
}
