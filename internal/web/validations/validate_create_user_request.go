package validations

import (
	"context"
	"logispro/internal/constants"
	"logispro/internal/db"
	"mime/multipart"
	"net/http"
	"strings"
)

func ValidateCreateUserRequest(r *http.Request, q *db.Queries, ctx context.Context) (ValidationErrors, *multipart.FileHeader, error) {
	errs := make(ValidationErrors)

	fullname := r.FormValue("fullname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	ValidateNonEmpty(fullname, "fullname", "required", errs)
	if fullname != "" {
		ValidateMinLength(fullname, "fullname", 3, errs)
	}

	ValidateNonEmpty(email, "email", "required", errs)
	if email != "" && !strings.Contains(email, "@") {
		errs.Add("email", "valid")
	}
	sameEmailUsers, err := q.CountUsersByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if sameEmailUsers > 0 {
		errs.Add("email", "unique")
	}

	ValidateNonEmpty(password, "password", "required", errs)
	if password != "" {
		ValidateMinLength(password, "password", 3, errs)
	}

	var logoHeader *multipart.FileHeader
	file, header, err := r.FormFile("agency_logo")
	if err != nil {
		if err != http.ErrMissingFile {
			errs.Add("agency_logo", "failed to read uploaded file")
		}
	} else {
		defer file.Close()
		ValidateFileIsImage(file, header, constants.AgencyLogoMaxSize, "agency_logo", errs)
		logoHeader = header
	}

	if errs.IsEmpty() {
		return nil, logoHeader, nil
	}
	return errs, logoHeader, errs
}
