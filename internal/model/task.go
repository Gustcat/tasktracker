package model

import (
	"time"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusTODO       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusWaiting    Status = "waiting"
	StatusDone       Status = "done"
)

type Task struct {
	ID          int64
	Author      int64
	Title       string
	Description *string
	Status      Status
	Operator    *int64
	DueDate     *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type TaskCreate struct {
	Author      int64
	Title       string
	Description *string
	Status      *Status
	Operator    *int64
	DueDate     *time.Time
	CompletedAt *time.Time
}

type TaskUpdate struct {
	Title       *string
	Description *string
	Status      *Status
	Operator    *int64
	DueDate     *time.Time
	CompletedAt *time.Time
	UpdatedAt   time.Time
}
