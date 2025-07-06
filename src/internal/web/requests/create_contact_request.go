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
	ClientType             string // e.g., "buyer", "tenant"
	SearchingFor           string // e.g., "villa a alger", "appartement a bordj el bahri"
	PreferredLocationType  string // e.g., "urban a alger", "suburban a blida"
	RentingFloorLookingFor string
	IsMarried              bool
	MinBudget              *int64
	MaxBudget              *int64
	PreferredBuildingTypes string // e.g., "villa,apartment"
	PreferredFeatures      string // e.g., JSON string: ["has_pool", "has_garden"]
	MinRooms               *int64
	MaxRooms               *int64
	MinSurface             *float64
	MaxSurface             *float64
	Furnished              *bool // Optional if also in features
	AcceptablePaymentType  string
	HouseFinishing         string // drop one if duplicated
	MaxYearBuilt           *int64
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

	if minBudget := r.FormValue("min_budget"); len(minBudget) > 0 {
		minBudgetInt := utils.ParseInt64(minBudget)
		req.MinBudget = &minBudgetInt
	}
	if maxBudget := r.FormValue("max_budget"); len(maxBudget) > 0 {
		maxBudgetInt := utils.ParseInt64(maxBudget)
		req.MaxBudget = &maxBudgetInt
	}

	req.PreferredBuildingTypes = r.FormValue("preferred_building_types")
	req.PreferredFeatures = r.FormValue("preferred_features")

	if minRooms := r.FormValue("min_rooms"); len(minRooms) > 0 {
		minRoomsInt := utils.ParseInt64(minRooms)
		req.MinRooms = &minRoomsInt
	}

	if maxRooms := r.FormValue("max_rooms"); len(maxRooms) > 0 {
		maxRoomsInt := utils.ParseInt64(maxRooms)
		req.MaxRooms = &maxRoomsInt
	}

	if minSurface := r.FormValue("min_surface"); len(minSurface) > 0 {
		minSurfaceInt := utils.ParseFloat(minSurface)
		req.MinSurface = &minSurfaceInt
	}

	if maxSurface := r.FormValue("max_surface"); len(maxSurface) > 0 {
		maxSurfaceInt := utils.ParseFloat(maxSurface)
		req.MaxSurface = &maxSurfaceInt
	}

	if furnished := r.FormValue("furnished"); len(furnished) > 0 {
		isFurnished := r.FormValue("furnished") == "1" || r.FormValue("furnished") == "true"
		req.Furnished = &isFurnished
	}

	req.AcceptablePaymentType = r.FormValue("acceptable_payment_type")
	req.HouseFinishing = r.FormValue("house_finishing")

	if maxYearBuilt := r.FormValue("max_year_built"); len(maxYearBuilt) > 0 {
		maxYearBuiltInt := utils.ParseInt64(maxYearBuilt)
		req.MaxYearBuilt = &maxYearBuiltInt
	}
	return req, nil, nil
}
