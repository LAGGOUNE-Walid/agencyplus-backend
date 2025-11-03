package worker

import (
	"context"
	"encoding/json"
	"log/slog"
	"logispro/internal/config"
	"logispro/internal/db"
	pdfservice "logispro/internal/services/pdf_service"
	"logispro/internal/sqlite"
	"logispro/internal/utils"
	"os"
	"path/filepath"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/mail.v2"
)

func StartEmailPaymentInvoicesWorker() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("sleeping 10s")
	// time.Sleep(time.Second * 10)
	config.LoadEnv()
	sqliteDb, err := sqlite.New("file", config.SqlitePath)
	if err != nil {
		panic(err)
	}
	defer sqliteDb.Close()
	queries := db.New(sqliteDb.GetDB())
	dialer := mail.NewDialer(config.MailHost, config.MailPort, config.MailUsername, config.MailPassword)
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
	err = rmq.DeclareQueue("created_invoices")
	if err != nil {
		logger.Error("Failed to declare queue created_invoices", slog.Any("error", err))
	}
	msgs, err := rmq.Channel.Consume("created_invoices", "", true, false, false, false, nil)
	if err != nil {
		logger.Error("Failed to consume queue created_invoices", slog.Any("error", err))
	}
	logger.Info("listening for created_invoices...")
	var generatedInv pdfservice.GeneratedInvoice
	for msg := range msgs {
		err := json.Unmarshal(msg.Body, &generatedInv)
		retryCount := rmq.GetRetryCount(msg.Headers)
		if err != nil {
			logger.Error("failed to unmarshal message", slog.Any("error", err))
			rmq.Publish("created_invoices", msg.Body, amqp.Table{"x-retry": retryCount + 1})
			continue
		}
		if retryCount > 3 {
			logger.Error("max attempts exceded", slog.Any("invoice id", generatedInv.Id))
			continue
		}
		ctx := context.Background()
		agencyUsers, err := queries.GetUsers(ctx, generatedInv.Users)
		if err != nil {
			logger.Error("filed to get users from database", slog.Any("users", generatedInv.Users))
			continue
		}
		var wg sync.WaitGroup
		throttle := time.Tick(1 * time.Second)
		for _, user := range agencyUsers {

			wg.Add(1)
			go func(user db.User, dilaer *mail.Dialer, wg *sync.WaitGroup, throttle <-chan time.Time) {
				defer wg.Done()
				<-throttle
				message := mail.NewMessage()
				message.SetHeader("From", "payment@logispro.com")
				// message.SetHeader("To", user.Email)
				message.SetHeader("To", "walidlaggoune159@gmail.com")
				message.SetHeader("Subject", "Your monthly payment")
				templatesAbsPath, _ := filepath.Abs("templates/email/")
				template, err := utils.ParseTemplate(templatesAbsPath+"/user_invoice.html", user)
				if err != nil {
					logger.Error("failed to parse user invoice template", slog.Any("err", err))
				}
				message.SetBody("text/html", string(template))
				message.Attach(generatedInv.Path)
				if err := dialer.DialAndSend(message); err != nil {
					logger.Error("error sending invoice email to user", slog.Any("user", user.Email), slog.Any("err", err))
				} else {
					logger.Info("sent invoice email to user", slog.Any("user", user.Email))
				}
			}(user, dialer, &wg, throttle)
		}
		wg.Wait()

	}
}
