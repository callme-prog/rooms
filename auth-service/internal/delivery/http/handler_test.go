package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	delivery "github.com/hotelapp/auth-service/internal/delivery/http"
	"github.com/hotelapp/auth-service/internal/domain"
	commonjwt "github.com/hotelapp/common/pkg/jwt"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockAuthUsecase — ручной мок usecase.AuthUsecase
type mockAuthUsecase struct{ mock.Mock }

func (m *mockAuthUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.TokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenResponse), args.Error(1)
}

func (m *mockAuthUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenResponse), args.Error(1)
}

func (m *mockAuthUsecase) Me(ctx context.Context, userID string) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func setupRouter(uc *mockAuthUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	jm := commonjwt.NewManager("test-secret", 15*time.Minute)
	h := delivery.NewHandler(uc, jm, zerolog.Nop())
	r := gin.New()
	r.POST("/api/v1/auth/register", h.Register)
	r.POST("/api/v1/auth/login", h.Login)
	r.GET("/api/v1/auth/me", h.JWTMiddleware(), h.Me)
	return r
}

func TestHandlerRegister_Created(t *testing.T) {
	uc := &mockAuthUsecase{}
	r := setupRouter(uc)

	uc.On("Register", mock.Anything, mock.Anything).Return(&domain.TokenResponse{
		AccessToken: "tok",
		User:        &domain.User{ID: "1", Role: domain.RoleGuest},
	}, nil)

	body, _ := json.Marshal(map[string]string{"name": "John", "email": "j@j.com", "password": "secret123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandlerRegister_Conflict(t *testing.T) {
	uc := &mockAuthUsecase{}
	r := setupRouter(uc)
	uc.On("Register", mock.Anything, mock.Anything).Return(nil, domain.ErrUserAlreadyExists)

	body, _ := json.Marshal(map[string]string{"name": "X", "email": "x@x.com", "password": "secret123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestHandlerLogin_OK(t *testing.T) {
	uc := &mockAuthUsecase{}
	r := setupRouter(uc)
	uc.On("Login", mock.Anything, mock.Anything).Return(&domain.TokenResponse{
		AccessToken: "tok",
		User:        &domain.User{ID: "1"},
	}, nil)

	body, _ := json.Marshal(map[string]string{"email": "a@a.com", "password": "pass"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandlerMe_Unauthorized(t *testing.T) {
	uc := &mockAuthUsecase{}
	r := setupRouter(uc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
