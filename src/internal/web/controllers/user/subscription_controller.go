package user

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"logispro/internal/config"
	"logispro/internal/constants"
	"logispro/internal/db"
	"logispro/internal/services/payment_service"
	"logispro/internal/shared/response_types"
	"logispro/internal/utils"
	"net/http"
	"time"
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

func (s *SubscriptionController) Cancel(w http.ResponseWriter, r *http.Request) response_types.Responder {
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
	err = s.SubscriptionService.UpdateUserSubscriptionStatus(ctx, payment_service.SUBS_STATUS_CANCELLED, agencyUsersId)
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

func (s *SubscriptionController) ChargilyWebhook(w http.ResponseWriter, r *http.Request) response_types.Responder {
	ctx := r.Context()
	signature := r.Header.Get("signature")
	if signature == "" {
		return response_types.ApiResponse{
			Content:    nil,
			Error:      fmt.Errorf("singature required"),
			StatusCode: http.StatusBadRequest,
		}
	}
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return response_types.ApiResponse{
			Content:    nil,
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	computedSignature := computeHMAC(payload, config.ChargiliySecretKey)
	if !hmac.Equal([]byte(computedSignature), []byte(signature)) {
		return response_types.ApiResponse{
			Content:    nil,
			Error:      fmt.Errorf("invalid signature"),
			StatusCode: http.StatusForbidden,
		}
	}
	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return response_types.ApiResponse{
			Content:    nil,
			Error:      fmt.Errorf("error deconding json payload"),
			StatusCode: http.StatusInternalServerError,
		}
	}
	switch event["type"] {
	case "checkout.paid":
		checkout := event["data"].(map[string]interface{})
		metadata, ok := checkout["metadata"].(map[string]interface{})
		if !ok {
			return response_types.ApiResponse{
				Content:    nil,
				Error:      fmt.Errorf("faild to get metadata"),
				StatusCode: http.StatusInternalServerError,
			}
		}
		amount, ok := checkout["amount"].(float64)
		if !ok {
			return response_types.ApiResponse{
				Content:    nil,
				Error:      fmt.Errorf("faild to decode amount to float64"),
				StatusCode: http.StatusInternalServerError,
			}
		}
		usersInterface, ok := metadata["users"].([]interface{})
		if !ok {
			return response_types.ApiResponse{
				Content:    nil,
				Error:      fmt.Errorf("faild to get users ids"),
				StatusCode: http.StatusInternalServerError,
			}
		}
		var usersSlice []int64
		for _, v := range usersInterface {
			if f, ok := v.(float64); ok {
				usersSlice = append(usersSlice, int64(f))
			}
		}
		for _, userID := range usersSlice {
			subscription := payment_service.Subscription{
				UserId:             userID,
				PlanId:             payment_service.PLAN_MONTH,
				Status:             payment_service.SUBS_STATUS_ACTIVE,
				CurrentPeriodStart: time.Now(),
				CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0),
				NextBillingDate:    time.Now().AddDate(0, 1, 0),
				TrialStart:         time.Time{},
				TrialEnd:           time.Time{},
				Ammount:            amount,
			}
			s.SubscriptionService.CreateSubscription(ctx, subscription)
			s.PaymentService.SavePaymentPayload(ctx, userID, string(payload))
		}
	}
	return response_types.ApiResponse{
		StatusCode: http.StatusOK,
		Content:    nil,
	}
}

func computeHMAC(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
