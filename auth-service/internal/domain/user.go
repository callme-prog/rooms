package domain

import (
	"errors"
	"time"
)

// Role constants.
const (
	RoleGuest = "guest"
	RoleAdmin = "admin"
)

// Sentinel errors.
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// User is the core domain entity.
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// RegisterRequest is the DTO for user registration.
type RegisterRequest struct {
	Name     string `json:"name"     binding:"required,min=2"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest is the DTO for user login.
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse is returned after successful auth.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	User        *User  `json:"user"`
}
