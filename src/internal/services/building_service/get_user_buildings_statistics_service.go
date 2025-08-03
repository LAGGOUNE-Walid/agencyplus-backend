package building_service

import (
	"context"
	"logispro/internal/db"
)

type GetBuildingsStatisticsService struct {
	Queries *db.Queries
}

type BuildingsStatisticsResult struct {
	Total   int64
	ForSale int64
}

func (s *GetBuildingsStatisticsService) Get(agencyUsers []int64, ctx context.Context) (BuildingsStatisticsResult, error) {
	var buildingsStatisticsResult BuildingsStatisticsResult
	total, err := s.Queries.CountUserBuildings(ctx, agencyUsers)
	if err != nil {
		return buildingsStatisticsResult, err
	}

	buildingsStatisticsResult.Total = total

	forSale, err := s.Queries.CountToSellUserBuildings(ctx, agencyUsers)
	if err != nil {
		return buildingsStatisticsResult, err
	}
	buildingsStatisticsResult.ForSale = forSale
	return buildingsStatisticsResult, nil
}

func (s *GetBuildingsStatisticsService) GetBuildingsTotalChangeRate(agencyUsers []int64, ctx context.Context) (db.GetBuildingsTotalChangeRateRow, error) {
	return s.Queries.GetBuildingsTotalChangeRate(ctx, agencyUsers)
}

func (s *GetBuildingsStatisticsService) GetBuildingsDairaDistribution(agencyUsers []int64, ctx context.Context) ([]db.GetBuildingsDairasRow, error) {
	return s.Queries.GetBuildingsDairas(ctx, agencyUsers)
}
func (s *GetBuildingsStatisticsService) GetBuildingsLocations(agencyUsers []int64, ctx context.Context) ([]db.GetBuildingsMapRow, error) {
	return s.Queries.GetBuildingsMap(ctx, agencyUsers)
}
