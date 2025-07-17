package building_service

import (
	"context"
	"database/sql"
	"logispro/internal/db"
)

type GetBuildingsStatisticsService struct {
	Queries *db.Queries
}

type BuildingsStatisticsResult struct {
	Total   int64
	ForSale int64
}

func (s *GetBuildingsStatisticsService) Get(userId int64, ctx context.Context) (BuildingsStatisticsResult, error) {
	var buildingsStatisticsResult BuildingsStatisticsResult
	total, err := s.Queries.CountUserBuildings(ctx, userId)
	if err != nil {
		return buildingsStatisticsResult, err
	}
	buildingsStatisticsResult.Total = total
	forSale, err := s.Queries.CountUserBuildingsByStatus(ctx, db.CountUserBuildingsByStatusParams{UserID: userId, Status: sql.NullString{String: "A vendre", Valid: true}})
	if err != nil {
		return buildingsStatisticsResult, err
	}
	buildingsStatisticsResult.ForSale = forSale
	return buildingsStatisticsResult, nil
}

func (s *GetBuildingsStatisticsService) GetBuildingsTotalChangeRate(userId int64, ctx context.Context) (db.GetBuildingsTotalChangeRateRow, error) {
	return s.Queries.GetBuildingsTotalChangeRate(ctx, userId)
}

func (s *GetBuildingsStatisticsService) GetBuildingsDairaDistribution(userId int64, ctx context.Context) ([]db.GetBuildingsDairasRow, error) {
	return s.Queries.GetBuildingsDairas(ctx, userId)
}
func (s *GetBuildingsStatisticsService) GetBuildingsLocations(userId int64, ctx context.Context) ([]db.GetBuildingsMapRow, error) {
	return s.Queries.GetBuildingsMap(ctx, userId)
}
