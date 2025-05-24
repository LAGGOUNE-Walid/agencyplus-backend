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

func (s *GetBuildingService) All(userId int64, offset int64, limit int64, ctx context.Context) ([]FullBuilding, error) {
	var full []FullBuilding
	buildings, err := s.Queries.ListPaginatedBuildings(ctx, db.ListPaginatedBuildingsParams{UserID: userId, Offset: offset, Limit: limit})
	if err != nil {
		return nil, err
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

	// Assemble final result

	for _, b := range buildings {
		full = append(full, FullBuilding{
			Building:  b,
			Images:    imageMap[b.ID],
			Documents: docMap[b.ID],
		})
	}

	return full, nil
}

func (s *GetBuildingService) Get(userId int64, id int64, ctx context.Context) (FullBuilding, error) {
	var full FullBuilding
	b, err := s.Queries.GetBuilding(ctx, db.GetBuildingParams{ID: id, UserID: userId})
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
