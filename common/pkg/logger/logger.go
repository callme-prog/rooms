// Package logger provides a structured logger based on zerolog.
package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// New returns a zerolog.Logger tagged with the given service name.
func New(serviceName string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger()
}
