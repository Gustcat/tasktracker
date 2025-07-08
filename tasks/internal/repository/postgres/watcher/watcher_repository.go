package watcher

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/task-server/internal/logger"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"log/slog"
)

const (
	tableName = "task_watchers"

	idColumn      = "id"
	taskIDColumn  = "task_id"
	watcherColumn = "watcher"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewWatcherRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Add(ctx context.Context, taskID, userID int64) error {
	const op = "watcher.Add"

	builder := sq.Insert(tableName).
		Columns(taskIDColumn, watcherColumn).
		Values(taskID, userID).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: building SQL failed: %w", op, err)
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil
		}
		return fmt.Errorf("%s: executing query failed: %w", op, err)
	}

	log := logger.LogFromContextAddOP(ctx, op)
	log.Info("Create watcher for task", slog.Int64("id", id))

	return nil
}

func (r *Repo) Remove(ctx context.Context, taskID, userID int64) error {
	const op = "watcher.Remove"

	return nil
}
