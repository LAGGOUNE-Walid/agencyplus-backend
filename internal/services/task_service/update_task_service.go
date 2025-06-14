package task_service

import (
	"context"
	"database/sql"
	"logispro/internal/db"
)

type UpdateTaskService struct {
	Queries *db.Queries
}

func (s *UpdateTaskService) MarkAsDone(id int64, userId int64, ctx context.Context) error {
	return s.Queries.MarkTaskAsDone(ctx, db.MarkTaskAsDoneParams{
		ID:     id,
		ToID:   userId,
		RootID: sql.NullInt64{Int64: userId, Valid: true},
	})
}
