package validations

import (
	"context"
	"database/sql"
	"logispro/internal/db"
	"net/http"
	"strings"
)

func ValidateCreateContactRequest(r *http.Request, q *db.Queries, ctx context.Context) (ValidationErrors, error) {
	errs := make(ValidationErrors)

	fullname := r.FormValue("fullname")
	email := r.FormValue("email")
	phone := r.FormValue("phone")

	ValidateNonEmpty(fullname, "fullname", "requis", errs)
	ValidateMinLength(fullname, "fullname", 3, errs)

	ValidateNonEmpty(email, "email", "requis", errs)
	if email != "" && !strings.Contains(email, "@") {
		errs.Add("email", "valid")
	} else if email != "" {
		sameEmailContacts, err := q.CountContactsByEmail(ctx, sql.NullString{String: email, Valid: true})
		if err != nil {
			return nil, err
		}
		if sameEmailContacts > 0 {
			errs.Add("email", "unique")
		}
	}

	ValidateNonEmpty(phone, "phone", "requis", errs)
	if phone != "" {
		samePhoneContacts, err := q.CountContactsByPhone(ctx, sql.NullString{String: phone, Valid: true})
		if err != nil {
			return nil, err
		}
		if samePhoneContacts > 0 {
			errs.Add("phone", "unique")
		}
	}

	if errs.IsEmpty() {
		return nil, nil
	}
	return errs, nil
}
