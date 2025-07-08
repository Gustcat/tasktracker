package middleware

import (
	"context"
	"errors"
	"github.com/Gustcat/task-server/internal/config"
	"github.com/Gustcat/task-server/internal/lib/ctxutils"
	"github.com/Gustcat/task-server/internal/lib/response"
	"github.com/Gustcat/task-server/internal/logger"
	"github.com/Gustcat/task-server/internal/model"
	"github.com/Gustcat/tasktracker/libs/jwt_auth"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

var ErrInvalidToken = errors.New("user is not authorised: invalid token")

func AuthMiddleware(conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := logger.LogFromContext(ctx)

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			log.Error("Missing token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(ErrInvalidToken.Error()))
			return
		}

		if !strings.HasPrefix(authHeader, conf.TokenConfig.AuthPrefix) {
			log.Error("authorization header is not bearer token", slog.String("authHeader", authHeader))
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(ErrInvalidToken.Error()))
			return
		}
		accessToken := strings.TrimPrefix(authHeader, conf.TokenConfig.AuthPrefix)

		claims, err := jwt_auth.VerifyToken(accessToken, []byte(conf.TokenConfig.AccessTokenSecretKey))
		if err != nil {
			log.Error("Error verifying access token",
				slog.String("accessToken", accessToken),
				slog.String("err", err.Error()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(ErrInvalidToken.Error()))
		}

		user := &model.User{
			ID:   claims.ID,
			Name: claims.Name,
			Role: model.Role(claims.Role),
		}
		log.Debug("put user in model", slog.Int64("id", claims.ID), slog.String("name", claims.Name), slog.Int("role", int(model.Role(claims.Role))))

		ctx = context.WithValue(c.Request.Context(), ctxutils.UserKey, user)
		c.Request = c.Request.WithContext(ctx)
	}
}
