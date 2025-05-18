package validations

import (
	"context"
	"logispro/internal/constants"
	"logispro/internal/db"
	"mime/multipart"
	"net/http"
	"strings"
)

func ValidateUpdateUserRequest(r *http.Request, q *db.Queries, ctx context.Context, currentUserID int64) (ValidationErrors, *multipart.FileHeader, error) {
	errs := make(ValidationErrors)

	fullname := r.FormValue("fullname")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	agencyName := r.FormValue("agency_name")
	agencyAddress := r.FormValue("agency_address")
	wilaya := r.FormValue("wilaya")
	daira := r.FormValue("daira")
	password := r.FormValue("password")

	ValidateNonEmpty(fullname, "fullname", "required", errs)
	if fullname != "" {
		ValidateMinLength(fullname, "fullname", 3, errs)
	}

	ValidateNonEmpty(email, "email", "required", errs)
	if email != "" && !strings.Contains(email, "@") {
		errs.Add("email", "valid")
	}
	if email != "" {
		existingUser, err := q.GetUserByEmail(ctx, email)
		if err == nil && existingUser.ID != currentUserID {
			errs.Add("email", "unique")
		}
	}

	ValidateNonEmpty(phone, "phone", "required", errs)
	if phone != "" {
		ValidateMinLength(phone, "phone", 6, errs)
	}

	ValidateNonEmpty(agencyName, "agency_name", "required", errs)
	ValidateNonEmpty(agencyAddress, "agency_address", "required", errs)
	ValidateNonEmpty(wilaya, "wilaya", "required", errs)
	ValidateNonEmpty(daira, "daira", "required", errs)

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
	return errs, logoHeader, nil
}
