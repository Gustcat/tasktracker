package main

import (
	"context"
	taskHandler "github.com/Gustcat/task-server/internal/api/handlers/task"
	"github.com/Gustcat/task-server/internal/config"
	"github.com/Gustcat/task-server/internal/logger"
	"github.com/Gustcat/task-server/internal/repository/postgres"
	taskService "github.com/Gustcat/task-server/internal/service/task"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
)

func main() {
	log := logger.SetupLogger(slog.LevelInfo)

	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Warn("doesn't load env file: %s", slog.String("error", err.Error()))
	}

	conf, err := config.New()
	if err != nil {
		log.Error("doesn't set config: %s", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if conf.Env == envLocal {
		log = logger.SetupLogger(slog.LevelDebug)
	}

	log.Debug("Try to connect to db", slog.String("DSN", conf.Postgres.DSN))

	repo, err := postgres.NewRepo(ctx, conf.Postgres.DSN)
	if err != nil {
		log.Error("doesn't create repo", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer repo.Close()

	err = repo.Ping(ctx)
	if err != nil {
		log.Error("doesn't ping db", slog.String("error", err.Error()))
		os.Exit(1)
	}

	service := taskService.NewService(repo)
	handler := taskHandler.NewHandler(service)

	log.Debug("Try to setup router")
	router := gin.Default()

	r := router.Group("/api/v1/tasks")
	{
		r.POST("/", handler.Create(ctx, log))
		//r.GET("/", task.List(ctx, log))
		r.GET("/:id", handler.Get(ctx, log))
		//r.PATCH("/:id", task.Update(ctx, log))
		//r.DELETE("/:id", task.Delete(ctx, log))
	}

	srv := &http.Server{
		Addr:         conf.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  conf.HTTPServer.Timeout,
		WriteTimeout: conf.HTTPServer.Timeout,
		IdleTimeout:  conf.HTTPServer.IdleTimeout,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start http server", slog.String("error", err.Error()))
		}
	}()

	log.Info("Server started", slog.String("address", conf.HTTPServer.Address))

	<-quit

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown http server", slog.String("error", err.Error()))
		return
	}

	log.Info("Server stopped", slog.String("address", conf.HTTPServer.Address))
}
