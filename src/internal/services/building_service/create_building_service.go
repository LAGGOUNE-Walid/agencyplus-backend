package building_service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"logispro/internal/constants"
	"logispro/internal/db"
	"logispro/internal/utils"
	"logispro/internal/web/requests"
	"mime/multipart"
	"path/filepath"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CreateBuildingService struct {
	Queries      *db.Queries
	DB           *sql.DB
	RabbitMqConn *amqp.Connection
}

type BuildingEmbedding struct {
	ID     int64
	Params db.CreateBuildingParams
}

func (s *CreateBuildingService) Create(ctx context.Context, req requests.CreateBuildingRequest, imageHeaders []*multipart.FileHeader, documentHeaders []*multipart.FileHeader) (int64, error) {

	arg := db.CreateBuildingParams{
		UserID:                     req.UserID,
		Title:                      sql.NullString{String: req.Title, Valid: req.Title != ""},
		Location:                   sql.NullString{String: req.Location, Valid: req.Location != ""},
		Wilaya:                     sql.NullString{String: req.Wilaya, Valid: req.Wilaya != ""},
		Daira:                      sql.NullString{String: req.Daira, Valid: req.Daira != ""},
		BuildingType:               sql.NullString{String: req.BuildingType, Valid: req.BuildingType != ""},
		IsPromotionBuilding:        sql.NullBool{Bool: req.IsPromotionBuilding, Valid: true}, // adjust validation if needed
		IsResidency:                sql.NullBool{Bool: req.IsResidency, Valid: true},
		Status:                     sql.NullString{String: req.Status, Valid: req.Status != ""},
		Price:                      sql.NullInt64{Int64: req.Price, Valid: req.Price != 0},
		SurfaceTotal:               sql.NullFloat64{Float64: req.SurfaceTotal, Valid: req.SurfaceTotal != 0},
		SurfaceBuilt:               sql.NullFloat64{Float64: req.SurfaceBuilt, Valid: req.SurfaceBuilt != 0},
		Rooms:                      sql.NullInt64{Int64: req.Rooms, Valid: req.Rooms != 0},
		Bathrooms:                  sql.NullInt64{Int64: req.Bathrooms, Valid: req.Bathrooms != 0},
		FloorsTotal:                sql.NullInt64{Int64: req.FloorsTotal, Valid: req.FloorsTotal != 0},
		ParkingSpaces:              sql.NullInt64{Int64: req.ParkingSpaces, Valid: req.ParkingSpaces != 0},
		IsByTheSea:                 sql.NullBool{Bool: req.IsByTheSea, Valid: true},
		HasWater:                   sql.NullBool{Bool: req.HasWater, Valid: true},
		HasElectricity:             sql.NullBool{Bool: req.HasElectricity, Valid: true},
		HasGas:                     sql.NullBool{Bool: req.HasGas, Valid: true},
		HasInternet:                sql.NullBool{Bool: req.HasInternet, Valid: true},
		HasGarden:                  sql.NullBool{Bool: req.HasGarden, Valid: true},
		HasPool:                    sql.NullBool{Bool: req.HasPool, Valid: true},
		HasElevator:                sql.NullBool{Bool: req.HasElevator, Valid: true},
		HasCentralHeating:          sql.NullBool{Bool: req.HasCentralHeating, Valid: true},
		HasWaterTank:               sql.NullBool{Bool: req.HasWaterTank, Valid: true},
		HasAirConditioner:          sql.NullBool{Bool: req.HasAirConditioner, Valid: true},
		HasEquippedKitchen:         sql.NullBool{Bool: req.HasEquippedKitchen, Valid: true},
		HasTerrace:                 sql.NullBool{Bool: req.HasTerrace, Valid: true},
		HasNotarialDeed:            sql.NullBool{Bool: req.HasNotarialDeed, Valid: true},
		HasLandBooklet:             sql.NullBool{Bool: req.HasLandBooklet, Valid: true},
		HasActInJointOwnership:     sql.NullBool{Bool: req.HasActInJointOwnership, Valid: true},
		HasCertificateOfConformity: sql.NullBool{Bool: req.HasCertificateOfConformity, Valid: true},
		HasDecision:                sql.NullBool{Bool: req.HasDecision, Valid: true},
		HasConcession:              sql.NullBool{Bool: req.HasConcession, Valid: true},
		HasStampedPaper:            sql.NullBool{Bool: req.HasStampedPaper, Valid: true},
		HasBuildingPermit:          sql.NullBool{Bool: req.HasBuildingPermit, Valid: true},
		HasOffPlanSalesContract:    sql.NullBool{Bool: req.HasOffPlanSalesContract, Valid: true},
		BuildingFinishedType:       sql.NullString{String: req.BuildingFinishedType, Valid: req.BuildingFinishedType != ""},
		AcceptablePaymentType:      sql.NullString{String: req.AcceptablePaymentType, Valid: req.AcceptablePaymentType != ""},
		Furnished:                  sql.NullBool{Bool: req.Furnished, Valid: true},
		YearBuilt:                  sql.NullInt64{Int64: req.YearBuilt, Valid: req.YearBuilt != 0},
		Description:                sql.NullString{String: req.Description, Valid: req.Description != ""},
		ShareableLink:              sql.NullString{String: req.ShareableLink, Valid: req.ShareableLink != ""},
	}
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	qtx := s.Queries.WithTx(tx)

	var filesToDelete []string
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			utils.DeleteFiles(filesToDelete...)
		}
	}()

	res, err := qtx.CreateBuilding(ctx, arg)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	for i, header := range req.ImageHeaders {
		imagePath, err := utils.SaveFile(req.ImageFiles[i], header, "uploads/", constants.MaxBuildingImageSize)
		if err != nil {
			return 0, fmt.Errorf("failed to save image: %w", err)
		}
		filesToDelete = append(filesToDelete, imagePath)

		err = qtx.CreateBuildingImage(ctx, db.CreateBuildingImageParams{
			UserID:     req.UserID,
			BuildingID: lastID,
			Path:       imagePath,
			Mimetype:   sql.NullString{String: header.Header.Get("Content-Type"), Valid: true},
			Size:       sql.NullInt64{Int64: header.Size, Valid: true},
		})
		if err != nil {
			return 0, err
		}
	}

	for i, header := range req.DocumentHeaders {
		docPath, err := utils.SaveFile(req.DocumentFiles[i], header, "uploads/", constants.MaxBuildingDocumentSize)
		if err != nil {
			return 0, fmt.Errorf("failed to save document: %w", err)
		}
		filesToDelete = append(filesToDelete, docPath)

		sourceAbsPath, err := filepath.Abs(fmt.Sprintf("uploads/%s", docPath))
		if err != nil {
			return 0, err
		}

		thumbPath := fmt.Sprintf("%s-thumb", sourceAbsPath)
		err = utils.GeneratePDFThumbnail(sourceAbsPath, thumbPath)
		if err != nil {
			return 0, fmt.Errorf("failed to generate thumbnail: %w", err)
		}
		filesToDelete = append(filesToDelete, thumbPath)

		err = qtx.CreateBuildingDocument(ctx, db.CreateBuildingDocumentParams{
			UserID:     req.UserID,
			BuildingID: sql.NullInt64{Valid: true, Int64: lastID},
			Path:       docPath,
			Mimetype:   sql.NullString{String: header.Header.Get("Content-Type"), Valid: true},
			Size:       sql.NullInt64{Int64: header.Size, Valid: true},
			Thumbnail:  sql.NullString{String: fmt.Sprintf("%s-thumb-1.png", docPath), Valid: true},
		})
		if err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		utils.DeleteFiles(filesToDelete...)
		return 0, err
	}
	err = s.EnqueueBuildingEmbeddingGeneration(arg, lastID)
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (s *CreateBuildingService) EnqueueBuildingEmbeddingGeneration(params db.CreateBuildingParams, id int64) error {
	ch, err := s.RabbitMqConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	var buildingEmbedding BuildingEmbedding
	buildingEmbedding.ID = id
	buildingEmbedding.Params = params
	rmq := &utils.RabbitMQ{Conn: s.RabbitMqConn, Channel: ch}
	data, err := json.Marshal(buildingEmbedding)
	if err != nil {
		return err
	}
	return rmq.Publish("created_buildings", data, amqp.Table{"x-retry": 1})
}
