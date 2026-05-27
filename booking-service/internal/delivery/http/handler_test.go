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
	commonjwt "github.com/hotelapp/common/pkg/jwt"
	"github.com/hotelapp/booking-service/internal/domain"
	delivery "github.com/hotelapp/booking-service/internal/delivery/http"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─── mocks ────────────────────────────────────────────────────────────────────

type mockRoomUC struct{ mock.Mock }

func (m *mockRoomUC) GetAll(ctx context.Context) ([]*domain.Room, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]*domain.Room), args.Error(1)
}
func (m *mockRoomUC) GetByID(ctx context.Context, id string) (*domain.Room, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*domain.Room), args.Error(1)
}

type mockBookingUC struct{ mock.Mock }

func (m *mockBookingUC) Create(ctx context.Context, uid string, req *domain.CreateBookingRequest) (*domain.Booking, error) {
	args := m.Called(ctx, uid, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*domain.Booking), args.Error(1)
}
func (m *mockBookingUC) GetAll(ctx context.Context) ([]*domain.Booking, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]*domain.Booking), args.Error(1)
}
func (m *mockBookingUC) GetByUser(ctx context.Context, uid string) ([]*domain.Booking, error) {
	args := m.Called(ctx, uid)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]*domain.Booking), args.Error(1)
}
func (m *mockBookingUC) Update(ctx context.Context, uid, id string, req *domain.UpdateBookingRequest) (*domain.Booking, error) {
	args := m.Called(ctx, uid, id, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*domain.Booking), args.Error(1)
}
func (m *mockBookingUC) Delete(ctx context.Context, uid, id string, isAdmin bool) error {
	args := m.Called(ctx, uid, id, isAdmin)
	return args.Error(0)
}
func (m *mockBookingUC) ChangeStatus(ctx context.Context, id string, s domain.BookingStatus) (*domain.Booking, error) {
	args := m.Called(ctx, id, s)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*domain.Booking), args.Error(1)
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func setupBookingRouter(rUC *mockRoomUC, bUC *mockBookingUC) (*gin.Engine, *commonjwt.Manager) {
	gin.SetMode(gin.TestMode)
	jm := commonjwt.NewManager("secret", 15*time.Minute)
	h := delivery.NewHandler(rUC, bUC, jm, zerolog.Nop())
	r := gin.New()
	r.GET("/api/v1/rooms", h.GetRooms)
	r.GET("/api/v1/rooms/:id", h.GetRoom)
	r.Use(h.JWTMiddleware())
	r.POST("/api/v1/bookings", h.CreateBooking)
	r.GET("/api/v1/bookings/my", h.GetMyBookings)
	r.PUT("/api/v1/bookings/:id", h.UpdateBooking)
	r.DELETE("/api/v1/bookings/:id", h.DeleteBooking)
	admin := r.Group("", h.AdminOnly())
	admin.GET("/api/v1/bookings", h.GetAllBookings)
	admin.PATCH("/api/v1/bookings/:id/status", h.ChangeBookingStatus)
	return r, jm
}

func bearerToken(jm *commonjwt.Manager, uid, role string) string {
	tok, _ := jm.GenerateAccessToken(uid, uid+"@test.com", role)
	return "Bearer " + tok
}

// ─── tests ────────────────────────────────────────────────────────────────────

func TestGetRooms_OK(t *testing.T) {
	rUC := &mockRoomUC{}
	bUC := &mockBookingUC{}
	r, _ := setupBookingRouter(rUC, bUC)

	rUC.On("GetAll", mock.Anything).Return([]*domain.Room{{ID: "1", Name: "Suite"}}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetRoom_NotFound(t *testing.T) {
	rUC := &mockRoomUC{}
	bUC := &mockBookingUC{}
	r, _ := setupBookingRouter(rUC, bUC)

	rUC.On("GetByID", mock.Anything, "nope").Return(nil, domain.ErrRoomNotFound)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms/nope", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateBooking_OK(t *testing.T) {
	rUC := &mockRoomUC{}
	bUC := &mockBookingUC{}
	r, jm := setupBookingRouter(rUC, bUC)

	now := time.Now()
	bUC.On("Create", mock.Anything, "user-1", mock.Anything).Return(&domain.Booking{
		ID: "b1", UserID: "user-1", Status: domain.StatusPending,
		CheckIn: now, CheckOut: now.Add(24 * time.Hour),
	}, nil)

	body, _ := json.Marshal(map[string]interface{}{
		"room_id":   "room-1",
		"check_in":  now.Format(time.RFC3339),
		"check_out": now.Add(24 * time.Hour).Format(time.RFC3339),
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(jm, "user-1", "guest"))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateBooking_Unauthorized(t *testing.T) {
	r, _ := setupBookingRouter(&mockRoomUC{}, &mockBookingUC{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetAllBookings_AdminOnly(t *testing.T) {
	rUC := &mockRoomUC{}
	bUC := &mockBookingUC{}
	r, jm := setupBookingRouter(rUC, bUC)

	// guest should get 403
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings", nil)
	req.Header.Set("Authorization", bearerToken(jm, "guest-1", "guest"))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestChangeStatus_ConfirmEdit(t *testing.T) {
	rUC := &mockRoomUC{}
	bUC := &mockBookingUC{}
	r, jm := setupBookingRouter(rUC, bUC)

	b := &domain.Booking{ID: "b1", Status: domain.StatusConfirmed, Room: &domain.Room{}}
	bUC.On("ChangeStatus", mock.Anything, "b1", domain.StatusConfirmed).Return(b, nil)

	body, _ := json.Marshal(map[string]string{"status": "confirmed"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/bookings/b1/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(jm, "admin-1", "admin"))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
