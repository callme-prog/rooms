package repository

import (
	"context"

	"github.com/hotelapp/booking-service/internal/domain"
)

// RoomRepository defines read-only persistence for rooms.
//
//go:generate mockery --name=RoomRepository --output=../repository/mocks --outpkg=mocks
type RoomRepository interface {
	GetAll(ctx context.Context) ([]*domain.Room, error)
	GetByID(ctx context.Context, id string) (*domain.Room, error)
}
