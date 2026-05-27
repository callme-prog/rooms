package usecase

import (
	"context"

	"github.com/hotelapp/booking-service/internal/domain"
	"github.com/hotelapp/booking-service/internal/repository"
	"github.com/rs/zerolog"
)

// RoomUsecase defines business logic for rooms.
type RoomUsecase interface {
	GetAll(ctx context.Context) ([]*domain.Room, error)
	GetByID(ctx context.Context, id string) (*domain.Room, error)
}

type roomUsecase struct {
	roomRepo repository.RoomRepository
	log      zerolog.Logger
}

func NewRoomUsecase(roomRepo repository.RoomRepository, log zerolog.Logger) RoomUsecase {
	return &roomUsecase{roomRepo: roomRepo, log: log}
}

func (uc *roomUsecase) GetAll(ctx context.Context) ([]*domain.Room, error) {
	return uc.roomRepo.GetAll(ctx)
}

func (uc *roomUsecase) GetByID(ctx context.Context, id string) (*domain.Room, error) {
	return uc.roomRepo.GetByID(ctx, id)
}
