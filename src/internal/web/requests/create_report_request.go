package requests

import (
	"logispro/internal/web/validations"
	"net/http"
)

type CreateReportRequest struct {
	UserID  int64
	Title   string
	Content string
}

func ParseCreateReportRequest(r *http.Request) (CreateReportRequest, validations.ValidationErrors) {
	var req CreateReportRequest
	validationErrors := validations.ValidateCreateReportRequest(r)
	if len(validationErrors) > 0 {
		return req, validationErrors
	}
	req.Title = r.FormValue("title")
	req.Content = r.FormValue("content")

	return req, validationErrors
}
