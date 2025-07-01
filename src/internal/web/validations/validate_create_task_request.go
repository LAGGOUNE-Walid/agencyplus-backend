package validations

import "net/http"

func ValidateCreateTaskRequest(r *http.Request) ValidationErrors {
	errs := make(ValidationErrors)
	title := r.FormValue("title")
	content := r.FormValue("content")
	to := r.FormValue("user_id")
	ValidateNonEmpty(title, "title", "required", errs)
	ValidateNonEmpty(content, "content", "required", errs)
	ValidateNonEmpty(to, "user_id", "required", errs)
	if len(r.FormValue("date")) > 0 {
		ValidateDate(r.FormValue("date"), "date", "date", errs)
	}
	return errs
}
