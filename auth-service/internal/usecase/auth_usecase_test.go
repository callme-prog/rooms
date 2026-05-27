package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	commonjwt "github.com/hotelapp/common/pkg/jwt"
	"github.com/hotelapp/auth-service/internal/domain"
	"github.com/hotelapp/auth-service/internal/repository/mocks"
	"github.com/hotelapp/auth-service/internal/usecase"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func newTestJWTManager() *commonjwt.Manager {
	return commonjwt.NewManager("test-secret", 15*time.Minute)
}

func TestRegister_Success(t *testing.T) {
	repo := &mocks.UserRepository{}
	uc := usecase.NewAuthUsecase(repo, newTestJWTManager(), zerolog.Nop())

	repo.On("GetByEmail", mock.Anything, "john@example.com").
		Return(nil, domain.ErrUserNotFound)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
		Return(nil)

	resp, err := uc.Register(context.Background(), &domain.RegisterRequest{
		Name:     "John",
		Email:    "john@example.com",
		Password: "secret123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.Equal(t, domain.RoleGuest, resp.User.Role)
	repo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := &mocks.UserRepository{}
	uc := usecase.NewAuthUsecase(repo, newTestJWTManager(), zerolog.Nop())

	existing := &domain.User{ID: uuid.NewString(), Email: "dup@example.com"}
	repo.On("GetByEmail", mock.Anything, "dup@example.com").Return(existing, nil)

	_, err := uc.Register(context.Background(), &domain.RegisterRequest{
		Name: "Dup", Email: "dup@example.com", Password: "secret123",
	})

	assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
	repo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	repo := &mocks.UserRepository{}
	uc := usecase.NewAuthUsecase(repo, newTestJWTManager(), zerolog.Nop())

	hash, _ := bcrypt.GenerateFromPassword([]byte("mypassword"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:           uuid.NewString(),
		Email:        "alice@example.com",
		PasswordHash: string(hash),
		Role:         domain.RoleGuest,
		CreatedAt:    time.Now(),
	}
	repo.On("GetByEmail", mock.Anything, "alice@example.com").Return(user, nil)

	resp, err := uc.Login(context.Background(), &domain.LoginRequest{
		Email:    "alice@example.com",
		Password: "mypassword",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	repo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := &mocks.UserRepository{}
	uc := usecase.NewAuthUsecase(repo, newTestJWTManager(), zerolog.Nop())

	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	user := &domain.User{ID: uuid.NewString(), Email: "bob@example.com", PasswordHash: string(hash)}
	repo.On("GetByEmail", mock.Anything, "bob@example.com").Return(user, nil)

	_, err := uc.Login(context.Background(), &domain.LoginRequest{
		Email: "bob@example.com", Password: "wrongpass",
	})

	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mocks.UserRepository{}
	uc := usecase.NewAuthUsecase(repo, newTestJWTManager(), zerolog.Nop())

	repo.On("GetByEmail", mock.Anything, "ghost@example.com").Return(nil, domain.ErrUserNotFound)

	_, err := uc.Login(context.Background(), &domain.LoginRequest{
		Email: "ghost@example.com", Password: "pass",
	})

	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestMe_Success(t *testing.T) {
	repo := &mocks.UserRepository{}
	uc := usecase.NewAuthUsecase(repo, newTestJWTManager(), zerolog.Nop())

	uid := uuid.NewString()
	expected := &domain.User{ID: uid, Name: "Test", Email: "t@t.com", Role: domain.RoleGuest}
	repo.On("GetByID", mock.Anything, uid).Return(expected, nil)

	got, err := uc.Me(context.Background(), uid)
	assert.NoError(t, err)
	assert.Equal(t, expected.Email, got.Email)
}
