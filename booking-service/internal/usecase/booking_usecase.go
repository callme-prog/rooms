package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hotelapp/booking-service/internal/domain"
	"github.com/hotelapp/booking-service/internal/repository"
	"github.com/rs/zerolog"
)

// BookingUsecase defines business logic for bookings.
type BookingUsecase interface {
	Create(ctx context.Context, userID string, req *domain.CreateBookingRequest) (*domain.Booking, error)
	GetAll(ctx context.Context) ([]*domain.Booking, error)
	GetByUser(ctx context.Context, userID string) ([]*domain.Booking, error)
	Update(ctx context.Context, userID, bookingID string, req *domain.UpdateBookingRequest) (*domain.Booking, error)
	Delete(ctx context.Context, userID, bookingID string, isAdmin bool) error
	ChangeStatus(ctx context.Context, bookingID string, status domain.BookingStatus) (*domain.Booking, error)
}

type bookingUsecase struct {
	bookingRepo repository.BookingRepository
	log         zerolog.Logger
}

func NewBookingUsecase(bookingRepo repository.BookingRepository, log zerolog.Logger) BookingUsecase {
	return &bookingUsecase{bookingRepo: bookingRepo, log: log}
}

func (uc *bookingUsecase) Create(ctx context.Context, userID string, req *domain.CreateBookingRequest) (*domain.Booking, error) {
	b := &domain.Booking{
		ID:        uuid.NewString(),
		UserID:    userID,
		RoomID:    req.RoomID,
		CheckIn:   req.CheckIn,
		CheckOut:  req.CheckOut,
		Status:    domain.StatusPending,
		Comment:   req.Comment,
		CreatedAt: time.Now(),
	}
	if err := uc.bookingRepo.Create(ctx, b); err != nil {
		uc.log.Error().Err(err).Msg("create booking failed")
		return nil, err
	}
	uc.log.Info().Str("booking_id", b.ID).Str("user_id", userID).Msg("booking created")
	return b, nil
}

func (uc *bookingUsecase) GetAll(ctx context.Context) ([]*domain.Booking, error) {
	return uc.bookingRepo.GetAll(ctx)
}

func (uc *bookingUsecase) GetByUser(ctx context.Context, userID string) ([]*domain.Booking, error) {
	return uc.bookingRepo.GetByUserID(ctx, userID)
}

func (uc *bookingUsecase) Update(ctx context.Context, userID, bookingID string, req *domain.UpdateBookingRequest) (*domain.Booking, error) {
	b, err := uc.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if b.UserID != userID {
		return nil, domain.ErrForbidden
	}
	b.CheckIn  = req.CheckIn
	b.CheckOut = req.CheckOut
	b.Comment  = req.Comment
	b.Status   = domain.StatusPendingEdit
	if err := uc.bookingRepo.Update(ctx, b); err != nil {
		return nil, err
	}
	uc.log.Info().Str("booking_id", bookingID).Msg("booking updated, pending edit")
	return b, nil
}

func (uc *bookingUsecase) Delete(ctx context.Context, userID, bookingID string, isAdmin bool) error {
	if !isAdmin {
		b, err := uc.bookingRepo.GetByID(ctx, bookingID)
		if err != nil { return err }
		if b.UserID != userID { return domain.ErrForbidden }
	}
	return uc.bookingRepo.Delete(ctx, bookingID)
}

func (uc *bookingUsecase) ChangeStatus(ctx context.Context, bookingID string, status domain.BookingStatus) (*domain.Booking, error) {
	if err := uc.bookingRepo.ChangeStatus(ctx, bookingID, status); err != nil {
		return nil, err
	}
	uc.log.Info().Str("booking_id", bookingID).Str("status", string(status)).Msg("status changed")
	return uc.bookingRepo.GetByID(ctx, bookingID)
}
