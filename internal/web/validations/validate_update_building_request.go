package validations

import (
	"net/http"
)

func ValidateUpdateBuildingInfoRequest(r *http.Request) ValidationErrors {
	errs := make(ValidationErrors)

	title := r.FormValue("title")
	price := r.FormValue("price")
	status := r.FormValue("status")

	ValidateNonEmpty(title, "title", "required", errs)
	ValidateNonEmpty(price, "price", "required", errs)
	ValidateNonEmpty(status, "status", "required", errs)

	return errs
}
