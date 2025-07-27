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

func (s *GetBuildingsStatisticsService) Get(userId int64, rootId *int64, ctx context.Context) (BuildingsStatisticsResult, error) {
	var buildingsStatisticsResult BuildingsStatisticsResult
	var paramsAll db.CountUserBuildingsParams
	paramsAll.UserID = userId
	if rootId != nil {
		paramsAll.UserID2 = sql.NullInt64{Valid: true, Int64: *rootId}
	}
	total, err := s.Queries.CountUserBuildings(ctx, paramsAll)
	if err != nil {
		return buildingsStatisticsResult, err
	}

	buildingsStatisticsResult.Total = total

	var paramsSale db.CountUserBuildingsByStatusParams
	paramsSale.UserID = userId
	paramsSale.Status = sql.NullString{String: "A vendre", Valid: true}
	if rootId != nil {
		paramsSale.UserID2 = sql.NullInt64{Valid: true, Int64: *rootId}
	}
	forSale, err := s.Queries.CountUserBuildingsByStatus(ctx, paramsSale)
	if err != nil {
		return buildingsStatisticsResult, err
	}

	buildingsStatisticsResult.ForSale = forSale
	return buildingsStatisticsResult, nil
}

func (s *GetBuildingsStatisticsService) GetBuildingsTotalChangeRate(userId int64, rootId *int64, ctx context.Context) (db.GetBuildingsTotalChangeRateRow, error) {
	var params db.GetBuildingsTotalChangeRateParams
	params.UserID = userId
	if rootId != nil {
		params.UserID2 = sql.NullInt64{Valid: true, Int64: *rootId}
	}
	return s.Queries.GetBuildingsTotalChangeRate(ctx, params)
}

func (s *GetBuildingsStatisticsService) GetBuildingsDairaDistribution(userId int64, rootId *int64, ctx context.Context) ([]db.GetBuildingsDairasRow, error) {
	var params db.GetBuildingsDairasParams
	params.UserID = userId
	if rootId != nil {
		params.UserID2 = sql.NullInt64{Valid: true, Int64: *rootId}
	}
	return s.Queries.GetBuildingsDairas(ctx, params)
}
func (s *GetBuildingsStatisticsService) GetBuildingsLocations(userId int64, rootId *int64, ctx context.Context) ([]db.GetBuildingsMapRow, error) {
	var params db.GetBuildingsMapParams
	params.UserID = userId
	if rootId != nil {
		params.UserID2 = sql.NullInt64{Valid: true, Int64: *rootId}
	}
	return s.Queries.GetBuildingsMap(ctx, params)
}
