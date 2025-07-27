package contact_service

import (
	"context"
	"logispro/internal/db"
)

type GetContactService struct {
	Queries *db.Queries
}

func (s *GetContactService) All(agencyUsers []int64, rootId *int64, ctx context.Context) ([]db.Contact, error) {
	return s.Queries.GetAllContacts(ctx, agencyUsers)
}

func (s *GetContactService) Get(id int64, agencyUsers []int64, ctx context.Context) (db.Contact, error) {
	var params db.GetContactParams
	params.UsersID = agencyUsers
	params.ID = id
	return s.Queries.GetContact(ctx, params)
}

func (s *GetContactService) FindAll(ids []int64, agencyUsers []int64, ctx context.Context) ([]db.Contact, error) {
	var params db.GetContactsListParams
	params.UsersID = agencyUsers
	params.ContactIds = ids
	return s.Queries.GetContactsList(ctx, params)
}

func (s *GetContactService) Count(agencyUsers []int64, ctx context.Context) (int64, error) {
	return s.Queries.CountUserContacts(ctx, agencyUsers)
}
