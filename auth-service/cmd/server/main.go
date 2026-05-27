// @title           Auth Service API
// @version         1.0
// @description     JWT-based authentication microservice for Hotel Booking App
// @host            localhost:8081
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	commonjwt "github.com/hotelapp/common/pkg/jwt"
	"github.com/hotelapp/common/pkg/logger"
	"github.com/hotelapp/auth-service/config"
	delivery "github.com/hotelapp/auth-service/internal/delivery/http"
	"github.com/hotelapp/auth-service/internal/repository/postgres"
	"github.com/hotelapp/auth-service/internal/usecase"
)

func main() {
	cfg := config.Load()
	log := logger.New("auth-service")

	// Database connection pool
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()
	log.Info().Msg("database connected")

	// Run migrations
	m, err := migrate.New(cfg.MigrationsPath, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init migrations")
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}
	log.Info().Msg("migrations applied")

	// Wiring
	jwtManager := commonjwt.NewManager(cfg.JWTSecret, cfg.JWTAccessTTL)
	userRepo   := postgres.NewUserRepo(pool, log)
	authUC     := usecase.NewAuthUsecase(userRepo, jwtManager, log)
	h          := delivery.NewHandler(authUC, jwtManager, log)
	router     := delivery.NewRouter(h)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.HTTPPort).Msg("auth-service started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("listen error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
	}
	log.Info().Msg("auth-service stopped")
}
