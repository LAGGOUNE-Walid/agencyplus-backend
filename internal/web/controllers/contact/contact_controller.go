package contact

import (
	"logispro/internal/services/contact_service"
	"logispro/internal/shared/response_types"
	"logispro/internal/web/requests"
	"net/http"
)

type ContactController struct {
	CreateContactService *contact_service.CreateContactService
}

func (c *ContactController) CreateContactHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	contactService := c.CreateContactService
	req, validationErrors, err := requests.ParseCreateContactRequest(r, contactService.Queries, r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}
	if len(validationErrors) > 0 {
		return response_types.ApiResponse{
			Content:    validationErrors,
			StatusCode: http.StatusBadRequest,
		}
	}
	contactId, err := contactService.Create(r.Context(), req)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content: map[string]any{
			"contact_id": contactId,
		},
		StatusCode: http.StatusCreated,
		Error:      nil,
	}
}
