package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/task-server/internal/repository"
	modelrepo "github.com/Gustcat/task-server/internal/repository/model"
	"github.com/Gustcat/task-server/internal/repository/postgres/watcher"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
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
	db          *pgxpool.Pool
	watcherRepo repository.WatcherRepository
}

func NewRepo(db *pgxpool.Pool, watcherRepo repository.WatcherRepository) (*Repo, error) {
	return &Repo{
		db:          db,
		watcherRepo: watcherRepo,
	}, nil
}

func (r *Repo) Create(ctx context.Context, task *modelrepo.TaskCreateDB, watcher *string) (int64, error) {
	const op = "postgres.task.Create"

	//TODO: транзакция
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
		return 0, fmt.Errorf("%s: building SQL failed: %w", op, err)
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, repository.ErrTaskExists
		}
		return 0, fmt.Errorf("%s: executing query failed: %w", op, err)
	}

	if watcher != nil {
		err = r.watcherRepo.Add(ctx, id, *watcher)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (r *Repo) Get(ctx context.Context, id int64) (*modelrepo.TaskDB, error) {
	const op = "postgres.task.Get"

	builder := sq.Select(idColumn, titleColumn, descriptionColumn, statusColumn, authorColumn,
		operatorColumn, dueDateColumn, completedAtColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: building SQL failed: %w", op, err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	var task modelrepo.TaskDB
	err = pgxscan.ScanOne(&task, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrTaskNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: executing query failed: %w", op, err)
	}

	return &task, nil
}

func (r *Repo) GetWithWatchers(ctx context.Context, id int64) (*modelrepo.FullTaskDB, error) {
	const op = "postgres.task.GetWithWatchers"

	builder := sq.Select(
		"t."+idColumn,
		"t."+titleColumn,
		"t."+descriptionColumn,
		"t."+statusColumn,
		"t."+authorColumn,
		"t."+operatorColumn,
		"t."+dueDateColumn,
		"t."+completedAtColumn,
		"t."+createdAtColumn,
		"t."+updatedAtColumn,
		"COALESCE(ARRAY_REMOVE(ARRAY_AGG(w.watcher), NULL), '{}') AS watchers",
	).
		From(tableName+" t").
		LeftJoin(watcher.TableName+" w ON w."+watcher.TaskIDColumn+" = t."+idColumn).
		Where(sq.Eq{"t." + idColumn: id}).
		GroupBy(
			"t."+idColumn,
			"t."+titleColumn,
			"t."+descriptionColumn,
			"t."+statusColumn,
			"t."+authorColumn,
			"t."+operatorColumn,
			"t."+dueDateColumn,
			"t."+completedAtColumn,
			"t."+createdAtColumn,
			"t."+updatedAtColumn,
		).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: building SQL failed: %w", op, err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	var task modelrepo.FullTaskDB
	err = pgxscan.ScanOne(&task, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrTaskNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: executing query failed: %w", op, err)
	}

	return &task, nil
}

func (r *Repo) List(ctx context.Context) ([]*modelrepo.TaskDB, error) {
	const op = "postgres.task.List"

	builder := sq.Select(idColumn, titleColumn, descriptionColumn, statusColumn, authorColumn,
		operatorColumn, dueDateColumn, completedAtColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: building SQL failed: %w", op, err)
	}

	// TODO: пагинация и фильтрация

	tasks := make([]*modelrepo.TaskDB, 0)

	err = pgxscan.Select(ctx, r.db, &tasks, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return tasks, nil
}

func (r *Repo) Update(ctx context.Context, id int64, task *modelrepo.TaskUpdateDB) (*modelrepo.TaskDB, error) {
	const op = "postgres.task.Update"
	alias := "p"

	builder := sq.Update(fmt.Sprintf("%s AS %s", tableName, alias)).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Set(updatedAtColumn, task.UpdatedAt)

	if task.Title.Valid {
		builder = builder.Set(titleColumn, task.Title)
	}

	if task.Description.Valid {
		builder = builder.Set(descriptionColumn, task.Description)
	}

	if task.Status.Valid {
		builder = builder.Set(statusColumn, task.Status).
			Set(completedAtColumn, task.CompletedAt)
	}

	if task.Operator.Valid {
		builder = builder.Set(operatorColumn, task.Operator)
	}

	if task.DueDate.Valid {
		builder = builder.Set(dueDateColumn, task.DueDate)
	}

	fields := []string{idColumn, titleColumn, descriptionColumn, statusColumn, operatorColumn,
		dueDateColumn, completedAtColumn, authorColumn, createdAtColumn, updatedAtColumn}
	row := strings.Join(fields, ", ")
	builder = builder.Suffix(fmt.Sprintf("RETURNING %s", row))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: building SQL failed: %w", op, err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	var updatedTask modelrepo.TaskDB
	err = pgxscan.ScanOne(&updatedTask, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrTaskNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: executing query failed: %w", op, err)
	}

	return &updatedTask, nil

}

func (r *Repo) Delete(ctx context.Context, id int64) error {
	const op = "postgres.task.Delete"

	builder := sq.Delete(tableName).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: building SQL failed: %w", op, err)
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: executing query failed: %w", op, err)
	}

	rows := res.RowsAffected()
	if rows == 0 {
		return repository.ErrTaskNotFound
	}

	return nil
}
