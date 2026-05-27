package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hotelapp/booking-service/internal/domain"
	"github.com/hotelapp/booking-service/internal/repository/mocks"
	"github.com/hotelapp/booking-service/internal/usecase"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makeBooking(userID string) *domain.Booking {
	return &domain.Booking{
		ID:       uuid.NewString(),
		UserID:   userID,
		RoomID:   uuid.NewString(),
		Room:     &domain.Room{ID: uuid.NewString(), Name: "Suite"},
		CheckIn:  time.Now(),
		CheckOut: time.Now().Add(48 * time.Hour),
		Status:   domain.StatusConfirmed,
	}
}

func TestBookingUsecase_Create(t *testing.T) {
	repo := &mocks.BookingRepository{}
	uc := usecase.NewBookingUsecase(repo, zerolog.Nop())

	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Booking")).Return(nil)

	b, err := uc.Create(context.Background(), "user-1", &domain.CreateBookingRequest{
		RoomID:   "room-1",
		CheckIn:  time.Now(),
		CheckOut: time.Now().Add(24 * time.Hour),
	})

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusPending, b.Status)
	assert.NotEmpty(t, b.ID)
	repo.AssertExpectations(t)
}

func TestBookingUsecase_Update_Forbidden(t *testing.T) {
	repo := &mocks.BookingRepository{}
	uc := usecase.NewBookingUsecase(repo, zerolog.Nop())

	b := makeBooking("owner-id")
	repo.On("GetByID", mock.Anything, b.ID).Return(b, nil)

	_, err := uc.Update(context.Background(), "other-user", b.ID, &domain.UpdateBookingRequest{
		CheckIn:  time.Now(),
		CheckOut: time.Now().Add(24 * time.Hour),
	})

	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestBookingUsecase_Update_SetsPendingEdit(t *testing.T) {
	repo := &mocks.BookingRepository{}
	uc := usecase.NewBookingUsecase(repo, zerolog.Nop())

	b := makeBooking("user-1")
	updated := *b
	updated.Status = domain.StatusPendingEdit

	repo.On("GetByID", mock.Anything, b.ID).Return(b, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Booking")).Return(nil)

	got, err := uc.Update(context.Background(), "user-1", b.ID, &domain.UpdateBookingRequest{
		CheckIn:  time.Now().Add(2 * time.Hour),
		CheckOut: time.Now().Add(50 * time.Hour),
	})

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusPendingEdit, got.Status)
	repo.AssertExpectations(t)
}

func TestBookingUsecase_ChangeStatus_Confirm(t *testing.T) {
	repo := &mocks.BookingRepository{}
	uc := usecase.NewBookingUsecase(repo, zerolog.Nop())

	b := makeBooking("user-1")
	b.Status = domain.StatusPendingEdit
	confirmed := *b
	confirmed.Status = domain.StatusConfirmed

	repo.On("ChangeStatus", mock.Anything, b.ID, domain.StatusConfirmed).Return(nil)
	repo.On("GetByID", mock.Anything, b.ID).Return(&confirmed, nil)

	got, err := uc.ChangeStatus(context.Background(), b.ID, domain.StatusConfirmed)

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusConfirmed, got.Status)
	repo.AssertExpectations(t)
}

func TestBookingUsecase_Delete_AdminBypass(t *testing.T) {
	repo := &mocks.BookingRepository{}
	uc := usecase.NewBookingUsecase(repo, zerolog.Nop())

	repo.On("Delete", mock.Anything, "booking-1").Return(nil)

	err := uc.Delete(context.Background(), "any-user", "booking-1", true)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBookingUsecase_GetAll(t *testing.T) {
	repo := &mocks.BookingRepository{}
	uc := usecase.NewBookingUsecase(repo, zerolog.Nop())

	expected := []*domain.Booking{makeBooking("u1"), makeBooking("u2")}
	repo.On("GetAll", mock.Anything).Return(expected, nil)

	got, err := uc.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, got, 2)
	repo.AssertExpectations(t)
}
