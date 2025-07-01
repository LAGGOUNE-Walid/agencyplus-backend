package requests

import (
	"logispro/internal/web/validations"
	"net/http"
	"time"
)

type CreateCalendarEventRequest struct {
	UserId  int64
	Title   string
	Content string
	ForDate time.Time
}

func ParseCreateCalendarEventRequest(r *http.Request) (CreateCalendarEventRequest, validations.ValidationErrors) {
	var req CreateCalendarEventRequest
	validationErrors := validations.ValidateCreateCalendarEventRequest(r)
	if len(validationErrors) > 0 {
		return req, validationErrors
	}
	req.Title = r.FormValue("title")
	req.Content = r.FormValue("content")
	req.ForDate, _ = time.Parse("2006-01-02 15:04:05", r.FormValue("for_date"))
	return req, validationErrors
}
