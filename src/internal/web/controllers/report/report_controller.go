package report

import (
	"fmt"
	"logispro/internal/constants"
	"logispro/internal/db"
	"logispro/internal/services/report_service"
	"logispro/internal/shared/response_types"
	"logispro/internal/utils"
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

	rootId, err := utils.GetRootIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsers, err := utils.GetAgencyUsers(r.Context(), c.ReportService.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsersId := utils.ExtractField(agencyUsers, func(u db.User) int64 {
		return u.ID
	})
	err = c.ReportService.Update(req, r.Context(), rootId, agencyUsersId)
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
	rootId, err := utils.GetRootIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsers, err := utils.GetAgencyUsers(r.Context(), c.ReportService.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsersId := utils.ExtractField(agencyUsers, func(u db.User) int64 {
		return u.ID
	})
	err = c.ReportService.Delete(r.Context(), id, userId, rootId, agencyUsersId)
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
	rootId, err := utils.GetRootIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsers, err := utils.GetAgencyUsers(r.Context(), c.ReportService.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsersId := utils.ExtractField(agencyUsers, func(u db.User) int64 {
		return u.ID
	})
	reports, err := c.ReportService.All(r.Context(), userId, rootId, agencyUsersId)
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
