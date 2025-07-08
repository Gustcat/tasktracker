package logger

import (
	"context"
	"github.com/Gustcat/task-server/internal/lib/ctxutils"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
)

func SetupLogger(levelLog slog.Level) *slog.Logger {
	var log *slog.Logger

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    2, // МБ
		MaxBackups: 1, // Кол-во старых файлов
		MaxAge:     2, // Дней хранить
		Compress:   true,
	}

	log = slog.New(
		slog.NewJSONHandler(lumberjackLogger, &slog.HandlerOptions{Level: levelLog}),
	)

	return log
}

func LogFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(ctxutils.LoggerKey).(*slog.Logger)
	if !ok {
		return SetupLogger(slog.LevelInfo)
	}
	return logger
}

func LogFromContextAddOP(ctx context.Context, op string) *slog.Logger {
	return LogFromContext(ctx).With(slog.String("op", op))
}
