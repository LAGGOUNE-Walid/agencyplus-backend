package contact_service

import (
	"context"
	"database/sql"
	"logispro/internal/db"
	"logispro/internal/web/requests"
)

type CreateContactService struct {
	Queries *db.Queries
}

func (s *CreateContactService) Create(ctx context.Context, req requests.CreateContactRequest) (int64, error) {
	arg := db.CreateContactParams{
		UserID:                 req.UserID,
		Fullname:               req.FullName,
		Phone:                  sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Email:                  sql.NullString{String: req.Email, Valid: req.Email != ""},
		Wilaya:                 sql.NullString{String: req.Wilaya, Valid: req.Wilaya != ""},
		Daira:                  sql.NullString{String: req.Daira, Valid: req.Daira != ""},
		ClientType:             sql.NullString{String: req.ClientType, Valid: req.ClientType != ""},
		SearchingFor:           sql.NullString{String: req.SearchingFor, Valid: req.SearchingFor != ""},
		PreferredLocationType:  sql.NullString{String: req.PreferredLocationType, Valid: req.PreferredLocationType != ""},
		HouseFinishing:         sql.NullString{String: req.HouseFinishing, Valid: req.HouseFinishing != ""},
		RentingFloorLookingFor: sql.NullString{String: req.RentingFloorLookingFor, Valid: req.RentingFloorLookingFor != ""},
		IsMarried:              sql.NullBool{Bool: req.IsMarried, Valid: true},
		MinBudget:              sql.NullInt64{Int64: req.MinBudget, Valid: req.MinBudget > 0},
		MaxBudget:              sql.NullInt64{Int64: req.MaxBudget, Valid: req.MaxBudget > 0},
	}

	contactID, err := s.Queries.CreateContact(ctx, arg)
	if err != nil {
		return 0, err
	}

	return contactID.LastInsertId()
}
