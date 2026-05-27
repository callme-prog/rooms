package repository

import (
	"context"

	"github.com/hotelapp/booking-service/internal/domain"
)

// BookingRepository defines the persistence contract for bookings.
//
//go:generate mockery --name=BookingRepository --output=../repository/mocks --outpkg=mocks
type BookingRepository interface {
	Create(ctx context.Context, b *domain.Booking) error
	GetByID(ctx context.Context, id string) (*domain.Booking, error)
	GetAll(ctx context.Context) ([]*domain.Booking, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Booking, error)
	Update(ctx context.Context, b *domain.Booking) error
	Delete(ctx context.Context, id string) error
	ChangeStatus(ctx context.Context, id string, status domain.BookingStatus) error
}
