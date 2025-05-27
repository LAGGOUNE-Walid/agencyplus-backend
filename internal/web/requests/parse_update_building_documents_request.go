package requests

import (
	"logispro/internal/web/validations"
	"mime/multipart"
	"net/http"
)

type UpdateBuildingDocuments struct {
	UserID          int64
	DocumentFiles   []multipart.File
	DocumentHeaders []*multipart.FileHeader
}

func ParseUpdateBuildingDocumentsRequest(r *http.Request, userID int64) (UpdateBuildingDocuments, validations.ValidationErrors, error) {
	var req UpdateBuildingDocuments
	req.UserID = userID
	validationErrors, documentHeaders, err := validations.ValidateCreateBuildingDocumentsRequest(r)
	if err != nil {
		return req, nil, err
	}
	if len(validationErrors) > 0 {
		return req, validationErrors, nil
	}
	req.DocumentHeaders = documentHeaders
	for _, hdr := range documentHeaders {
		file, err := hdr.Open()
		if err != nil {
			return req, nil, err
		}
		req.DocumentFiles = append(req.DocumentFiles, file)
		defer file.Close()
	}
	return req, validationErrors, nil
}
