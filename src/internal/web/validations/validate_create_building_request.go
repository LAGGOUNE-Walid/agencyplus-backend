package validations

import (
	"context"
	"logispro/internal/constants"
	"logispro/internal/db"
	"mime/multipart"
	"net/http"
	"sync"
)

type fieldError struct {
	Field   string
	Message string
}

func ValidateCreateBuildingRequest(
	r *http.Request,
	q *db.Queries,
	ctx context.Context,
) (ValidationErrors, []*multipart.FileHeader, []*multipart.FileHeader, error) {
	errs := make(ValidationErrors)

	title := r.FormValue("title")
	price := r.FormValue("price")
	status := r.FormValue("status")

	ValidateNonEmpty(title, "title", "requis", errs)
	ValidateNonEmpty(price, "price", "requis", errs)
	ValidateNonEmpty(status, "status", "requis", errs)

	var (
		wg                sync.WaitGroup
		fileErrs          = make(chan fieldError, 20)
		validImageHeaders []*multipart.FileHeader
		validDocHeaders   []*multipart.FileHeader
	)
	if r.MultipartForm == nil {
		errs.Add("images", "requis")
		return errs, nil, nil, nil
	}
	// Validate images[]
	if images, ok := r.MultipartForm.File["images"]; ok {

		if len(images) > constants.MaxBuildingImages {
			errs.Add("images[]", "max 50 fichiers")
		} else {
			for _, fileHeader := range images {
				if fileHeader.Size > constants.MaxBuildingImageSize {
					errs.Add("images[]", "taille max 5MB")
					continue
				}
				validImageHeaders = append(validImageHeaders, fileHeader)
				wg.Add(1)
				go func(fh *multipart.FileHeader) {
					defer wg.Done()
					mime, err := detectMimeType(fh)
					if err != nil || !isValidImageMime(mime) {
						fileErrs <- fieldError{Field: "images[]", Message: "image invalide"}
					}
				}(fileHeader)
			}
		}
	}
	if r.MultipartForm == nil {
		errs.Add("documents", "requis")
		return errs, nil, nil, nil
	}
	// Validate documents[]
	if docs, ok := r.MultipartForm.File["documents"]; ok {
		if len(docs) > constants.MaxBuildingDocuments {
			errs.Add("documents[]", "max 50 fichiers")
		} else {
			for _, fileHeader := range docs {
				if fileHeader.Size > constants.MaxBuildingDocumentSize {
					errs.Add("documents[]", "taille max 5MB")
					continue
				}
				validDocHeaders = append(validDocHeaders, fileHeader)
				wg.Add(1)
				go func(fh *multipart.FileHeader) {
					defer wg.Done()
					mime, err := detectMimeType(fh)
					if err != nil || !isValidPDFMime(mime) {
						fileErrs <- fieldError{Field: "documents[]", Message: "pdf invalide"}
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
		return nil, validImageHeaders, validDocHeaders, nil
	}

	return errs, nil, nil, nil
}

// detectMimeType reads first 512 bytes to detect content type
func detectMimeType(fh *multipart.FileHeader) (string, error) {
	file, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buffer[:n]), nil
}

func isValidImageMime(mime string) bool {
	return mime == "image/jpeg" || mime == "image/png"
}

func isValidPDFMime(mime string) bool {
	return mime == "application/pdf"
}
