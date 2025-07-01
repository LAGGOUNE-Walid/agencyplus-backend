package requests

import (
	"logispro/internal/web/validations"
	"mime/multipart"
	"net/http"
)

type UpdateBuildingImages struct {
	UserID       int64
	ImageFiles   []multipart.File
	ImageHeaders []*multipart.FileHeader
}

func ParseUpdateBuildingImagesRequest(r *http.Request, userID int64) (UpdateBuildingImages, validations.ValidationErrors, error) {
	var req UpdateBuildingImages
	req.UserID = userID

	validationErrors, imageHeaders, err := validations.ValidateCreateBuildingImagesRequest(r)
	if err != nil {
		return req, nil, err
	}
	if len(validationErrors) > 0 {
		return req, validationErrors, nil
	}
	req.ImageHeaders = imageHeaders
	for _, hdr := range imageHeaders {
		file, err := hdr.Open()
		if err != nil {
			return req, nil, err
		}
		req.ImageFiles = append(req.ImageFiles, file)
		defer file.Close()
	}
	return req, validationErrors, nil
}
