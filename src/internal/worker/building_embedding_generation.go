package worker

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

type BuildingEmbedding struct {
	ID     int64
	Params db.CreateBuildingParams
}
type OllamaEmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

func StartBuildingEmbdGenerationWorker() {
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
	err = rmq.DeclareQueue("created_buildings")
	if err != nil {
		logger.Error("Failed to declare queue created_buildings", slog.Any("error", err))
	}
	msgs, err := rmq.Channel.Consume("created_buildings", "", true, false, false, false, nil)
	if err != nil {
		logger.Error("Failed to consume queue created_buildings", slog.Any("error", err))
	}
	logger.Info("listening for created_buildings...")

	var b BuildingEmbedding
	for msg := range msgs {

		err := json.Unmarshal(msg.Body, &b)
		retryCount := rmq.GetRetryCount(msg.Headers)
		if err != nil {
			logger.Error("failed to unmarshal message", slog.Any("error", err))
			rmq.Publish("created_buildings", msg.Body, amqp.Table{"x-retry": retryCount + 1})
			continue
		}

		if retryCount > 3 {
			logger.Error("max attempts exceded", slog.Any("building id", b.ID))
			continue
		}

		var parts []string

		appendStr := func(label string, v sql.NullString) {
			if v.Valid && v.String != "" {
				parts = append(parts, fmt.Sprintf("%s: %s", label, v.String))
			}
		}
		appendBool := func(label string, v sql.NullBool) {
			if v.Valid && v.Bool {
				parts = append(parts, fmt.Sprintf("Feature: %s", label))
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

		// Matching attributes
		appendStr("Wilaya", b.Params.Wilaya)
		appendStr("Daira", b.Params.Daira)
		appendStr("BuildingType", b.Params.BuildingType)
		appendStr("BuildingFinishing", b.Params.BuildingFinishedType)
		appendStr("Payment", b.Params.AcceptablePaymentType)
		appendBool("Furnished", b.Params.Furnished)
		appendInt("Rooms", b.Params.Rooms)
		appendFloat("SurfaceTotal", b.Params.SurfaceTotal)
		appendInt("Price", b.Params.Price)
		appendInt("YearBuilt", b.Params.YearBuilt)

		// Feature flags
		appendBool("has_water", b.Params.HasWater)
		appendBool("has_electricity", b.Params.HasElectricity)
		appendBool("has_gas", b.Params.HasGas)
		appendBool("has_internet", b.Params.HasInternet)
		appendBool("has_garden", b.Params.HasGarden)
		appendBool("has_pool", b.Params.HasPool)
		appendBool("has_elevator", b.Params.HasElevator)
		appendBool("has_central_heating", b.Params.HasCentralHeating)
		appendBool("has_water_tank", b.Params.HasWaterTank)
		appendBool("has_air_conditioner", b.Params.HasAirConditioner)
		appendBool("has_equipped_kitchen", b.Params.HasEquippedKitchen)
		appendBool("has_terrace", b.Params.HasTerrace)
		appendBool("has_notarial_deed", b.Params.HasNotarialDeed)
		appendBool("has_land_booklet", b.Params.HasLandBooklet)
		appendBool("has_act_in_joint_ownership", b.Params.HasActInJointOwnership)
		appendBool("has_certificate_of_conformity", b.Params.HasCertificateOfConformity)
		appendBool("has_decision", b.Params.HasDecision)
		appendBool("has_concession", b.Params.HasConcession)
		appendBool("has_stamped_paper", b.Params.HasStampedPaper)
		appendBool("has_building_permit", b.Params.HasBuildingPermit)
		appendBool("has_off_plan_sales_contract", b.Params.HasOffPlanSalesContract)

		buildingEmbeddingText := strings.Join(parts, ", ")

		payload := map[string]string{
			"model":  "mxbai-embed-large",
			"prompt": buildingEmbeddingText,
		}

		jsonData, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/embeddings", config.OllamaHost), bytes.NewBuffer(jsonData))
		if err != nil {
			logger.Error("failed to init request to ollama server", slog.Any("error", err), slog.Any("building id", b.ID))
			rmq.Publish("created_buildings", msg.Body, amqp.Table{"x-retry": retryCount + 1})
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error("failed to send request to ollama server", slog.Any("error", err), slog.Any("building id", b.ID))
			rmq.Publish("created_buildings", msg.Body, amqp.Table{"x-retry": retryCount + 1})
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		var OllamaEmbeddingResponse OllamaEmbeddingResponse
		json.Unmarshal(body, &OllamaEmbeddingResponse)
		jsonEmbedding, _ := json.Marshal(OllamaEmbeddingResponse.Embedding)

		err = queries.InsertEmbeddings(context.Background(), db.InsertEmbeddingsParams{
			BuildingID: b.ID,
			Embedding:  string(jsonEmbedding),
		})
		if err != nil {
			logger.Error("failed to insert embedding to db", slog.Any("error", err), slog.Any("building id", b.ID))
			rmq.Publish("created_buildings", msg.Body, amqp.Table{"x-retry": retryCount + 1})
			continue
		}
		logger.Info("created for building ", slog.Any("id", b.ID))
	}
}
