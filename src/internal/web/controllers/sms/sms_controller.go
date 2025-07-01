package sms

import (
	"fmt"
	"logispro/internal/constants"
	"logispro/internal/services/sms_service"
	"logispro/internal/shared/response_types"
	"logispro/internal/web/requests"
	"net/http"
)

type SmsController struct {
	CreateSmsService *sms_service.CreateSmsService
}

func (c *SmsController) CreateSmsHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	userId, ok := r.Context().Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", r.Context().Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	req, validationErrors := requests.ParseCreateSmsRequest(r)
	if len(validationErrors) > 0 {
		return response_types.ApiResponse{
			Content:    validationErrors,
			StatusCode: http.StatusBadRequest,
		}
	}
	_, err := c.CreateSmsService.Create(req, userId, r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Content:    nil,
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    nil,
		StatusCode: http.StatusCreated,
	}
}
