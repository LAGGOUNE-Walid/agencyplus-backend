package validations

import (
	"logispro/internal/constants"
	"mime/multipart"
	"net/http"
	"sync"
)

func ValidateCreateBuildingImagesRequest(r *http.Request) (ValidationErrors, []*multipart.FileHeader, error) {
	errs := make(ValidationErrors)

	var (
		wg                sync.WaitGroup
		fileErrs          = make(chan fieldError, 20)
		validImageHeaders []*multipart.FileHeader
	)

	// Validate images[]
	if images, ok := r.MultipartForm.File["images"]; ok {

		if len(images) > constants.MaxBuildingImages {
			errs.Add("images[]", "max_50_files")
		} else {
			for _, fileHeader := range images {
				if fileHeader.Size > constants.MaxBuildingImageSize {
					errs.Add("images[]", "max_size_5MB")
					continue
				}
				validImageHeaders = append(validImageHeaders, fileHeader)
				wg.Add(1)
				go func(fh *multipart.FileHeader) {
					defer wg.Done()
					mime, err := detectMimeType(fh)
					if err != nil || !isValidImageMime(mime) {
						fileErrs <- fieldError{Field: "images[]", Message: "invalid_image"}
					}
				}(fileHeader)
			}
		}
	}

	// Wait for all goroutines
	go func() {
		wg.Wait()
		close(fileErrs)
	}()

	for fe := range fileErrs {
		errs.Add(fe.Field, fe.Message)
	}

	if errs.IsEmpty() {
		return nil, validImageHeaders, nil
	}

	return errs, nil, nil
}
