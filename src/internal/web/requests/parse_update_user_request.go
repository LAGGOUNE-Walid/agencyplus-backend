package requests

import (
	"context"
	"logispro/internal/constants"
	"logispro/internal/db"
	"logispro/internal/web/validations"
	"mime/multipart"
	"net/http"
)

type UpdateUserRequest struct {
	UserID        int64
	FullName      string
	Email         string
	Phone         string
	AgencyName    string
	AgencyAddress string
	Wilaya        string
	Daira         string
	Password      string
	AgencyLogo    multipart.File
	LogoHeader    *multipart.FileHeader
	HasPassword   bool
	HasLogo       bool
}

func ParseUpdateUserRequest(r *http.Request, q *db.Queries, ctx context.Context) (UpdateUserRequest, validations.ValidationErrors, error) {
	var req UpdateUserRequest
	validationErrors, agencyLogoHeader, err := validations.ValidateUpdateUserRequest(r, q, ctx, ctx.Value(constants.UserIDContextKey).(int64))
	if err != nil {
		return req, nil, err
	}
	if len(validationErrors) > 0 {
		return req, validationErrors, nil
	}

	// assuming the user id is available in context (from middleware)
	userID := ctx.Value("user_id").(int64)
	req.UserID = userID
	req.FullName = r.FormValue("fullname")
	req.Email = r.FormValue("email")
	req.Phone = r.FormValue("phone")
	req.AgencyName = r.FormValue("agency_name")
	req.AgencyAddress = r.FormValue("agency_address")
	req.Wilaya = r.FormValue("wilaya")
	req.Daira = r.FormValue("daira")

	if pw := r.FormValue("password"); pw != "" {
		req.Password = pw
		req.HasPassword = true
	}

	if agencyLogoHeader != nil && agencyLogoHeader.Size > 0 {
		file, err := agencyLogoHeader.Open()
		if err != nil {
			return req, nil, err
		}
		req.AgencyLogo = file
		req.LogoHeader = agencyLogoHeader
		req.HasLogo = true
		defer file.Close()
	}

	return req, validationErrors, nil
}
