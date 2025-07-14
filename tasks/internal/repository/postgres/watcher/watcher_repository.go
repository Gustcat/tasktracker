package watcher

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/task-server/internal/client/db"
	"github.com/Gustcat/task-server/internal/logger"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"log/slog"
)

const (
	TableName = "task_watchers"

	TaskIDColumn  = "task_id"
	watcherColumn = "watcher"
)

type Repo struct {
	db db.Client
}

func NewWatcherRepo(db db.Client) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Add(ctx context.Context, taskID int64, username string) error {
	const op = "watcher.Add"

	builder := sq.Insert(TableName).
		Columns(TaskIDColumn, watcherColumn).
		Values(taskID, username).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: building SQL failed: %w", op, err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil
		}
		return fmt.Errorf("%s: executing query failed: %w", op, err)
	}

	log := logger.LogFromContextAddOP(ctx, op)
	log.Info("Create watcher for task",
		slog.Int64("task_id", taskID),
		slog.String("username", username))

	return nil
}

func (r *Repo) Remove(ctx context.Context, taskID int64, username string) error {
	const op = "watcher.Remove"

	builder := sq.Delete(TableName).
		Where(sq.Eq{watcherColumn: username, TaskIDColumn: taskID}).
		PlaceholderFormat(sq.Dollar).Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: building SQL failed: %w", op, err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil
		}
		return fmt.Errorf("%s: executing query failed: %w", op, err)
	}

	log := logger.LogFromContextAddOP(ctx, op)
	log.Info("Create watcher for task",
		slog.Int64("task_id", taskID),
		slog.String("username", username))

	return nil
}
