package task_service

import (
	"context"
	"database/sql"
	"logispro/internal/db"
	"logispro/internal/web/requests"
)

type CreateTaskService struct {
	Queries *db.Queries
}

func (s *CreateTaskService) Create(req requests.CreateTaskRequest, rootId int64, ctx context.Context) (db.Task, error) {
	return s.Queries.CreateTask(ctx, db.CreateTaskParams{
		ToID:    req.To,
		Title:   req.Title,
		Content: req.Content,
		DueDate: sql.NullTime{Time: req.Date, Valid: !req.Date.IsZero()},
		RootID:  rootId,
	})
}
