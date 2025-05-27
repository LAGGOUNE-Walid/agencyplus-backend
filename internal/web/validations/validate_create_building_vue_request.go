package validations

import (
	"net/http"
)

func ValidateCreateBuildingVueRequest(r *http.Request) ValidationErrors {
	errs := make(ValidationErrors)
	ipAddress := r.FormValue("ip_address")
	userAgent := r.FormValue("user_agent")

	ValidateNonEmpty(ipAddress, "ip_address", "required", errs)
	ValidateIp(ipAddress, "ip_address", "valid ip", errs)
	ValidateNonEmpty(userAgent, "user_agent", "required", errs)

	return errs
}
