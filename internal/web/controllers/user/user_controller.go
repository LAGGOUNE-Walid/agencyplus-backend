package user

import (
	"logispro/internal/services/user_services"
	"logispro/internal/shared/response_types"
	"logispro/internal/web/requests"
	"net/http"
)

type UserController struct {
	CreateUserService *user_services.CreateUserService
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

func (c *UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	return response_types.ApiResponse{}
}
