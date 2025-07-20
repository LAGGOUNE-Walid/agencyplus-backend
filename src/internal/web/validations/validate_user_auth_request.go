package validations

import (
	"net/http"
	"strings"
)

func ValidateUserAuthRequest(r *http.Request) ValidationErrors {
	errs := make(ValidationErrors)
	email := r.FormValue("email")
	password := r.FormValue("password")
	ValidateNonEmpty(email, "email", "requis", errs)
	if email != "" && !strings.Contains(email, "@") {
		errs.Add("email", "valid")
	}
	ValidateNonEmpty(password, "password", "requis", errs)
	return errs
}
