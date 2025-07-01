package requests

import (
	"logispro/internal/web/validations"
	"net/http"
)

type UpdateReportRequest struct {
	ID      int64
	UserID  int64
	Title   string
	Content string
}

func ParseUpdateReportRequest(r *http.Request) (UpdateReportRequest, validations.ValidationErrors) {
	var req UpdateReportRequest
	validationErrors := validations.ValidateUpdateReportRequest(r)
	if len(validationErrors) > 0 {
		return req, validationErrors
	}
	req.Title = r.FormValue("title")
	req.Content = r.FormValue("content")

	return req, validationErrors
}
