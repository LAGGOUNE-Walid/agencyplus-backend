package requests

import (
	"context"
	"logispro/internal/db"
	"logispro/internal/utils"
	"logispro/internal/web/validations"
	"net/http"
)

type CreateContactRequest struct {
	UserID                 int64
	FullName               string
	Phone                  string
	Email                  string
	Wilaya                 string
	Daira                  string
	ClientType             string
	SearchingFor           string
	PreferredLocationType  string
	HouseFinishing         string
	RentingFloorLookingFor string
	IsMarried              bool
	MinBudget              int64
	MaxBudget              int64
}

func ParseCreateContactRequest(r *http.Request, q *db.Queries, ctx context.Context) (CreateContactRequest, validations.ValidationErrors, error) {
	var req CreateContactRequest

	validationErrors, err := validations.ValidateCreateContactRequest(r, q, ctx)
	if err != nil {
		return req, nil, err
	}
	if len(validationErrors) > 0 {
		return req, validationErrors, nil
	}

	req.FullName = r.FormValue("fullname")
	req.Phone = r.FormValue("phone")
	req.Email = r.FormValue("email")
	req.Wilaya = r.FormValue("wilaya")
	req.Daira = r.FormValue("daira")
	req.ClientType = r.FormValue("client_type")
	req.SearchingFor = r.FormValue("searching_for")
	req.PreferredLocationType = r.FormValue("preferred_location_type")
	req.HouseFinishing = r.FormValue("house_finishing")
	req.RentingFloorLookingFor = r.FormValue("renting_floor_looking_for")
	req.IsMarried = r.FormValue("is_married") == "1" || r.FormValue("is_married") == "true"

	// Handle budgets safely
	if minBudget := r.FormValue("min_budget"); minBudget != "" {
		req.MinBudget = utils.ParseInt64(minBudget)
	}
	if maxBudget := r.FormValue("max_budget"); maxBudget != "" {
		req.MaxBudget = utils.ParseInt64(maxBudget)
	}

	return req, nil, nil
}
