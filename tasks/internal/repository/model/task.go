package model

import (
	"database/sql"
	"github.com/Gustcat/task-server/internal/model"
	"time"
)

type TaskCreateDB struct {
	Title       string              `db:"title"`
	Description sql.Null[string]    `db:"description"`
	Status      model.Status        `db:"status"`
	Author      int64               `db:"author"`
	Watcher     sql.Null[int64]     `db:"watcher"`
	Operator    sql.Null[int64]     `db:"operator"`
	DueDate     sql.Null[time.Time] `db:"due_date"`
	CompletedAt sql.Null[time.Time] `db:"completed_at"`
}

type TaskDB struct {
	ID        int64               `db:"id"`
	CreatedAt time.Time           `db:"created_at"`
	UpdatedAt sql.Null[time.Time] `db:"updated_at"`
	TaskCreateDB
}

type TaskUpdateDB struct {
	Title       sql.Null[string]       `db:"title"`
	Description sql.Null[string]       `db:"description"`
	Status      sql.Null[model.Status] `db:"status"`
	Operator    sql.Null[int64]        `db:"operator"`
	DueDate     sql.Null[time.Time]    `db:"due_date"`
	CompletedAt sql.Null[time.Time]    `db:"completed_at"`
	UpdatedAt   time.Time              `db:"updated_at"`
}
