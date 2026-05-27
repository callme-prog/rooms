package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	commonjwt "github.com/hotelapp/common/pkg/jwt"
	"github.com/hotelapp/booking-service/internal/domain"
	"github.com/hotelapp/booking-service/internal/usecase"
	"github.com/rs/zerolog"
)

// Handler holds HTTP handlers for booking-service.
type Handler struct {
	roomUC    usecase.RoomUsecase
	bookingUC usecase.BookingUsecase
	jwt       *commonjwt.Manager
	log       zerolog.Logger
}

func NewHandler(roomUC usecase.RoomUsecase, bookingUC usecase.BookingUsecase,
	jwtMgr *commonjwt.Manager, log zerolog.Logger) *Handler {
	return &Handler{roomUC: roomUC, bookingUC: bookingUC, jwt: jwtMgr, log: log}
}

// ─── JWT Middleware ────────────────────────────────────────────────────────────

func (h *Handler) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}
		claims, err := h.jwt.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func (h *Handler) AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, _ := c.Get("role"); role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}
		c.Next()
	}
}

// ─── Room Handlers ─────────────────────────────────────────────────────────────

// GetRooms godoc
// @Summary      List all rooms
// @Tags         rooms
// @Produce      json
// @Success      200  {array}   domain.Room
// @Router       /rooms [get]
func (h *Handler) GetRooms(c *gin.Context) {
	rooms, err := h.roomUC.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

// GetRoom godoc
// @Summary      Get room by ID
// @Tags         rooms
// @Param        id  path  string  true  "Room ID"
// @Produce      json
// @Success      200  {object}  domain.Room
// @Failure      404  {object}  map[string]string
// @Router       /rooms/{id} [get]
func (h *Handler) GetRoom(c *gin.Context) {
	room, err := h.roomUC.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, domain.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, room)
}

// ─── Booking Handlers ──────────────────────────────────────────────────────────

// CreateBooking godoc
// @Summary      Create a booking
// @Tags         bookings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.CreateBookingRequest  true  "Booking payload"
// @Success      201   {object}  domain.Booking
// @Failure      400   {object}  map[string]string
// @Router       /bookings [post]
func (h *Handler) CreateBooking(c *gin.Context) {
	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("userID")
	b, err := h.bookingUC.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, b)
}

// GetMyBookings godoc
// @Summary      Get current user bookings
// @Tags         bookings
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}  domain.Booking
// @Router       /bookings/my [get]
func (h *Handler) GetMyBookings(c *gin.Context) {
	bookings, err := h.bookingUC.GetByUser(c.Request.Context(), c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// GetAllBookings godoc
// @Summary      Get all bookings (admin)
// @Tags         bookings
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}  domain.Booking
// @Router       /bookings [get]
func (h *Handler) GetAllBookings(c *gin.Context) {
	bookings, err := h.bookingUC.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// UpdateBooking godoc
// @Summary      Edit a booking (client: sets pending_edit)
// @Tags         bookings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string                      true  "Booking ID"
// @Param        body  body  domain.UpdateBookingRequest true  "Update payload"
// @Success      200   {object}  domain.Booking
// @Failure      403   {object}  map[string]string
// @Router       /bookings/{id} [put]
func (h *Handler) UpdateBooking(c *gin.Context) {
	var req domain.UpdateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b, err := h.bookingUC.Update(c.Request.Context(), c.GetString("userID"), c.Param("id"), &req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBookingNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		case errors.Is(err, domain.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	c.JSON(http.StatusOK, b)
}

// ChangeBookingStatus godoc
// @Summary      Change booking status (admin)
// @Tags         bookings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string                      true  "Booking ID"
// @Param        body  body  domain.ChangeStatusRequest  true  "Status payload"
// @Success      200   {object}  domain.Booking
// @Router       /bookings/{id}/status [patch]
func (h *Handler) ChangeBookingStatus(c *gin.Context) {
	var req domain.ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b, err := h.bookingUC.ChangeStatus(c.Request.Context(), c.Param("id"), req.Status)
	if err != nil {
		if errors.Is(err, domain.ErrBookingNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, b)
}

// DeleteBooking godoc
// @Summary      Delete a booking
// @Tags         bookings
// @Security     BearerAuth
// @Param        id  path  string  true  "Booking ID"
// @Success      204
// @Failure      403  {object}  map[string]string
// @Router       /bookings/{id} [delete]
func (h *Handler) DeleteBooking(c *gin.Context) {
	isAdmin := c.GetString("role") == "admin"
	if err := h.bookingUC.Delete(c.Request.Context(), c.GetString("userID"), c.Param("id"), isAdmin); err != nil {
		switch {
		case errors.Is(err, domain.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		case errors.Is(err, domain.ErrBookingNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
