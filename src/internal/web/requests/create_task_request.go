package requests

import (
	"logispro/internal/web/validations"
	"net/http"
	"strconv"
	"time"
)

type CreateTaskRequest struct {
	Title   string
	Content string
	To      int64
	Date    time.Time
}

func ParseCreateTaskRequest(r *http.Request) (CreateTaskRequest, validations.ValidationErrors) {
	var req CreateTaskRequest
	validationErrors := validations.ValidateCreateTaskRequest(r)
	if len(validationErrors) > 0 {
		return req, validationErrors
	}
	to, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		return req, validationErrors
	}
	req.To = int64(to)
	req.Title = r.FormValue("title")
	req.Content = r.FormValue("content")
	if len(r.FormValue("date")) > 0 {
		date, _ := time.Parse("2006-01-02", r.FormValue("date"))
		req.Date = date
	}
	return req, validationErrors
}
