package access

import (
	"context"
	"github.com/Gustcat/auth/internal/client/db"
	"github.com/Gustcat/auth/internal/repository"
	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "accesses"

	endpointColumn = "endpoint"
	roleColumn     = "role"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AccessRepository {
	return &repo{db: db}
}

func (r *repo) Check(ctx context.Context, role int32, endpoint string) error {
	builder := sq.Select("1").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.And{
			sq.Eq{endpointColumn: endpoint},
			sq.Eq{roleColumn: role},
		}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "auth_repository.Login",
		QueryRaw: query,
	}

	ct, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return err
	}

	return nil
}
