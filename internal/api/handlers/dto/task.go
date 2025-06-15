package dto

import (
	"github.com/Gustcat/task-server/internal/model"
	"time"
)

type CreateTaskRequest struct {
	Title       string       `json:"title" binding:"required, min=2, max=250"`
	Description *string      `json:"description"`
	Status      model.Status `json:"status" binding:"required,oneof=new in_progress done todo"`
	Operator    *int64       `json:"operator" binding:"gte=0"`
	DueDate     *time.Time   `json:"due_date" time_format:"2006-01-02"` //TODO: кастомная валидация, не меньше текущей даты
}

type CreateResponse struct {
	ID int64 `json:"id"`
}
