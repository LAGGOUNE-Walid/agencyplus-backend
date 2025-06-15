package validations

import (
	"net/http"
)

func ValidateCreateReportRequest(r *http.Request) ValidationErrors {
	errs := make(ValidationErrors)
	title := r.FormValue("title")
	content := r.FormValue("content")
	ValidateNonEmpty(title, "title", "required", errs)
	ValidateNonEmpty(content, "content", "required", errs)
	return errs
}
