package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/hotelapp/booking-service/docs"
)

func NewRouter(h *Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	// Public
	rooms := api.Group("/rooms")
	{
		rooms.GET("", h.GetRooms)
		rooms.GET("/:id", h.GetRoom)
	}

	// Authenticated
	auth := api.Group("", h.JWTMiddleware())
	{
		bookings := auth.Group("/bookings")
		{
			bookings.POST("", h.CreateBooking)
			bookings.GET("/my", h.GetMyBookings)
			bookings.PUT("/:id", h.UpdateBooking)
			bookings.DELETE("/:id", h.DeleteBooking)
		}

		// Admin only
		admin := auth.Group("", h.AdminOnly())
		{
			admin.GET("/bookings", h.GetAllBookings)
			admin.PATCH("/bookings/:id/status", h.ChangeBookingStatus)
		}
	}

	return r
}
