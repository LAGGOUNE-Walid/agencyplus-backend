package contact_service

import (
	"context"
	"logispro/internal/db"
)

type DeleteContactService struct {
	Queries *db.Queries
}

func (s DeleteContactService) Delete(id int64, userId int64, ctx context.Context) error {
	return s.Queries.DeleteContact(ctx, db.DeleteContactParams{ID: id, UserID: userId})
}
