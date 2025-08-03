package user

import (
	"fmt"
	"logispro/internal/constants"
	"logispro/internal/services/user_services"
	"logispro/internal/shared/response_types"
	"logispro/internal/utils"
	"logispro/internal/web/requests"
	"net/http"
)

type UserController struct {
	CreateUserService *user_services.CreateUserService
	AuthService       *user_services.AuthService
	UpdateUserService *user_services.UpdateUserService
}

func (c *UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	req, validationErrors, err := requests.ParseCreateUserRequest(r, c.CreateUserService.Queries, r.Context())
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
	userId, jwtToken, err := c.CreateUserService.Create(req)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return response_types.ApiResponse{
		Content: map[string]any{
			"user":  userId,
			"token": jwtToken,
		},
		StatusCode: http.StatusCreated,
		Error:      nil,
	}
}

func (c *UserController) Auth(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	req, validationErrors := requests.ParseAuthRequest(r)
	if len(validationErrors) > 0 {
		return response_types.ApiResponse{
			Content:    validationErrors,
			StatusCode: http.StatusBadRequest,
		}
	}
	user, token, err := c.AuthService.Authenticate(r.Context(), req)
	if err != nil {
		return response_types.ApiResponse{
			Content:    err.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}
	return response_types.ApiResponse{
		Content: map[string]any{
			"user": map[string]any{
				"id":       user.ID,
				"fullname": user.Fullname,
				"email":    user.Email,
				"role":     user.Role,
				"rootId":   user.RootID,
			},
			"token": token,
		},
		StatusCode: http.StatusOK,
	}
}

func (c *UserController) GetAgencyUsers(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	ctx := r.Context()
	userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
	if !ok {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
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
	agencyUsers, err := utils.GetAgencyUsers(r.Context(), c.CreateUserService.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    agencyUsers,
		StatusCode: http.StatusOK,
	}
}

func (c *UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	req, validationErrors, err := requests.ParseUpdateUserRequest(r, c.UpdateUserService.Queries, r.Context())
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

	err = c.UpdateUserService.Update(r.Context(), req)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return response_types.ApiResponse{
		Content:    "user updated successfully",
		StatusCode: http.StatusOK,
	}
}
