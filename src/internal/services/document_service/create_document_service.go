package document_service

import (
	"context"
	"database/sql"
	"fmt"
	"logispro/internal/constants"
	"logispro/internal/db"
	"logispro/internal/utils"
	"logispro/internal/web/requests"
	"path/filepath"
)

type CreateDocumentService struct {
	Queries *db.Queries
}

func (s *CreateDocumentService) Create(ctx context.Context, req requests.UpdateBuildingDocuments) error {
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
			BuildingID: sql.NullInt64{Valid: false},
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
