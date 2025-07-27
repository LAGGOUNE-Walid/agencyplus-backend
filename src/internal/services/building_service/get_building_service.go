package building_service

import (
	"context"
	"logispro/internal/db"
)

type GetBuildingService struct {
	Queries *db.Queries
}
type FullBuilding struct {
	Building  db.Building           `json:"building"`
	Images    []db.BuildingImage    `json:"images"`
	Documents []db.BuildingDocument `json:"documents"`
}
type PaginatedBuildingsResponse struct {
	Data    []FullBuilding `json:"data"`
	HasMore bool           `json:"has_more"`
}

func (s *GetBuildingService) All(agencyUsers []int64, offset int64, limit int64, ctx context.Context) (*PaginatedBuildingsResponse, error) {
	var params db.ListPaginatedBuildingsParams
	params.UsersID = agencyUsers
	params.Offset = offset
	params.Limit = limit + 1 // fetch one extra to check hasMore

	buildings, err := s.Queries.ListPaginatedBuildings(ctx, params)
	if err != nil {
		return nil, err
	}

	hasMore := len(buildings) > int(limit)
	if hasMore {
		buildings = buildings[:limit] // trim to actual limit
	}

	ids := make([]int64, len(buildings))
	for i, b := range buildings {
		ids[i] = b.ID
	}

	images, err := s.Queries.ListImagesForBuildingIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	docs, err := s.Queries.ListDocumentsForBuildingIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	imageMap := make(map[int64][]db.BuildingImage)
	for _, img := range images {
		imageMap[img.BuildingID] = append(imageMap[img.BuildingID], img)
	}

	docMap := make(map[int64][]db.BuildingDocument)
	for _, doc := range docs {
		docMap[doc.BuildingID] = append(docMap[doc.BuildingID], doc)
	}

	var full []FullBuilding
	for _, b := range buildings {
		full = append(full, FullBuilding{
			Building:  b,
			Images:    imageMap[b.ID],
			Documents: docMap[b.ID],
		})
	}

	return &PaginatedBuildingsResponse{
		Data:    full,
		HasMore: hasMore,
	}, nil
}

func (s *GetBuildingService) Get(agencyUsers []int64, id int64, ctx context.Context) (FullBuilding, error) {
	var full FullBuilding
	var params db.GetBuildingParams
	params.ID = id
	params.UsersID = agencyUsers
	b, err := s.Queries.GetBuilding(ctx, params)
	if err != nil {
		return full, err
	}
	ids := make([]int64, 1)
	ids = append(ids, b.ID)

	images, err := s.Queries.ListImagesForBuildingIDs(ctx, ids)
	if err != nil {
		return full, err
	}
	docs, err := s.Queries.ListDocumentsForBuildingIDs(ctx, ids)
	if err != nil {
		return full, err
	}
	full.Building = b
	full.Documents = docs
	full.Images = images
	return full, nil
}
