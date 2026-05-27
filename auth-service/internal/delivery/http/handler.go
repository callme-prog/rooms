package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	commonjwt "github.com/hotelapp/common/pkg/jwt"
	"github.com/hotelapp/auth-service/internal/domain"
	"github.com/hotelapp/auth-service/internal/usecase"
	"github.com/rs/zerolog"
)

// Handler holds all HTTP handlers for auth-service.
type Handler struct {
	authUC     usecase.AuthUsecase
	jwtManager *commonjwt.Manager
	log        zerolog.Logger
}

// NewHandler constructs a Handler.
func NewHandler(authUC usecase.AuthUsecase, jwtManager *commonjwt.Manager, log zerolog.Logger) *Handler {
	return &Handler{authUC: authUC, jwtManager: jwtManager, log: log}
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a guest account and returns an access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.RegisterRequest  true  "Registration payload"
// @Success      201   {object}  domain.TokenResponse
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.authUC.Register(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		default:
			h.log.Error().Err(err).Msg("register error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary      Authenticate user
// @Description  Returns a JWT access token for valid credentials
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.LoginRequest  true  "Login payload"
// @Success      200   {object}  domain.TokenResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.authUC.Login(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Me godoc
// @Summary      Get current user
// @Description  Returns the authenticated user profile
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  domain.User
// @Failure      401  {object}  map[string]string
// @Router       /auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, err := h.authUC.Me(c.Request.Context(), userID.(string))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// JWTMiddleware validates Bearer token and sets user claims into gin context.
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
		claims, err := h.jwtManager.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}
