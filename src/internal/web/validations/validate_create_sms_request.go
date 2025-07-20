package validations

import (
	"net/http"
)

func ValidateCreateSmsRequest(r *http.Request) ValidationErrors {
	errs := make(ValidationErrors)
	title := r.FormValue("title")
	content := r.FormValue("content")
	contacts := r.FormValue("contacts")

	ValidateNonEmpty(title, "title", "requis", errs)
	ValidateNonEmpty(content, "content", "requis", errs)
	ValidateNonEmpty(contacts, "contacts", "requis", errs)
	ValidJsonOfIntegers(contacts, "contacts", "json", errs)

	// check if ids are in db
	// check quota

	return errs
}
