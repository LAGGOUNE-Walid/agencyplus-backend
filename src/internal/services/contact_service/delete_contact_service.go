package contact_service

import (
	"context"
	"logispro/internal/db"
)

type DeleteContactService struct {
	Queries *db.Queries
}

func (s DeleteContactService) Delete(id int64, agencyUsers []int64, ctx context.Context) error {
	return s.Queries.DeleteContact(ctx, db.DeleteContactParams{ID: id, UsersID: agencyUsers})
}
