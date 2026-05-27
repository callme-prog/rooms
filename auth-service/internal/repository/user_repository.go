package repository

import (
	"context"

	"github.com/hotelapp/auth-service/internal/domain"
)

// UserRepository defines the persistence contract for users.
//
//go:generate mockery --name=UserRepository --output=../repository/mocks --outpkg=mocks
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}
