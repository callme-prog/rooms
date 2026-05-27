package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	commonjwt "github.com/hotelapp/common/pkg/jwt"
	"github.com/hotelapp/auth-service/internal/domain"
	"github.com/hotelapp/auth-service/internal/repository"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

// AuthUsecase defines business logic for authentication.
type AuthUsecase interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.TokenResponse, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenResponse, error)
	Me(ctx context.Context, userID string) (*domain.User, error)
}

type authUsecase struct {
	userRepo   repository.UserRepository
	jwtManager *commonjwt.Manager
	log        zerolog.Logger
}

// NewAuthUsecase constructs an AuthUsecase.
func NewAuthUsecase(
	userRepo repository.UserRepository,
	jwtManager *commonjwt.Manager,
	log zerolog.Logger,
) AuthUsecase {
	return &authUsecase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		log:        log,
	}
}

func (uc *authUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.TokenResponse, error) {
	// Check uniqueness
	_, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, domain.ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           uuid.NewString(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         domain.RoleGuest,
		CreatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.log.Error().Err(err).Msg("failed to create user")
		return nil, err
	}

	token, err := uc.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	uc.log.Info().Str("user_id", user.ID).Msg("user registered")
	return &domain.TokenResponse{AccessToken: token, User: user}, nil
}

func (uc *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := uc.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	uc.log.Info().Str("user_id", user.ID).Msg("user logged in")
	return &domain.TokenResponse{AccessToken: token, User: user}, nil
}

func (uc *authUsecase) Me(ctx context.Context, userID string) (*domain.User, error) {
	return uc.userRepo.GetByID(ctx, userID)
}
