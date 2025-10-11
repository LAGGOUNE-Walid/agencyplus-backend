package payment_service

import (
	"fmt"
	"logispro/internal/config"

	"github.com/Chargily/chargily-pay-go/pkg/chargily"
	"github.com/Chargily/chargily-pay-go/pkg/models"
)

type PaymentService struct {
	Client *chargily.Client
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
		SuccessURL:      fmt.Sprintf("%s/%s", config.AppUrl, "checkout-success"),
		FailureURL:      fmt.Sprintf("%s/%s", config.AppUrl, "checkout-fail"),
		WebhookEndpoint: fmt.Sprintf("%s/%s", config.AppUrl, "checkout-webhook"),
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
