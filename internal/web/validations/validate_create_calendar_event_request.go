package validations

import "net/http"

func ValidateCreateCalendarEventRequest(r *http.Request) ValidationErrors {
	errs := make(ValidationErrors)
	title := r.FormValue("title")
	content := r.FormValue("content")
	for_date := r.FormValue("for_date")
	ValidateNonEmpty(title, "title", "required", errs)
	ValidateNonEmpty(content, "content", "required", errs)
	ValidateNonEmpty(for_date, "for_date", "required", errs)
	ValidateDateTime(for_date, "for_date", "datetime", errs)
	ValidateDateTimeInFuture(for_date, "for_date", "future", errs)
	return errs
}
