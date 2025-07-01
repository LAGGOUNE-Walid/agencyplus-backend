package main

import (
	"encoding/json"
	"log"
	"logispro/internal/config"
	"logispro/internal/sqlite"
	"logispro/internal/utils"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SMSPayload struct {
	SMSID     int64  `json:"sms_id"`
	ContactID int64  `json:"contact_id"`
	Number    string `json:"number"`
	Content   string `json:"content"`
}

func sendSMS(to string, content string) error {
	log.Println("Sending SMS to:", to)
	return nil // simulate send
}

func main() {
	log.Printf("Sleeping 10s")
	time.Sleep(time.Second * 10)
	config.LoadEnv()
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	sqliteDb, err := sqlite.New("file", config.SqlitePath)
	// queries := db.New(sqliteDb.GetDB())
	if err != nil {
		panic(err)
	}
	defer sqliteDb.Close()

	var rabbitMqConn *amqp.Connection

	// Retry loop for RabbitMQ connection
	for i := 1; i <= 10; i++ {
		rabbitMqConn, err = amqp.Dial(config.RabbitMqHost)
		if err == nil {
			break
		}
		log.Printf("❌ Failed to connect to RabbitMQ (attempt %d): %v", i, err)
		time.Sleep(time.Duration(i) * 2 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ Could not connect to RabbitMQ after retries: %v", err)
	}
	defer rabbitMqConn.Close()

	ch, err := rabbitMqConn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open channel: %v", err)
	}
	defer ch.Close()

	rmq := &utils.RabbitMQ{Conn: rabbitMqConn, Channel: ch}

	err = rmq.DeclareQueue("sms_send")
	if err != nil {
		log.Fatalf("❌ Failed to declare queue sms_send: %v", err)
	}
	msgs, err := rmq.Channel.Consume("sms_send", "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ Failed to consume queue sms_send: %v", err)
	}
	// ctx := context.Background()
	log.Println("✅ Worker A listening for sms_prepare...")
	for msg := range msgs {
		var payload SMSPayload
		err := json.Unmarshal(msg.Body, &payload)
		if err != nil {
			log.Printf("⚠️ Failed to unmarshal message: %v", err)
			continue
		}
		sendSMS(payload.Number, payload.Content)
	}
}
