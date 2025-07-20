package validations

import "net/http"

func ValidateCreateTaskRequest(r *http.Request) ValidationErrors {
	errs := make(ValidationErrors)
	title := r.FormValue("title")
	content := r.FormValue("content")
	to := r.FormValue("user_id")
	ValidateNonEmpty(title, "title", "requis", errs)
	ValidateNonEmpty(content, "content", "requis", errs)
	ValidateNonEmpty(to, "user_id", "requis", errs)
	if len(r.FormValue("date")) > 0 {
		ValidateDate(r.FormValue("date"), "date", "doit être une date valide", errs)
	}
	return errs
}
