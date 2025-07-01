package requests

import (
	"logispro/internal/web/validations"
	"net/http"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ParseAuthRequest(r *http.Request) (AuthRequest, validations.ValidationErrors) {
	var req AuthRequest
	validationErrors := validations.ValidateUserAuthRequest(r)

	req.Email = r.FormValue("email")
	req.Password = r.FormValue("password")
	return req, validationErrors
}
