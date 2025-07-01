package requests

import (
	"logispro/internal/web/validations"
	"net"
	"net/http"
)

type CreateBuildingVueRequest struct {
	BuildingId int64
	IpAddress  net.IP
	UserAgent  string
}

func ParseCreateBuildingVueRequest(r *http.Request, buildingId int64) (CreateBuildingVueRequest, validations.ValidationErrors) {
	var req CreateBuildingVueRequest
	req.BuildingId = buildingId
	validationErrors := validations.ValidateCreateBuildingVueRequest(r)
	if len(validationErrors) > 0 {
		return req, validationErrors
	}
	req.IpAddress = net.ParseIP(r.FormValue("ip_address"))
	req.UserAgent = r.FormValue("user_agent")

	return req, nil
}
