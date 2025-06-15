package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/task-server/internal/repository"
	modelrepo "github.com/Gustcat/task-server/internal/repository/model"
	"github.com/jackc/pgconn"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	tableName = "task"

	idColumn          = "id"
	titleColumn       = "title"
	descriptionColumn = "description"
	statusColumn      = "status"
	authorColumn      = "author"
	operatorColumn    = "operator"
	dueDateColumn     = "due_date"
	completedAtColumn = "completed_at"
	createdAtColumn   = "created_at"
	updatedAtColumn   = "updated_at"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(ctx context.Context, DSN string) (*Repo, error) {
	const op = "repository.postgres.NewRepo"

	db, err := pgxpool.Connect(ctx, DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db, %s: %w", op, err)
	}
	return &Repo{db: db}, nil
}

func (p *Repo) Ping(ctx context.Context) error {
	return p.db.Ping(ctx)
}

func (r *Repo) Close() {
	r.db.Close()
}

func (r *Repo) Create(ctx context.Context, task *modelrepo.TaskCreateDB) (int64, error) {
	const op = "postgres.Create"

	builder := sq.Insert(tableName).
		Columns(titleColumn,
			descriptionColumn,
			statusColumn,
			authorColumn,
			operatorColumn,
			dueDateColumn,
			completedAtColumn).
		Values(task.Title,
			task.Description,
			task.Status,
			task.Author,
			task.Operator,
			task.DueDate,
			task.CompletedAt).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, repository.ErrTaskExists
		}
		return 0, err
	}

	return id, nil
}
