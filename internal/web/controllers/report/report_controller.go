package report

import (
	"fmt"
	"logispro/internal/constants"
	"logispro/internal/services/report_service"
	"logispro/internal/shared/response_types"
	"logispro/internal/web/requests"
	"net/http"
	"strconv"
)

type ReportController struct {
	ReportService *report_service.ReportService
}

func (c *ReportController) CreateReportHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	userId, ok := r.Context().Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", r.Context().Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	req, validationErrors := requests.ParseCreateReportRequest(r)
	if len(validationErrors) > 0 {
		return response_types.ApiResponse{
			Content:    validationErrors,
			StatusCode: http.StatusBadRequest,
		}
	}
	req.UserID = userId
	report, err := c.ReportService.Create(req, r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Content:    nil,
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    report,
		StatusCode: http.StatusCreated,
	}
}

func (c *ReportController) UpdateReportHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	userId, ok := r.Context().Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", r.Context().Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("invalid building ID"),
		}
	}
	req, validationErrors := requests.ParseUpdateReportRequest(r)
	if len(validationErrors) > 0 {
		return response_types.ApiResponse{
			Content:    validationErrors,
			StatusCode: http.StatusBadRequest,
		}
	}
	req.UserID = userId
	req.ID = id
	err = c.ReportService.Update(req, r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Content:    nil,
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    nil,
		StatusCode: http.StatusNoContent,
	}
}

func (c *ReportController) DeleteReportHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	userId, ok := r.Context().Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", r.Context().Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("invalid building ID"),
		}
	}
	err = c.ReportService.Delete(r.Context(), id, userId)
	if err != nil {
		return response_types.ApiResponse{
			Content:    nil,
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    nil,
		StatusCode: http.StatusNoContent,
	}
}

func (c *ReportController) GetReportsHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	userId, ok := r.Context().Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", r.Context().Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	reports, err := c.ReportService.All(r.Context(), userId)
	if err != nil {
		return response_types.ApiResponse{
			Content:    nil,
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		StatusCode: http.StatusOK,
		Content:    reports,
	}
}
