package validations

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
	"time"
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

func ValidateDate(value, field, errorMessage string, errs ValidationErrors) {
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		errs.Add(field, errorMessage)
	}
}

func ValidateDateTime(value, field, errorMessage string, errs ValidationErrors) {
	_, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		errs.Add(field, errorMessage)
	}
}
func ValidateDateTimeInFuture(value, field, errorMessage string, errs ValidationErrors) {
	datetime, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		errs.Add(field, errorMessage)
	}
	now := time.Now()
	if !now.Before(datetime) || now.Equal(datetime) {
		errs.Add(field, errorMessage)
	}

}

func ValidateIp(value, field, errorMessage string, errs ValidationErrors) {
	ip := net.ParseIP(value)
	if ip == nil {
		errs.Add(field, errorMessage)
	}
}

func ValidateMinLength(value, field string, min int, errs ValidationErrors) {
	if len(value) < min {
		errs.Add(field, fmt.Sprintf("%s doit contenir au moins %d caractères", field, min))
	}
}

func ValidateFileIsImage(file multipart.File, header *multipart.FileHeader, maxSize int64, field string, errs ValidationErrors) {
	if header.Size > maxSize {
		errs.Add(field, fmt.Sprintf("la taille du fichier dépasse %d octets", maxSize))
		return
	}

	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		errs.Add(field, "échec de lecture du contenu du fichier")
		return
	}
	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		errs.Add(field, "le fichier téléchargé n'est pas une image valide")
	}
	if seeker, ok := file.(interface {
		Seek(int64, int) (int64, error)
	}); ok {
		seeker.Seek(0, 0)
	}
}

func ValidJsonOfIntegers(value, field, errorMessage string, errs ValidationErrors) {
	if value != "" {
		var ids []int64
		if err := json.Unmarshal([]byte(value), &ids); err != nil {
			errs.Add(field, errorMessage)
		} else if len(ids) == 0 {
			errs.Add(field, errorMessage)
		}
	}
}
