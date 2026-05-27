package domain

import "errors"

var ErrRoomNotFound = errors.New("room not found")

// Room represents a hotel room available for booking.
type Room struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Capacity    int     `json:"capacity"`
	ImageURL    string  `json:"image_url"`
}
