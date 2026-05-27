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

type bookingRepo struct {
	db  *pgxpool.Pool
	log zerolog.Logger
}

func NewBookingRepo(db *pgxpool.Pool, log zerolog.Logger) repository.BookingRepository {
	return &bookingRepo{db: db, log: log}
}

const bookingCols = `b.id, b.user_id, b.room_id, b.check_in, b.check_out, b.status, b.comment, b.created_at,
	u.name as guest_name,
	r.name, r.description, r.price, r.capacity, r.image_url`

const bookingJoins = `
	FROM bookings b
	JOIN users    u ON u.id = b.user_id
	JOIN rooms    r ON r.id = b.room_id`

func scanBooking(row pgx.Row) (*domain.Booking, error) {
	b := &domain.Booking{Room: &domain.Room{}}
	err := row.Scan(
		&b.ID, &b.UserID, &b.RoomID,
		&b.CheckIn, &b.CheckOut, &b.Status, &b.Comment, &b.CreatedAt,
		&b.GuestName,
		&b.Room.Name, &b.Room.Description, &b.Room.Price, &b.Room.Capacity, &b.Room.ImageURL,
	)
	if err != nil {
		return nil, err
	}
	b.Room.ID = b.RoomID
	return b, nil
}

func (r *bookingRepo) Create(ctx context.Context, b *domain.Booking) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO bookings (id, user_id, room_id, check_in, check_out, status, comment, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		b.ID, b.UserID, b.RoomID, b.CheckIn, b.CheckOut, b.Status, b.Comment, b.CreatedAt)
	return err
}

func (r *bookingRepo) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+bookingCols+bookingJoins+` WHERE b.id=$1`, id)
	b, err := scanBooking(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, err
	}
	return b, nil
}

func (r *bookingRepo) GetAll(ctx context.Context) ([]*domain.Booking, error) {
	return r.queryBookings(ctx, `SELECT `+bookingCols+bookingJoins+` ORDER BY b.created_at DESC`)
}

func (r *bookingRepo) GetByUserID(ctx context.Context, userID string) ([]*domain.Booking, error) {
	return r.queryBookings(ctx,
		`SELECT `+bookingCols+bookingJoins+` WHERE b.user_id=$1 ORDER BY b.created_at DESC`, userID)
}

func (r *bookingRepo) queryBookings(ctx context.Context, query string, args ...interface{}) ([]*domain.Booking, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		b, err := scanBooking(rows)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *bookingRepo) Update(ctx context.Context, b *domain.Booking) error {
	_, err := r.db.Exec(ctx,
		`UPDATE bookings SET check_in=$2, check_out=$3, comment=$4, status=$5 WHERE id=$1`,
		b.ID, b.CheckIn, b.CheckOut, b.Comment, b.Status)
	return err
}

func (r *bookingRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM bookings WHERE id=$1`, id)
	return err
}

func (r *bookingRepo) ChangeStatus(ctx context.Context, id string, status domain.BookingStatus) error {
	_, err := r.db.Exec(ctx, `UPDATE bookings SET status=$2 WHERE id=$1`, id, status)
	return err
}
