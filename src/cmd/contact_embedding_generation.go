package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"logispro/internal/config"
	"logispro/internal/db"
	"logispro/internal/sqlite"
	"logispro/internal/utils"
	"net/http"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ContactEmbedding struct {
	ID     int64
	Params db.CreateContactParams
}
type OllamaContactResponse struct {
	Embedding []float64 `json:"embedding"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("sleeping 10s")
	time.Sleep(time.Second * 10)
	config.LoadEnv()

	sqliteDb, err := sqlite.New("file", config.SqlitePath)
	if err != nil {
		panic(err)
	}
	defer sqliteDb.Close()
	queries := db.New(sqliteDb.GetDB())

	var rabbitMqConn *amqp.Connection
	for i := 1; i <= 10; i++ {
		rabbitMqConn, err = amqp.Dial(config.RabbitMqHost)
		if err == nil {
			break
		}
		logger.Error("Failed to connect to RabbitMQ", slog.Any("attempt", i), slog.Any("error", err))
		time.Sleep(time.Duration(i) * 2 * time.Second)
	}
	if err != nil {
		logger.Error("Could not connect to RabbitMQ after retries", slog.Any("error", err))
	}
	defer rabbitMqConn.Close()
	ch, err := rabbitMqConn.Channel()
	if err != nil {
		logger.Error("Failed to open channel", slog.Any("error", err))
	}
	defer ch.Close()
	rmq := &utils.RabbitMQ{Conn: rabbitMqConn, Channel: ch}
	err = rmq.DeclareQueue("created_contacts")
	if err != nil {
		logger.Error("Failed to declare queue created_contacts", slog.Any("error", err))
	}
	msgs, err := rmq.Channel.Consume("created_contacts", "", true, false, false, false, nil)
	if err != nil {
		logger.Error("Failed to consume queue created_contacts", slog.Any("error", err))
	}
	logger.Info("listening for created_contacts...")

	var u ContactEmbedding
	for msg := range msgs {

		err := json.Unmarshal(msg.Body, &u)
		retryCount := rmq.GetRetryCount(msg.Headers)
		if err != nil {
			logger.Error("failed to unmarshal message", slog.Any("error", err))
			rmq.Publish("created_contacts", msg.Body, amqp.Table{"x-retry": retryCount + 1})
			continue
		}

		if retryCount > 3 {
			logger.Error("max attempts exceded", slog.Any("contact id", u.ID))
			continue
		}

		var parts []string

		appendStr := func(label string, v sql.NullString) {
			if v.Valid && v.String != "" {
				parts = append(parts, fmt.Sprintf("%s: %s", label, v.String))
			}
		}
		appendBool := func(label string, v sql.NullBool) {
			if v.Valid {
				parts = append(parts, fmt.Sprintf("%s: %t", label, v.Bool))
			}
		}
		appendInt := func(label string, v sql.NullInt64) {
			if v.Valid {
				parts = append(parts, fmt.Sprintf("%s: %d", label, v.Int64))
			}
		}
		appendFloat := func(label string, v sql.NullFloat64) {
			if v.Valid {
				parts = append(parts, fmt.Sprintf("%s: %.2f", label, v.Float64))
			}
		}

		// Core matching fields
		appendStr("Wilaya", u.Params.Wilaya)
		appendStr("Daira", u.Params.Daira)
		appendStr("BuildingFinishing", u.Params.HouseFinishing)
		appendStr("Payment", u.Params.AcceptablePaymentType)
		appendBool("Furnished", u.Params.Furnished)
		appendInt("RoomsMin", u.Params.MinRooms)
		appendInt("RoomsMax", u.Params.MaxRooms)
		appendFloat("SurfaceMin", u.Params.MinSurface)
		appendFloat("SurfaceMax", u.Params.MaxSurface)
		appendInt("BudgetMin", u.Params.MinBudget)
		appendInt("BudgetMax", u.Params.MaxBudget)
		appendInt("MaxYearBuilt", u.Params.MaxYearBuilt)

		// Normalize features
		if u.Params.PreferredFeatures.Valid && u.Params.PreferredFeatures.String != "" {
			var features []string
			_ = json.Unmarshal([]byte(u.Params.PreferredFeatures.String), &features)
			for _, f := range features {
				parts = append(parts, fmt.Sprintf("Feature: %s", f))
			}
		}

		if u.Params.PreferredBuildingTypes.Valid && u.Params.PreferredBuildingTypes.String != "" {
			var buildingTypes []string
			_ = json.Unmarshal([]byte(u.Params.PreferredBuildingTypes.String), &buildingTypes)
			for _, b := range buildingTypes {
				parts = append(parts, fmt.Sprintf("BuildingType: %s", b))
			}
		}

		contactEmbeddingText := strings.Join(parts, ", ")

		payload := map[string]string{
			"model":  "mxbai-embed-large",
			"prompt": contactEmbeddingText,
		}
		jsonData, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/embeddings", config.OllamaHost), bytes.NewBuffer(jsonData))
		if err != nil {
			logger.Error("failed to init request to ollama server", slog.Any("error", err), slog.Any("contact id", u.ID))
			rmq.Publish("created_contacts", msg.Body, amqp.Table{"x-retry": retryCount + 1})
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error("failed to send request to ollama server", slog.Any("error", err), slog.Any("contact id", u.ID))
			rmq.Publish("created_contacts", msg.Body, amqp.Table{"x-retry": retryCount + 1})
			continue
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var OllamaContactResponse OllamaContactResponse
		json.Unmarshal(body, &OllamaContactResponse)
		jsonEmbedding, _ := json.Marshal(OllamaContactResponse.Embedding)
		err = queries.InsertContactEmbeddings(context.Background(), db.InsertContactEmbeddingsParams{
			ContactID: u.ID,
			Embedding: string(jsonEmbedding),
		})
		if err != nil {
			logger.Error("failed to insert embedding to db", slog.Any("error", err), slog.Any("contact id", u.ID))
			rmq.Publish("created_contacts", msg.Body, amqp.Table{"x-retry": retryCount + 1})
		}
		logger.Info("created for contact ", slog.Any("id", u.ID))
	}
}
