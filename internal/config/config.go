package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JwtSecret []byte
var SqlitePath string
var RabbitMqHost string

func LoadEnv() {
	_ = godotenv.Load(".env")
	secret := os.Getenv("JWT_SECRET")
	sqlitePath := os.Getenv("SQLITE_PATH")
	rabbitMqHost := os.Getenv("RABBITMQ_URL")
	if secret == "" {
		log.Fatal("JWT_SECRET not set in environment")
	}
	if sqlitePath == "" {
		log.Fatal("SQLITE_PATH not set in environment")
	}
	if rabbitMqHost == "" {
		log.Fatal("RABBITMQ_URL not set in environment")
	}
	JwtSecret = []byte(secret)
	SqlitePath = sqlitePath
	RabbitMqHost = rabbitMqHost
}
