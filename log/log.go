package log

import (
	"context"
	"github.com/tesla59/blaze/config"
	"log/slog"
	"os"
	"sync"
)

var (
	once   sync.Once
	Logger *slog.Logger
)

func Init() {
	once.Do(func() {
		var handler slog.Handler
		if config.GetConfig().Environment == "production" {
			handler = slog.NewJSONHandler(os.Stdout, nil)
		} else {
			handler = slog.NewTextHandler(os.Stdout, nil)
		}
		Logger = slog.New(handler)
	})
}

// WithContext adds fields (e.g., request ID) to the logger.
func WithContext(ctx context.Context) *slog.Logger {
	if l := ctx.Value(loggerKey); l != nil {
		if logger, ok := l.(*slog.Logger); ok {
			return logger
		}
	}
	return Logger
}

type ctxKey string

const loggerKey ctxKey = "logger"

// Inject adds a logger into the context (for per-request logging).
func Inject(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
