package validations

import (
	"logispro/internal/constants"
	"mime/multipart"
	"net/http"
	"sync"
)

func ValidateCreateBuildingDocumentsRequest(r *http.Request) (ValidationErrors, []*multipart.FileHeader, error) {
	errs := make(ValidationErrors)
	var (
		wg              sync.WaitGroup
		fileErrs        = make(chan fieldError, 20)
		validDocHeaders []*multipart.FileHeader
	)
	// Validate documents[]
	if docs, ok := r.MultipartForm.File["documents"]; ok {
		if len(docs) > constants.MaxBuildingDocuments {
			errs.Add("documents[]", "max_50_files")
		} else {
			for _, fileHeader := range docs {
				if fileHeader.Size > constants.MaxBuildingDocumentSize {
					errs.Add("documents[]", "max_size_5MB")
					continue
				}
				validDocHeaders = append(validDocHeaders, fileHeader)
				wg.Add(1)
				go func(fh *multipart.FileHeader) {
					defer wg.Done()
					mime, err := detectMimeType(fh)
					if err != nil || !isValidPDFMime(mime) {
						fileErrs <- fieldError{Field: "documents[]", Message: "invalid_pdf"}
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
		return nil, validDocHeaders, nil
	}
	return errs, nil, nil
}
