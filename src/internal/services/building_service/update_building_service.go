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
	"path/filepath"

	amqp "github.com/rabbitmq/amqp091-go"
)

type UpdateBuildingService struct {
	Queries      *db.Queries
	DB           *sql.DB
	RabbitMqConn *amqp.Connection
}

func (s *UpdateBuildingService) Delete(ctx context.Context, agencyUsers []int64, buildingId int64) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := s.Queries.WithTx(tx)
	var params1 db.DeleteBuildingParams
	params1.UsersID = agencyUsers
	params1.ID = buildingId
	err = qtx.DeleteBuilding(ctx, params1)
	if err != nil {
		return err
	}
	var params2 db.DeleteBuildingImagesParams
	params2.UsersID = agencyUsers
	params2.BuildingID = buildingId
	err = qtx.DeleteBuildingImages(ctx, params2)
	if err != nil {
		return err
	}
	var params3 db.DeleteBuildingDocumentsParams
	params3.UsersID = agencyUsers
	params3.BuildingID = sql.NullInt64{Valid: true, Int64: buildingId}
	err = qtx.DeleteBuildingDocuments(ctx, params3)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *UpdateBuildingService) UpdateBasicInfo(ctx context.Context, agencyUsers []int64, req requests.UpdateBuildingRequest, buildingId int64) error {
	params := db.UpdateBuildingParams{
		Title:                      sql.NullString{String: req.Title, Valid: req.Title != ""},
		Status:                     sql.NullString{String: req.Status, Valid: req.Status != ""},
		Location:                   sql.NullString{String: req.Location, Valid: req.Location != ""},
		Wilaya:                     sql.NullString{String: req.Wilaya, Valid: req.Wilaya != ""},
		Daira:                      sql.NullString{String: req.Daira, Valid: req.Daira != ""},
		BuildingType:               sql.NullString{String: req.BuildingType, Valid: req.BuildingType != ""},
		IsPromotionBuilding:        sql.NullBool{Bool: req.IsPromotionBuilding, Valid: true}, // adjust validation if needed
		IsResidency:                sql.NullBool{Bool: req.IsResidency, Valid: true},
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
		ID:                         buildingId,
		UsersID:                    agencyUsers,
	}
	fmt.Println(params)
	err := s.Queries.UpdateBuilding(ctx, params)
	if err != nil {
		return err
	}
	s.Queries.DeleteEmbeddings(ctx, buildingId)
	ch, err := s.RabbitMqConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	var buildingEmbedding BuildingEmbedding
	buildingEmbedding.ID = buildingId
	createParams := db.CreateBuildingParams{
		Title:                      sql.NullString{String: req.Title, Valid: req.Title != ""},
		Status:                     sql.NullString{String: req.Status, Valid: req.Status != ""},
		Location:                   sql.NullString{String: req.Location, Valid: req.Location != ""},
		Wilaya:                     sql.NullString{String: req.Wilaya, Valid: req.Wilaya != ""},
		Daira:                      sql.NullString{String: req.Daira, Valid: req.Daira != ""},
		BuildingType:               sql.NullString{String: req.BuildingType, Valid: req.BuildingType != ""},
		IsPromotionBuilding:        sql.NullBool{Bool: req.IsPromotionBuilding, Valid: true}, // adjust validation if needed
		IsResidency:                sql.NullBool{Bool: req.IsResidency, Valid: true},
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
	}
	buildingEmbedding.Params = createParams
	rmq := &utils.RabbitMQ{Conn: s.RabbitMqConn, Channel: ch}
	data, err := json.Marshal(buildingEmbedding)
	if err != nil {
		return err
	}
	return rmq.Publish("created_buildings", data, amqp.Table{"x-retry": 1})
}

func (s *UpdateBuildingService) AddImages(ctx context.Context, req requests.UpdateBuildingImages, buildingId int64) error {
	for i, header := range req.ImageHeaders {
		imagePath, err := utils.SaveFile(req.ImageFiles[i], header, "uploads/", constants.MaxBuildingImageSize)
		if err != nil {
			return fmt.Errorf("failed to save image: %w", err)
		}
		err = s.Queries.CreateBuildingImage(ctx, db.CreateBuildingImageParams{
			UserID:     req.UserID,
			BuildingID: buildingId,
			Path:       imagePath,
			Mimetype:   sql.NullString{String: header.Header.Get("Content-Type"), Valid: true},
			Size:       sql.NullInt64{Int64: header.Size, Valid: true},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UpdateBuildingService) DeleteImage(ctx context.Context, agencyUsers []int64, buildingId int64, imageId int64) error {
	params := db.DeleteBuildingImageParams{
		BuildingID: buildingId,
		UsersID:    agencyUsers,
		ID:         imageId,
	}
	return s.Queries.DeleteBuildingImage(ctx, params)
}

func (s *UpdateBuildingService) AddDocuments(ctx context.Context, req requests.UpdateBuildingDocuments, buildingId int64) error {
	for i, header := range req.DocumentHeaders {
		docPath, err := utils.SaveFile(req.DocumentFiles[i], header, "uploads/", constants.MaxBuildingDocumentSize)
		if err != nil {
			return fmt.Errorf("failed to save document: %w", err)
		}
		sourceAbsPath, err := filepath.Abs(fmt.Sprintf("uploads/%s", docPath))
		if err != nil {
			return err
		}

		thumbPath := fmt.Sprintf("%s-thumb", sourceAbsPath)
		err = utils.GeneratePDFThumbnail(sourceAbsPath, thumbPath)
		if err != nil {
			return fmt.Errorf("failed to generate thumbnail of file %s : %w", sourceAbsPath, err)
		}
		err = s.Queries.CreateBuildingDocument(ctx, db.CreateBuildingDocumentParams{
			UserID:     req.UserID,
			BuildingID: sql.NullInt64{Valid: true, Int64: buildingId},
			Path:       docPath,
			Mimetype:   sql.NullString{String: header.Header.Get("Content-Type"), Valid: true},
			Size:       sql.NullInt64{Int64: header.Size, Valid: true},
			Thumbnail:  sql.NullString{String: fmt.Sprintf("%s-thumb.png", docPath), Valid: true},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UpdateBuildingService) DeleteDocument(ctx context.Context, agencyUsers []int64, buildingId int64, documentId int64) error {
	params := db.DeleteBuildingDocumentParams{
		BuildingID: sql.NullInt64{Valid: true, Int64: buildingId},
		UsersID:    agencyUsers,
		ID:         documentId,
	}
	return s.Queries.DeleteBuildingDocument(ctx, params)
}

func (s *UpdateBuildingService) AddVue(ctx context.Context, req requests.CreateBuildingVueRequest) error {
	return s.Queries.CreateBuildingVue(ctx, db.CreateBuildingVueParams{
		BuildingID: req.BuildingId,
		IpAddress:  req.IpAddress.String(),
		UserAgent:  req.UserAgent,
	})
}
