package requests

import (
	"context"
	"logispro/internal/db"
	"logispro/internal/web/validations"
	"mime/multipart"
	"net/http"
	"strconv"
)

type CreateUserRequest struct {
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
	RootId        *int64
}

func ParseCreateUserRequest(r *http.Request, q *db.Queries, ctx context.Context) (CreateUserRequest, validations.ValidationErrors, error) {
	var req CreateUserRequest
	validationErrors, agencyLogoHeader, err := validations.ValidateCreateUserRequest(r, q, ctx)
	if err != nil {
		return req, nil, err
	}
	if len(validationErrors) > 0 {
		return req, validationErrors, nil
	}

	req.FullName = r.FormValue("fullname")
	req.Email = r.FormValue("email")
	req.Phone = r.FormValue("phone")
	req.AgencyName = r.FormValue("agency_name")
	req.AgencyAddress = r.FormValue("agency_address")
	req.Wilaya = r.FormValue("wilaya")
	req.Daira = r.FormValue("daira")
	req.Password = r.FormValue("password")
	req.LogoHeader = agencyLogoHeader
	if agencyLogoHeader != nil && agencyLogoHeader.Size > 0 {
		file, err := agencyLogoHeader.Open()
		if err != nil {
			return req, nil, err
		}
		req.AgencyLogo = file
		defer file.Close()
	}

	if rootIdStr := r.FormValue("root_id"); rootIdStr != "" {
		if rootIdInt, err := strconv.ParseInt(rootIdStr, 10, 64); err == nil {
			req.RootId = &rootIdInt
		} else {
			return req, nil, err
		}
	} else {
		req.RootId = nil
	}

	return req, validationErrors, nil
}
