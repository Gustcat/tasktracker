package main

import (
	"context"
	descUser "github.com/Gustcat/auth/pkg/user_v1"
	taskHandler "github.com/Gustcat/task-server/internal/api/handlers/task"
	"github.com/Gustcat/task-server/internal/client/authgrpc/user"
	"github.com/Gustcat/task-server/internal/client/db/pg"
	"github.com/Gustcat/task-server/internal/client/db/transaction"
	"github.com/Gustcat/task-server/internal/config"
	"github.com/Gustcat/task-server/internal/kafka_consumer"
	"github.com/Gustcat/task-server/internal/logger"
	"github.com/Gustcat/task-server/internal/middleware"
	taskRepository "github.com/Gustcat/task-server/internal/repository/postgres/task"
	"github.com/Gustcat/task-server/internal/repository/postgres/watcher"
	taskService "github.com/Gustcat/task-server/internal/service/task"
	"github.com/Gustcat/task-server/internal/validation"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	pgClient, err := pg.New(ctx, conf.Postgres.DSN)
	if err != nil {
		log.Error("failed to connect to db: %w", err)
	}
	defer pgClient.Close()

	err = pgClient.DB().Ping(ctx)
	if err != nil {
		log.Error("doesn't ping db", slog.String("error", err.Error()))
		os.Exit(1)
	}

	watcherRepo := watcher.NewWatcherRepo(pgClient)
	taskRepo := taskRepository.NewRepo(pgClient)

	txManager := transaction.NewTransactionManager(pgClient.DB())

	conn, err := grpc.NewClient(conf.AuthGRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("can't connect to auth server", slog.String("error", err.Error()))
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Error("failed to close to server", slog.String("error", err.Error()))
		}
	}(conn)

	authClient := user.NewClient(descUser.NewUserV1Client(conn))

	service := taskService.NewService(taskRepo, watcherRepo, authClient, txManager)
	handler := taskHandler.NewHandler(service)

	consumer := kafka_consumer.NewConsumer(conf.ConsumerConfig, service)
	go consumer.Run(ctx) // TODO: закончить работу

	log.Debug("Try to setup router")
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware(log))
	router.Use(middleware.AuthMiddleware(conf))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("not_before_now", validation.NotBeforeNowValidator)
		if err != nil {
			return
		}
	}

	r := router.Group("/api/v1/tasks")
	{
		r.POST("/", handler.Create)
		r.GET("/", handler.List)
		r.GET("/:id", handler.Get)
		r.PATCH("/:id", handler.Update)
		r.DELETE("/:id", handler.Delete)
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
