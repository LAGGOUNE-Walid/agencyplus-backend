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

func LoadEnv() {
	_ = godotenv.Load(".env")
	secret := os.Getenv("JWT_SECRET")
	sqlitePath := os.Getenv("SQLITE_PATH")
	rabbitMqHost := os.Getenv("RABBITMQ_URL")
	ollamaHost := os.Getenv("OLLAMA_URL")
	gorseHost := os.Getenv("GORSE_URL")
	gorseApiKey := os.Getenv("GORSE_API_KEY")
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
	RabbitMqHost = rabbitMqHost
	OllamaHost = ollamaHost
	GorseHost = gorseHost
	GorseApiKey = gorseApiKey
}
