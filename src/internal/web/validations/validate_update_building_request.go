package validations

import (
	"net/http"
)

func ValidateUpdateBuildingInfoRequest(r *http.Request) ValidationErrors {
	errs := make(ValidationErrors)

	title := r.FormValue("title")
	price := r.FormValue("price")
	status := r.FormValue("status")

	ValidateNonEmpty(title, "title", "requis", errs)
	ValidateNonEmpty(price, "price", "requis", errs)
	ValidateNonEmpty(status, "status", "requis", errs)

	return errs
}
