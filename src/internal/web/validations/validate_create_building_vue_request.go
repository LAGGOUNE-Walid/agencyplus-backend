package validations

import (
	"net/http"
)

func ValidateCreateBuildingVueRequest(r *http.Request) ValidationErrors {
	errs := make(ValidationErrors)
	ipAddress := r.FormValue("ip_address")
	userAgent := r.FormValue("user_agent")

	ValidateNonEmpty(ipAddress, "ip_address", "requis", errs)
	ValidateIp(ipAddress, "ip_address", "IP valide", errs)
	ValidateNonEmpty(userAgent, "user_agent", "requis", errs)

	return errs
}
