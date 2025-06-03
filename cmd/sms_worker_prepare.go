package main

import (
	"context"
	"encoding/json"
	"log"
	"logispro/internal/config"
	"logispro/internal/db"
	"logispro/internal/sqlite"
	"logispro/internal/utils"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SMSPreparePayload struct {
	SMSID int64 `json:"sms_id"`
}

type Contact struct {
	ID     int64
	Number string
}

func main() {
	time.Sleep(time.Second * 10)
	config.LoadEnv()
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	sqliteDb, err := sqlite.New("file", config.SqlitePath)
	queries := db.New(sqliteDb.GetDB())
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

	err = rmq.DeclareQueue("sms_prepare")
	if err != nil {
		log.Fatalf("❌ Failed to declare queue sms_prepare: %v", err)
	}
	err = rmq.DeclareQueue("sms_send")
	if err != nil {
		log.Fatalf("❌ Failed to declare queue sms_send: %v", err)
	}

	msgs, err := rmq.Channel.Consume("sms_prepare", "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ Failed to consume queue sms_prepare: %v", err)
	}
	ctx := context.Background()
	log.Println("✅ Worker A listening for sms_prepare...")
	for msg := range msgs {
		var payload SMSPreparePayload
		err := json.Unmarshal(msg.Body, &payload)
		if err != nil {
			log.Printf("⚠️ Failed to unmarshal message: %v", err)
			continue
		}

		smsQueue, err := queries.GetSmsQueue(ctx, payload.SMSID)
		if err != nil {
			log.Println("A")
			log.Println(payload.SMSID)
			log.Println(err)
		} else {
			contacts, err := queries.GetSmsContacts(ctx, payload.SMSID)
			if err != nil {
				log.Println("B")
				log.Println(err)
			} else {
				for _, c := range contacts {
					msg := map[string]interface{}{
						"sms_id":     payload.SMSID,
						"contact_id": c.ID,
						"number":     c.PhoneNumber,
						"Content":    smsQueue.Content,
					}
					body, _ := json.Marshal(msg)
					rmq.Publish("sms_send", body)
				}
			}
		}

	}
}
