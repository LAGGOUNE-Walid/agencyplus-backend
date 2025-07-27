package contact_service

import (
	"context"
	"database/sql"
	"logispro/internal/db"
)

type GetContactService struct {
	Queries *db.Queries
}

func (s *GetContactService) All(userId int64, rootId *int64, ctx context.Context) ([]db.Contact, error) {
	var params db.GetAllContactsParams
	params.UserID = userId
	if rootId == nil {
		params.UserID2 = sql.NullInt64{
			Valid: false,
		}
	} else {
		params.UserID2 = sql.NullInt64{
			Valid: true,
			Int64: *rootId,
		}
	}
	return s.Queries.GetAllContacts(ctx, params)
}

func (s *GetContactService) Get(id int64, userId int64, rootId *int64, ctx context.Context) (db.Contact, error) {
	var params db.GetContactParams
	params.UserID = userId
	params.ID = id
	if rootId == nil {
		params.UserID2 = sql.NullInt64{
			Valid: false,
		}
	} else {
		params.UserID2 = sql.NullInt64{
			Valid: true,
			Int64: *rootId,
		}
	}
	return s.Queries.GetContact(ctx, params)
}

func (s *GetContactService) FindAll(ids []int64, userId int64, rootId *int64, ctx context.Context) ([]db.Contact, error) {
	var params db.GetContactsListParams
	params.UserID = userId
	params.ContactIds = ids
	if rootId == nil {
		params.UserID2 = sql.NullInt64{
			Valid: false,
		}
	} else {
		params.UserID2 = sql.NullInt64{
			Valid: true,
			Int64: *rootId,
		}
	}
	return s.Queries.GetContactsList(ctx, params)
}

func (s *GetContactService) Count(userId int64, rootId *int64, ctx context.Context) (int64, error) {
	var params db.CountUserContactsParams
	params.UserID = userId
	if rootId == nil {
		params.UserID2 = sql.NullInt64{
			Valid: false,
		}
	} else {
		params.UserID2 = sql.NullInt64{
			Valid: true,
			Int64: *rootId,
		}
	}
	return s.Queries.CountUserContacts(ctx, params)
}
