package building

import (
	"database/sql"
	"errors"
	"fmt"
	"logispro/internal/constants"
	"logispro/internal/services/building_service"
	"logispro/internal/shared/response_types"
	"logispro/internal/utils"
	"logispro/internal/web/requests"
	"net/http"
	"strconv"
)

type BuildingController struct {
	CreateBuildingService         *building_service.CreateBuildingService
	GetBuildingService            *building_service.GetBuildingService
	UpdateBuildingService         *building_service.UpdateBuildingService
	GetBuildingsStatisticsService *building_service.GetBuildingsStatisticsService
}

func (c *BuildingController) CreateBuildingHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	ctx := r.Context()
	userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", r.Context().Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	req, validationErrors, err := requests.ParseCreateBuildingRequest(r, c.CreateBuildingService.Queries, ctx, userId)
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
	buildingId, err := c.CreateBuildingService.Create(ctx, req, req.ImageHeaders, req.DocumentHeaders)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content: map[string]any{
			"building": buildingId,
		},
		StatusCode: http.StatusCreated,
		Error:      nil,
	}
}

func (c *BuildingController) GetBuildingsHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	userId, ok := r.Context().Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", r.Context().Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	pageString := r.URL.Query().Get("page")
	page := utils.ParseInt64(pageString)
	offset := (page - 1) * 20
	buildings, err := c.GetBuildingService.All(userId, offset, 20, r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    buildings,
		StatusCode: http.StatusOK,
	}
}

func (c *BuildingController) GetBuildingHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
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
	b, err := c.GetBuildingService.Get(userId, id, r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response_types.ApiResponse{
				StatusCode: http.StatusNotFound,
				Error:      fmt.Errorf("building not found"),
			}
		}
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		StatusCode: http.StatusOK,
		Content:    b,
	}
}

func (c *BuildingController) UpdateBuildingHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	ctx := r.Context()
	userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
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

	req, validationErrors, err := requests.ParseUpdateBuildingInfoRequest(r, userId)
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

	err = c.UpdateBuildingService.UpdateBasicInfo(ctx, req, id)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return response_types.ApiResponse{
		Content:    nil,
		StatusCode: http.StatusNoContent,
	}
}

func (c *BuildingController) DeleteBuildingHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	ctx := r.Context()
	userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
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
	err = c.UpdateBuildingService.Delete(ctx, userId, id)
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}
	return response_types.ApiResponse{
		StatusCode: http.StatusNoContent,
	}
}

func (c *BuildingController) CreateBuildingImagesHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	ctx := r.Context()
	userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
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
	req, validationErrors, err := requests.ParseUpdateBuildingImagesRequest(r, userId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	if len(validationErrors) > 0 {
		return response_types.ApiResponse{
			Content:    validationErrors,
			StatusCode: http.StatusBadRequest,
		}
	}
	err = c.UpdateBuildingService.AddImages(ctx, req, id)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    nil,
		StatusCode: http.StatusNoContent,
	}
}

func (c *BuildingController) DeleteBuildingImageHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	ctx := r.Context()
	userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	buildingId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("invalid building ID"),
		}
	}
	imageId, err := strconv.ParseInt(r.PathValue("imageId"), 10, 64)
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("invalid building imageId"),
		}
	}
	err = c.UpdateBuildingService.DeleteImage(ctx, userId, buildingId, imageId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    nil,
		StatusCode: http.StatusNoContent,
	}
}

func (c *BuildingController) CreateBuildingDocumentsHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	ctx := r.Context()
	userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
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
	req, validationErrors, err := requests.ParseUpdateBuildingDocumentsRequest(r, userId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	if len(validationErrors) > 0 {
		return response_types.ApiResponse{
			Content:    validationErrors,
			StatusCode: http.StatusBadRequest,
		}
	}
	err = c.UpdateBuildingService.AddDocuments(ctx, req, id)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    nil,
		StatusCode: http.StatusNoContent,
	}
}

func (c *BuildingController) DeleteBuildingDocumentHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	ctx := r.Context()
	userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	buildingId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("invalid building ID"),
		}
	}
	documentId, err := strconv.ParseInt(r.PathValue("documentId"), 10, 64)
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("invalid building documentId"),
		}
	}
	err = c.UpdateBuildingService.DeleteDocument(ctx, userId, buildingId, documentId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    nil,
		StatusCode: http.StatusNoContent,
	}
}

func (c *BuildingController) AddVueHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	buildingId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("invalid building ID"),
		}
	}
	req, validationErrors := requests.ParseCreateBuildingVueRequest(r, buildingId)
	if len(validationErrors) > 0 {
		return response_types.ApiResponse{
			Content:    validationErrors,
			StatusCode: http.StatusBadRequest,
		}
	}
	err = c.UpdateBuildingService.AddVue(r.Context(), req)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    nil,
		StatusCode: http.StatusCreated,
	}
}

func (c *BuildingController) GetBuildingsStatisticsHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	ctx := r.Context()
	userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	statistics, err := c.GetBuildingsStatisticsService.Get(userId, ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response_types.ApiResponse{
				Content:    nil,
				StatusCode: http.StatusCreated,
			}
		}
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    statistics,
		StatusCode: http.StatusCreated,
	}
}

func (c *BuildingController) GetBuildingsGainHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	ctx := r.Context()
	userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	gain, err := c.GetBuildingsStatisticsService.GetBuildingsTotalChangeRate(userId, ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response_types.ApiResponse{
				Content:    nil,
				StatusCode: http.StatusCreated,
			}
		}
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    gain,
		StatusCode: http.StatusCreated,
	}
}

func (c *BuildingController) GetDairaDistributionHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	ctx := r.Context()
	userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	dairas, err := c.GetBuildingsStatisticsService.GetBuildingsDairaDistribution(userId, ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response_types.ApiResponse{
				Content:    nil,
				StatusCode: http.StatusCreated,
			}
		}
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    dairas,
		StatusCode: http.StatusCreated,
	}
}

func (c *BuildingController) GetMapHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	ctx := r.Context()
	userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	location, err := c.GetBuildingsStatisticsService.GetBuildingsLocations(userId, ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response_types.ApiResponse{
				Content:    nil,
				StatusCode: http.StatusCreated,
			}
		}
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    location,
		StatusCode: http.StatusCreated,
	}
}
