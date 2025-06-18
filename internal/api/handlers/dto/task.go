package dto

import (
	"github.com/Gustcat/task-server/internal/model"
	"strings"
	"time"
)

type CreateTaskRequest struct {
	Title       string       `json:"title" binding:"required,min=2,max=250"`
	Description *string      `json:"description"`
	Status      model.Status `json:"status" binding:"required,oneof=new in_progress done todo"`
	Operator    *int64       `json:"operator" binding:"gte=0"`
	DueDate     *Date        `json:"due_date" time_format:"2006-01-02 15:04:05"` //TODO: кастомная валидация, не меньше текущей даты
}

type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	if str == "null" {
		return nil
	}

	date, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}

	*d = Date(date)
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	if d == nil {
		return []byte("null"), nil
	}

	date := time.Time(*d)

	return []byte(`"` + date.Format("2006-01-02") + `"`), nil
}

type IdResponse struct {
	ID int64 `json:"id"`
}

type TaskResponse struct {
	ID          int64        `json:"id"`
	Author      int64        `json:"author"`
	Title       string       `json:"title"`
	Description *string      `json:"description"`
	Status      model.Status `json:"status" binding:"required,oneof=new in_progress done todo"`
	Operator    *int64       `json:"operator"`
	DueDate     *Date        `json:"due_date" time_format:"2006-01-02"`
	CompletedAt *time.Time   `json:"completed_at" time_format:"2006-01-02 15:04:05"`
	CreatedAt   time.Time    `json:"created_at" time_format:"2006-01-02 15:04:05"`
	UpdatedAt   *time.Time   `json:"updated_at" time_format:"2006-01-02 15:04:05"`
}

type IdUri struct {
	ID int64 `uri:"id" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       *string       `json:"title" binding:"omitempty,min=2,max=250"`
	Description *string       `json:"description"`
	Status      *model.Status `json:"status" binding:"omitempty,oneof=new in_progress done todo"`
	Operator    *int64        `json:"operator" binding:"omitempty,gte=0"`
	DueDate     *Date         `json:"due_date" time_format:"2006-01-02 15:04:05"`
}
