package postgres

import (
	"context"
	"errors"

	"github.com/hotelapp/booking-service/internal/domain"
	"github.com/hotelapp/booking-service/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type roomRepo struct {
	db  *pgxpool.Pool
	log zerolog.Logger
}

func NewRoomRepo(db *pgxpool.Pool, log zerolog.Logger) repository.RoomRepository {
	return &roomRepo{db: db, log: log}
}

func (r *roomRepo) GetAll(ctx context.Context) ([]*domain.Room, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, description, price, capacity, image_url FROM rooms ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*domain.Room
	for rows.Next() {
		m := &domain.Room{}
		if err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.Price, &m.Capacity, &m.ImageURL); err != nil {
			return nil, err
		}
		rooms = append(rooms, m)
	}
	return rooms, nil
}

func (r *roomRepo) GetByID(ctx context.Context, id string) (*domain.Room, error) {
	m := &domain.Room{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, description, price, capacity, image_url FROM rooms WHERE id=$1`, id).
		Scan(&m.ID, &m.Name, &m.Description, &m.Price, &m.Capacity, &m.ImageURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRoomNotFound
		}
		return nil, err
	}
	return m, nil
}
