package domain

import (
	"errors"
	"time"
)

// BookingStatus represents the lifecycle state of a booking.
type BookingStatus string

const (
	StatusPending     BookingStatus = "pending"
	StatusConfirmed   BookingStatus = "confirmed"
	StatusRejected    BookingStatus = "rejected"
	StatusPendingEdit BookingStatus = "pending_edit"
)

var (
	ErrBookingNotFound  = errors.New("booking not found")
	ErrForbidden        = errors.New("access forbidden")
	ErrConflict         = errors.New("dates conflict with existing booking")
)

// Booking is the core booking entity.
type Booking struct {
	ID        string        `json:"id"`
	UserID    string        `json:"user_id"`
	RoomID    string        `json:"room_id"`
	Room      *Room         `json:"room,omitempty"`
	CheckIn   time.Time     `json:"check_in"`
	CheckOut  time.Time     `json:"check_out"`
	Status    BookingStatus `json:"status"`
	Comment   string        `json:"comment"`
	GuestName string        `json:"guest_name,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
}

// CreateBookingRequest is the DTO for creating a new booking.
type CreateBookingRequest struct {
	RoomID   string    `json:"room_id"   binding:"required"`
	CheckIn  time.Time `json:"check_in"  binding:"required"`
	CheckOut time.Time `json:"check_out" binding:"required"`
	Comment  string    `json:"comment"`
}

// UpdateBookingRequest is the DTO for editing a booking (client side).
type UpdateBookingRequest struct {
	CheckIn  time.Time `json:"check_in"  binding:"required"`
	CheckOut time.Time `json:"check_out" binding:"required"`
	Comment  string    `json:"comment"`
}

// ChangeStatusRequest is the DTO for admin status changes.
type ChangeStatusRequest struct {
	Status BookingStatus `json:"status" binding:"required"`
}
