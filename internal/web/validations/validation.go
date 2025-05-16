package validations

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
)

type ValidationErrors map[string][]string

func (ve ValidationErrors) Add(field, message string) {
	ve[field] = append(ve[field], message)
}

func (ve ValidationErrors) Error() string {
	var sb strings.Builder
	for field, errs := range ve {
		for _, e := range errs {
			sb.WriteString(fmt.Sprintf("%s: %s; ", field, e))
		}
	}
	return strings.TrimSpace(sb.String())
}

func (ve ValidationErrors) IsEmpty() bool {
	return len(ve) == 0
}

func ValidateNonEmpty(value, field, errorMessage string, errs ValidationErrors) {
	if strings.TrimSpace(value) == "" {
		errs.Add(field, errorMessage)
	}
}

func ValidateMinLength(value, field string, min int, errs ValidationErrors) {
	if len(value) < min {
		errs.Add(field, fmt.Sprintf("%s must be at least %d characters", field, min))
	}
}

func ValidateFileIsImage(file multipart.File, header *multipart.FileHeader, maxSize int64, field string, errs ValidationErrors) {
	if header.Size > maxSize {
		errs.Add(field, fmt.Sprintf("file size exceeds %d bytes", maxSize))
		return
	}

	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		errs.Add(field, "failed to read file content")
		return
	}
	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		errs.Add(field, "uploaded file is not a valid image")
	}
	if seeker, ok := file.(interface {
		Seek(int64, int) (int64, error)
	}); ok {
		seeker.Seek(0, 0)
	}
}
