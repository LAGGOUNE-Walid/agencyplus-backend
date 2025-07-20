package requests

import (
	"context"
	"logispro/internal/db"
	"logispro/internal/utils"
	"logispro/internal/web/validations"
	"net/http"
)

type CreateContactRequest struct {
	UserID                 int64    `json:"user_id"`
	FullName               string   `json:"fullname"`
	Phone                  string   `json:"phone"`
	Email                  string   `json:"email"`
	Wilaya                 string   `json:"wilaya"`
	Daira                  string   `json:"daira"`
	ClientType             string   `json:"client_type"`
	PreferredLocationType  string   `json:"preferred_location_type"`
	HouseFinishing         string   `json:"house_finishing"`
	RentingFloorLookingFor string   `json:"renting_floor_looking_for"`
	IsMarried              bool     `json:"is_married"`
	MinBudget              *int64   `json:"min_budget"`
	MaxBudget              *int64   `json:"max_budget"`
	PreferredBuildingTypes string   `json:"preferred_building_types"`
	PreferredFeatures      string   `json:"preferred_features"`
	MinRooms               *int64   `json:"min_rooms"`
	MaxRooms               *int64   `json:"max_rooms"`
	MinSurface             *float64 `json:"min_surface"`
	MaxSurface             *float64 `json:"max_surface"`
	Furnished              *bool    `json:"furnished"`
	AcceptablePaymentType  string   `json:"acceptable_payment_type"`
	MaxYearBuilt           *int64   `json:"max_year_built"`
	PurchaseUrgency        string   `json:"purchase_urgency"`
	Comments               string   `json:"comments"`
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
	req.PreferredLocationType = r.FormValue("preferred_location_type")
	req.HouseFinishing = r.FormValue("house_finishing")
	req.RentingFloorLookingFor = r.FormValue("renting_floor_looking_for")
	req.IsMarried = r.FormValue("is_married") == "yes" || r.FormValue("is_married") == "true"
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
	if minBudget := r.FormValue("min_budget"); len(minBudget) > 0 {
		minBudgetInt := utils.ParseInt64(minBudget)
		req.MinBudget = &minBudgetInt
	}
	if maxBudget := r.FormValue("max_budget"); len(maxBudget) > 0 {
		maxBudgetInt := utils.ParseInt64(maxBudget)
		req.MaxBudget = &maxBudgetInt
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
		isFurnished := r.FormValue("furnished") == "yes" || r.FormValue("furnished") == "true"
		req.Furnished = &isFurnished
	}
	req.AcceptablePaymentType = r.FormValue("acceptable_payment_type")
	if maxYearBuilt := r.FormValue("max_year_built"); len(maxYearBuilt) > 0 {
		maxYearBuiltInt := utils.ParseInt64(maxYearBuilt)
		req.MaxYearBuilt = &maxYearBuiltInt
	}
	req.PurchaseUrgency = r.FormValue("purchase_urgency")
	req.Comments = r.FormValue("comments")
	return req, nil, nil
}
