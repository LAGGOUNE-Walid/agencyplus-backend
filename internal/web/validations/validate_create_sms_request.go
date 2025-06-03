package validations

import (
	"net/http"
)

func ValidateCreateSmsRequest(r *http.Request) ValidationErrors {
	errs := make(ValidationErrors)
	title := r.FormValue("title")
	content := r.FormValue("content")
	contacts := r.FormValue("contacts")

	ValidateNonEmpty(title, "title", "required", errs)
	ValidateNonEmpty(content, "content", "required", errs)
	ValidateNonEmpty(contacts, "contacts", "required", errs)
	ValidJsonOfIntegers(contacts, "contacts", "json", errs)

	// check if ids are in db
	// check quota

	return errs
}
