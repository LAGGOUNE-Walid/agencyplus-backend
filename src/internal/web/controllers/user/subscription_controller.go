package user

import (
	"context"
	"fmt"
	"logispro/internal/config"
	"logispro/internal/constants"
	"logispro/internal/db"
	"logispro/internal/services/payment_service"
	"logispro/internal/shared/response_types"
	"logispro/internal/utils"
	"net/http"
)

type SubscriptionController struct {
	SubscriptionService *payment_service.SubscriptionService
	PaymentService      *payment_service.PaymentService
}

func (s *SubscriptionController) GetStatus(ctx context.Context, userId int64) (payment_service.Status, error) {
	return s.SubscriptionService.GetSubscriptionStatus(ctx, userId)
}

func (s *SubscriptionController) CreateCheckoutLink(w http.ResponseWriter, r *http.Request) response_types.Responder {
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
	agencyUsers, err := utils.GetAgencyUsers(r.Context(), s.SubscriptionService.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsersId := utils.ExtractField(agencyUsers, func(u db.GetAgencyUsersRow) int64 {
		return u.ID
	})
	checkoutData, err := s.PaymentService.CreateCheckoutLink(agencyUsersId, config.MonthlyPriceId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    checkoutData,
		StatusCode: http.StatusCreated,
	}
}
