package contact_service

import (
	"context"
	"database/sql"
	"encoding/json"
	"logispro/internal/db"
	"logispro/internal/utils"
	"logispro/internal/web/requests"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CreateContactService struct {
	Queries      *db.Queries
	RabbitMqConn *amqp.Connection
}

type ContactEmbedding struct {
	ID     int64
	Params db.CreateContactParams
}

func (s *CreateContactService) Create(ctx context.Context, req requests.CreateContactRequest) (int64, error) {
	arg := db.CreateContactParams{
		UserID:   req.UserID,
		Fullname: req.FullName,
		Phone: sql.NullString{
			String: req.Phone,
			Valid:  req.Phone != "",
		},
		Email: sql.NullString{
			String: req.Email,
			Valid:  req.Email != "",
		},
		Wilaya: sql.NullString{
			String: req.Wilaya,
			Valid:  req.Wilaya != "",
		},
		Daira: sql.NullString{
			String: req.Daira,
			Valid:  req.Daira != "",
		},
		ClientType: sql.NullString{
			String: req.ClientType,
			Valid:  req.ClientType != "",
		},
		SearchingFor: sql.NullString{
			String: req.SearchingFor,
			Valid:  req.SearchingFor != "",
		},
		PreferredLocationType: sql.NullString{
			String: req.PreferredLocationType,
			Valid:  req.PreferredLocationType != "",
		},
		HouseFinishing: sql.NullString{
			String: req.HouseFinishing,
			Valid:  req.HouseFinishing != "",
		},
		RentingFloorLookingFor: sql.NullString{
			String: req.RentingFloorLookingFor,
			Valid:  req.RentingFloorLookingFor != "",
		},
		IsMarried: sql.NullBool{
			Bool:  req.IsMarried,
			Valid: true, // always valid since it's a non-pointer bool
		},
		MinBudget: sql.NullInt64{
			Int64: func() int64 {
				if req.MinBudget != nil {
					return *req.MinBudget
				}
				return 0
			}(),
			Valid: req.MinBudget != nil,
		},
		MaxBudget: sql.NullInt64{
			Int64: func() int64 {
				if req.MaxBudget != nil {
					return *req.MaxBudget
				}
				return 0
			}(),
			Valid: req.MaxBudget != nil,
		},
		PreferredBuildingTypes: sql.NullString{
			String: req.PreferredBuildingTypes,
			Valid:  req.PreferredBuildingTypes != "",
		},
		PreferredFeatures: sql.NullString{
			String: req.PreferredFeatures,
			Valid:  req.PreferredFeatures != "",
		},
		MinRooms: sql.NullInt64{
			Int64: func() int64 {
				if req.MinRooms != nil {
					return int64(*req.MinRooms)
				}
				return 0
			}(),
			Valid: req.MinRooms != nil,
		},
		MaxRooms: sql.NullInt64{
			Int64: func() int64 {
				if req.MaxRooms != nil {
					return int64(*req.MaxRooms)
				}
				return 0
			}(),
			Valid: req.MaxRooms != nil,
		},
		MinSurface: sql.NullFloat64{
			Float64: func() float64 {
				if req.MinSurface != nil {
					return *req.MinSurface
				}
				return 0
			}(),
			Valid: req.MinSurface != nil,
		},
		MaxSurface: sql.NullFloat64{
			Float64: func() float64 {
				if req.MaxSurface != nil {
					return *req.MaxSurface
				}
				return 0
			}(),
			Valid: req.MaxSurface != nil,
		},
		Furnished: sql.NullBool{
			Bool: func() bool {
				if req.Furnished != nil {
					return *req.Furnished
				}
				return false
			}(),
			Valid: req.Furnished != nil,
		},
		AcceptablePaymentType: sql.NullString{
			String: req.AcceptablePaymentType,
			Valid:  req.AcceptablePaymentType != "",
		},
		MaxYearBuilt: sql.NullInt64{
			Int64: func() int64 {
				if req.MaxYearBuilt != nil {
					return int64(*req.MaxYearBuilt)
				}
				return 0
			}(),
			Valid: req.MaxYearBuilt != nil,
		},
	}

	contactID, err := s.Queries.CreateContact(ctx, arg)
	if err != nil {
		return 0, err
	}
	lastId, err := contactID.LastInsertId()
	if err != nil {
		return 0, err
	}
	err = s.EnqueueContactEmbeddingGeneration(arg, lastId)
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

func (s *CreateContactService) EnqueueContactEmbeddingGeneration(params db.CreateContactParams, id int64) error {
	ch, err := s.RabbitMqConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	var contactEmbedding ContactEmbedding
	contactEmbedding.ID = id
	contactEmbedding.Params = params
	rmq := &utils.RabbitMQ{Conn: s.RabbitMqConn, Channel: ch}
	data, err := json.Marshal(contactEmbedding)
	if err != nil {
		return err
	}
	return rmq.Publish("created_contacts", data, amqp.Table{"x-retry": 1})
}
