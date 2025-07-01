package calendar

import (
	"fmt"
	"logispro/internal/constants"
	"logispro/internal/services/calendar_service"
	"logispro/internal/shared/response_types"
	"logispro/internal/web/requests"
	"net/http"
	"strconv"
)

type CalendarController struct {
	CalendarService *calendar_service.CalendarService
}

func (c *CalendarController) CreateCalendarEventHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	userId, ok := r.Context().Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", r.Context().Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	req, validationErrors := requests.ParseCreateCalendarEventRequest(r)
	if len(validationErrors) > 0 {
		return response_types.ApiResponse{
			Content:    validationErrors,
			StatusCode: http.StatusBadRequest,
		}
	}
	req.UserId = userId
	report, err := c.CalendarService.Create(r.Context(), req)
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

func (c *CalendarController) DeleteCalendarEventHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
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
	err = c.CalendarService.Delete(r.Context(), id, userId)
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

func (c *CalendarController) GetCalendarEventsHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	userId, ok := r.Context().Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", r.Context().Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	reports, err := c.CalendarService.All(r.Context(), userId)
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
