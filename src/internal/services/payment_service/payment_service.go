package payment_service

import (
	"context"
	"fmt"
	"logispro/internal/config"
	"logispro/internal/db"

	"github.com/Chargily/chargily-pay-go/pkg/chargily"
	"github.com/Chargily/chargily-pay-go/pkg/models"
)

type PaymentService struct {
	Client  *chargily.Client
	Queries *db.Queries
}

func (p *PaymentService) CreateCheckoutLink(users []int64, priceId string) (*models.Checkout, error) {
	items := []models.CItems{
		{
			Price:    priceId,
			Quantity: 1,
		},
	}
	checkout := &models.CheckoutParams{
		Items:           items,
		PaymentMethod:   "edahabia",
		SuccessURL:      fmt.Sprintf("%s/%s", config.AppUrl, "chargily-webhhook"),
		FailureURL:      fmt.Sprintf("%s/%s", config.AppUrl, "chargily-webhhook"),
		WebhookEndpoint: fmt.Sprintf("%s/%s", config.AppUrl, "chargily-webhhook"),
		Description:     "Checkout for Order #12345",
		Locale:          "en",
		Metadata: map[string]any{
			"users": users,
		},
	}
	checkoutData, err := p.Client.Checkouts.Create(checkout)
	if err != nil {
		return nil, err
	}
	return checkoutData, nil
}

func (p *PaymentService) SavePaymentPayload(ctx context.Context, userId int64, payload string) error {
	return p.Queries.CreatePayment(ctx, db.CreatePaymentParams{
		UserID:  userId,
		Payload: payload,
	})
}
