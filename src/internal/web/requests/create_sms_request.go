package requests

import (
	"encoding/json"
	"logispro/internal/web/validations"
	"net/http"
)

type CreateSmsRequest struct {
	Title    string
	Content  string
	Contacts []int64
}

func ParseCreateSmsRequest(r *http.Request) (CreateSmsRequest, validations.ValidationErrors) {
	var req CreateSmsRequest
	validationErrors := validations.ValidateCreateSmsRequest(r)
	if len(validationErrors) > 0 {
		return req, validationErrors
	}
	req.Title = r.FormValue("title")
	req.Content = r.FormValue("content")
	var ids []int64
	if err := json.Unmarshal([]byte(r.FormValue("contacts")), &ids); err == nil {
		req.Contacts = ids
	}
	return req, validationErrors
}
