package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Gustcat/auth/internal/client/db"
	"github.com/Gustcat/auth/internal/model"
	"github.com/Gustcat/auth/internal/repository"
	"github.com/Gustcat/auth/internal/repository/user/converter"
	modelRepo "github.com/Gustcat/auth/internal/repository/user/model"
	sq "github.com/Masterminds/squirrel"
	"time"
)

const (
	tableName = "auth_user"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	passwordColumn  = "password"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.UserInfo, pwd string) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, roleColumn, passwordColumn).
		Values(info.Name, info.Email, info.Role, pwd).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (int64, *model.UserInfo, time.Time, sql.NullTime, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn, passwordColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, nil, time.Time{}, sql.NullTime{Time: time.Time{}, Valid: false}, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return 0, nil, time.Time{}, sql.NullTime{Time: time.Time{}, Valid: false}, err
	}

	id, userInfo, createdAt, updatedAt := converter.ToUserFromRepo(&user)
	return id, userInfo, createdAt, updatedAt, nil
}

func (r *repo) Update(ctx context.Context, id int64, name string, email string) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{"id": id})

	if name != "" {
		builder = builder.Set(nameColumn, name)
	}

	if email != "" {
		builder = builder.Set(emailColumn, email)

	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	ct, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("запись с id %d не найдена", id)
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	ct, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("запись с id %d не найдена", id)
	}

	return nil
}
