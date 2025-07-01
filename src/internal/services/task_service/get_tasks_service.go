package task_service

import (
	"context"
	"logispro/internal/constants"
	"logispro/internal/db"
)

type GetTasksService struct {
	Queries *db.Queries
}

func (s *GetTasksService) GetForCurrentUser(userId int64, role int64, ctx context.Context) ([]db.Task, error) {
	if role == constants.ROLE_NORMAL {
		return s.Queries.GetCurrentUserTasks(ctx, userId)
	}
	return s.Queries.GetRootUserCreatedTasks(ctx, db.GetRootUserCreatedTasksParams{
		RootID: userId,
		ToID:   userId,
	})
}
