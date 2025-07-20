package validations

import (
	"net/http"
)

func ValidateCreateReportRequest(r *http.Request) ValidationErrors {
	errs := make(ValidationErrors)
	title := r.FormValue("title")
	content := r.FormValue("content")
	ValidateNonEmpty(title, "title", "requis", errs)
	ValidateNonEmpty(content, "content", "requis", errs)
	return errs
}
