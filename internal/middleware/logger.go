package middleware

import (
	"context"
	"github.com/Gustcat/task-server/internal/logger"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func LoggerMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := uuid.NewString()
		requestLogger := log.With(slog.String("request_id", reqID))

		ctx := context.WithValue(c.Request.Context(), logger.LoggerKey, requestLogger)
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		log.Info("request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", duration),
		)
	}
}
