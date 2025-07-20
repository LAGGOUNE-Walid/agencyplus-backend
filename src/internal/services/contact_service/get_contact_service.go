package contact_service

import (
	"context"
	"logispro/internal/db"
)

type GetContactService struct {
	Queries *db.Queries
}

func (s *GetContactService) All(userId int64, ctx context.Context) ([]db.Contact, error) {
	return s.Queries.GetAllContacts(ctx, userId)
}

func (s *GetContactService) Get(id int64, userId int64, ctx context.Context) (db.Contact, error) {
	return s.Queries.GetContact(ctx, db.GetContactParams{UserID: userId, ID: id})
}

func (s *GetContactService) FinAll(ids []int64, userId int64, ctx context.Context) ([]db.Contact, error) {
	return s.Queries.GetContactsList(ctx, db.GetContactsListParams{ContactIds: ids, UserID: userId})
}

func (s *GetContactService) Count(userId int64, ctx context.Context) (int64, error) {
	return s.Queries.CountUserContacts(ctx, userId)
}
