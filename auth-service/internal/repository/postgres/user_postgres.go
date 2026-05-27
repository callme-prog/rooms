package postgres

import (
	"context"
	"errors"

	"github.com/hotelapp/auth-service/internal/domain"
	"github.com/hotelapp/auth-service/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type userRepo struct {
	db  *pgxpool.Pool
	log zerolog.Logger
}

// NewUserRepo constructs a Postgres-backed UserRepository.
func NewUserRepo(db *pgxpool.Pool, log zerolog.Logger) repository.UserRepository {
	return &userRepo{db: db, log: log}
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, u.ID, u.Name, u.Email, u.PasswordHash, u.Role, u.CreatedAt)
	if err != nil {
		r.log.Error().Err(err).Str("email", u.Email).Msg("create user failed")
		return err
	}
	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, role, created_at FROM users WHERE email=$1`
	row := r.db.QueryRow(ctx, query, email)
	u := &domain.User{}
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, role, created_at FROM users WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	u := &domain.User{}
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}
