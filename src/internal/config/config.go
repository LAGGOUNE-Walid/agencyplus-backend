package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JwtSecret []byte
var SqlitePath string
var RabbitMqHost string
var OllamaHost string
var GorseHost string
var GorseApiKey string
var ChargiliyPublicKey string
var ChargiliySecretKey string
var ProductId string
var MonthlyPriceId string
var PaymentEndpoint string
var AppUrl string

func LoadEnv() {
	_ = godotenv.Load(".env")
	secret := os.Getenv("JWT_SECRET")
	sqlitePath := os.Getenv("SQLITE_PATH")
	rabbitMqHost := os.Getenv("RABBITMQ_URL")
	ollamaHost := os.Getenv("OLLAMA_URL")
	gorseHost := os.Getenv("GORSE_URL")
	gorseApiKey := os.Getenv("GORSE_API_KEY")
	chargiliyPublicKey := os.Getenv("CHARGILIY_PUBLIC_KEY")
	chargiliySecretKey := os.Getenv("CHARGILIY_SECRET_KEY")
	productId := os.Getenv("CHARGILY_PRODUCT_ID")
	monthlyPriceId := os.Getenv("CHARGILY_MONTHLY_PRICE_ID")
	paymentEndpoint := os.Getenv("CHARGILY_ENDPOINT")
	appUrl := os.Getenv("APP_URL")

	if appUrl == "" {
		log.Fatal("APP_URL not set in environment")
	}
	if paymentEndpoint == "" {
		log.Fatal("CHARGILY_ENDPOINT not set in environment")
	}
	if monthlyPriceId == "" {
		log.Fatal("CHARGILY_MONTHLY_PRICE_ID not set in environment")
	}
	if productId == "" {
		log.Fatal("CHARGILY_PRODUCT_ID not set in environment")
	}
	if chargiliyPublicKey == "" {
		log.Fatal("CHARGILIY_PUBLIC_KEY not set in environment")
	}
	if chargiliySecretKey == "" {
		log.Fatal("CHARGILIY_SECRET_KEY not set in environment")
	}
	if secret == "" {
		log.Fatal("JWT_SECRET not set in environment")
	}
	if sqlitePath == "" {
		log.Fatal("SQLITE_PATH not set in environment")
	}
	if rabbitMqHost == "" {
		log.Fatal("RABBITMQ_URL not set in environment")
	}
	if ollamaHost == "" {
		log.Fatal("OLLAMA_URL not set in environment")
	}
	if gorseHost == "" {
		log.Fatal("GORSE_URL not set in environment")
	}
	if gorseApiKey == "" {
		log.Fatal("GORSE_API_KEY not set in environment")
	}
	JwtSecret = []byte(secret)
	SqlitePath = sqlitePath
	SqlitePath = sqlitePath
	RabbitMqHost = rabbitMqHost
	OllamaHost = ollamaHost
	GorseHost = gorseHost
	GorseApiKey = gorseApiKey
	ChargiliyPublicKey = chargiliyPublicKey
	ChargiliySecretKey = chargiliySecretKey
	ProductId = productId
	MonthlyPriceId = monthlyPriceId
	PaymentEndpoint = paymentEndpoint
	AppUrl = appUrl
}
