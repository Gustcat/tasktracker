package main

import (
	"context"
	"database/sql"
	"github.com/brianvoe/gofakeit"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	dbDSN = "host=localhost port=54321 dbname=user-local user=user-user-local password=user-password-local sslmode=disable"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	builderInsert := sq.Insert("auth_user").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password").
		Values(gofakeit.Name(),
			gofakeit.Email(),
			gofakeit.Password(true, true, true, false, false, 10)).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int
	err = pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with id: %d", userID)

	builderSelect := sq.Select("id", "name", "email", "role", "password", "created_at", "updated_at").
		From("auth_user").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(10)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to select users: %v", err)
	}

	var id int64
	var name, email, password string
	var role int32
	var createdAt time.Time
	var updatedAt sql.NullTime

	for rows.Next() {
		err = rows.Scan(&id, &name, &email, &role, &password, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan user: %v", err)
		}

		log.Printf("id: %d, name: %s, email: %s, role: %d, password: %s, created_at: %v, updated_at: %v\n",
			id, name, email, role, password, createdAt, updatedAt)
	}

	// Делаем запрос на обновление записи в таблице
	builderUpdate := sq.Update("auth_user").
		PlaceholderFormat(sq.Dollar).
		Set("name", gofakeit.Name()).
		Set("email", gofakeit.Email()).
		Set("password", gofakeit.Password(true, true, true, true, false, 10)).
		Set("role", gofakeit.Number(0, 1)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": userID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	builderSelectOne := sq.Select("id", "name", "email", "role", "password", "created_at", "updated_at").
		From("auth_user").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": userID}).
		Limit(1)

	query, args, err = builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &password, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to select users: %v", err)
	}

	log.Printf("id: %d, name: %s, email: %s, role: %d, password: %s, created_at: %v, updated_at: %v\n",
		id, name, email, role, password, createdAt, updatedAt)
}
