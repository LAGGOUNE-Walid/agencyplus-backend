package requests

import (
	"logispro/internal/utils"
	"logispro/internal/web/validations"
	"net/http"
)

type UpdateBuildingRequest struct {
	UserID                     int64
	Location                   string
	Title                      string
	Wilaya                     string
	Daira                      string
	BuildingType               string
	IsPromotionBuilding        bool
	IsResidency                bool
	Status                     string
	Price                      int64
	SurfaceTotal               float64
	SurfaceBuilt               float64
	Rooms                      int64
	Bathrooms                  int64
	FloorsTotal                int64
	ParkingSpaces              int64
	IsByTheSea                 bool
	HasWater                   bool
	HasElectricity             bool
	HasGas                     bool
	HasInternet                bool
	HasGarden                  bool
	HasPool                    bool
	HasElevator                bool
	HasCentralHeating          bool
	HasWaterTank               bool
	HasAirConditioner          bool
	HasEquippedKitchen         bool
	HasTerrace                 bool
	HasNotarialDeed            bool
	HasLandBooklet             bool
	HasActInJointOwnership     bool
	HasCertificateOfConformity bool
	HasDecision                bool
	HasConcession              bool
	HasStampedPaper            bool
	HasBuildingPermit          bool
	HasOffPlanSalesContract    bool
	BuildingFinishedType       string
	AcceptablePaymentType      string
	Furnished                  bool
	YearBuilt                  int64
	Description                string
}

func ParseUpdateBuildingInfoRequest(r *http.Request, userID int64) (UpdateBuildingRequest, validations.ValidationErrors, error) {
	var req UpdateBuildingRequest
	req.UserID = userID

	validationErrors := validations.ValidateUpdateBuildingInfoRequest(r)
	if len(validationErrors) > 0 {
		return req, validationErrors, nil
	}
	// Required fields
	req.Title = r.FormValue("title")
	req.Status = r.FormValue("status")
	req.Price = utils.ParseInt64(r.FormValue("price"))

	// Optional string fields
	req.Location = r.FormValue("location")
	req.Wilaya = r.FormValue("wilaya")
	req.Daira = r.FormValue("daira")
	req.BuildingType = r.FormValue("building_type")
	req.BuildingFinishedType = r.FormValue("building_finished_type")
	req.AcceptablePaymentType = r.FormValue("acceptable_payment_type")
	req.Description = r.FormValue("description")

	// Optional number fields
	req.SurfaceTotal = parseFloat(r.FormValue("surface_total"))
	req.SurfaceBuilt = parseFloat(r.FormValue("surface_built"))
	req.Rooms = utils.ParseInt64(r.FormValue("rooms"))
	req.Bathrooms = utils.ParseInt64(r.FormValue("bathrooms"))
	req.FloorsTotal = utils.ParseInt64(r.FormValue("floors_total"))
	req.ParkingSpaces = utils.ParseInt64(r.FormValue("parking_spaces"))
	req.YearBuilt = utils.ParseInt64(r.FormValue("year_built"))

	// Optional boolean fields
	req.IsPromotionBuilding = parseBool(r.FormValue("is_promotion_building"))
	req.IsResidency = parseBool(r.FormValue("is_residency"))
	req.IsByTheSea = parseBool(r.FormValue("is_by_the_sea"))
	req.HasWater = parseBool(r.FormValue("has_water"))
	req.HasElectricity = parseBool(r.FormValue("has_electricity"))
	req.HasGas = parseBool(r.FormValue("has_gas"))
	req.HasInternet = parseBool(r.FormValue("has_internet"))
	req.HasGarden = parseBool(r.FormValue("has_garden"))
	req.HasPool = parseBool(r.FormValue("has_pool"))
	req.HasElevator = parseBool(r.FormValue("has_elevator"))
	req.HasCentralHeating = parseBool(r.FormValue("has_central_heating"))
	req.HasWaterTank = parseBool(r.FormValue("has_water_tank"))
	req.HasAirConditioner = parseBool(r.FormValue("has_air_conditioner"))
	req.HasEquippedKitchen = parseBool(r.FormValue("has_equipped_kitchen"))
	req.HasTerrace = parseBool(r.FormValue("has_terrace"))
	req.HasNotarialDeed = parseBool(r.FormValue("has_notarial_deed"))
	req.HasLandBooklet = parseBool(r.FormValue("has_land_booklet"))
	req.HasActInJointOwnership = parseBool(r.FormValue("has_act_in_joint_ownership"))
	req.HasCertificateOfConformity = parseBool(r.FormValue("has_certificate_of_conformity"))
	req.HasDecision = parseBool(r.FormValue("has_decision"))
	req.HasConcession = parseBool(r.FormValue("has_concession"))
	req.HasStampedPaper = parseBool(r.FormValue("has_stamped_paper"))
	req.HasBuildingPermit = parseBool(r.FormValue("has_building_permit"))
	req.HasOffPlanSalesContract = parseBool(r.FormValue("has_off_plan_sales_contract"))
	req.Furnished = parseBool(r.FormValue("furnished"))

	return req, nil, nil
}
