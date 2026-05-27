package usecase_test

import (
	"context"
	"testing"

	"github.com/hotelapp/booking-service/internal/domain"
	"github.com/hotelapp/booking-service/internal/repository/mocks"
	"github.com/hotelapp/booking-service/internal/usecase"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoomUsecase_GetAll(t *testing.T) {
	repo := &mocks.RoomRepository{}
	uc := usecase.NewRoomUsecase(repo, zerolog.Nop())

	expected := []*domain.Room{
		{ID: "1", Name: "Suite", Price: 150.0},
		{ID: "2", Name: "Standard", Price: 80.0},
	}
	repo.On("GetAll", mock.Anything).Return(expected, nil)

	rooms, err := uc.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, rooms, 2)
	assert.Equal(t, "Suite", rooms[0].Name)
	repo.AssertExpectations(t)
}

func TestRoomUsecase_GetByID_Found(t *testing.T) {
	repo := &mocks.RoomRepository{}
	uc := usecase.NewRoomUsecase(repo, zerolog.Nop())

	room := &domain.Room{ID: "abc", Name: "Deluxe", Price: 200}
	repo.On("GetByID", mock.Anything, "abc").Return(room, nil)

	got, err := uc.GetByID(context.Background(), "abc")

	assert.NoError(t, err)
	assert.Equal(t, "Deluxe", got.Name)
	repo.AssertExpectations(t)
}

func TestRoomUsecase_GetByID_NotFound(t *testing.T) {
	repo := &mocks.RoomRepository{}
	uc := usecase.NewRoomUsecase(repo, zerolog.Nop())

	repo.On("GetByID", mock.Anything, "missing").Return(nil, domain.ErrRoomNotFound)

	_, err := uc.GetByID(context.Background(), "missing")

	assert.ErrorIs(t, err, domain.ErrRoomNotFound)
}
