package main

import (
	"database/sql"
	"log/slog"
	"logispro/internal/config"
	"logispro/internal/db"
	"logispro/internal/services/building_service"
	"logispro/internal/services/calendar_service"
	"logispro/internal/services/contact_service"
	"logispro/internal/services/document_service"
	"logispro/internal/services/payment_service"
	pdfservice "logispro/internal/services/pdf_service"
	"logispro/internal/services/report_service"
	"logispro/internal/services/sms_service"
	"logispro/internal/services/task_service"
	"logispro/internal/services/user_services"
	"logispro/internal/sqlite"
	"logispro/internal/web"
	"logispro/internal/web/controllers"
	"logispro/internal/web/controllers/building"
	"logispro/internal/web/controllers/calendar"
	"logispro/internal/web/controllers/contact"
	"logispro/internal/web/controllers/document"
	"logispro/internal/web/controllers/recommendation"
	"logispro/internal/web/controllers/report"
	"logispro/internal/web/controllers/shareable"
	"logispro/internal/web/controllers/sms"
	"logispro/internal/web/controllers/task"
	"logispro/internal/web/controllers/user"
	"os"
	"path/filepath"

	"github.com/Chargily/chargily-pay-go/pkg/chargily"
	amqp "github.com/rabbitmq/amqp091-go"
)

func InitServices(logger *slog.Logger, db *sql.DB, queries *db.Queries, rabbitMqConn *amqp.Connection, paymentService payment_service.PaymentService) controllers.Controller {
	return controllers.Controller{
		UserController: &user.UserController{
			CreateUserService: &user_services.CreateUserService{
				Queries: queries,
			},
			AuthService: &user_services.AuthService{
				Queries: queries,
			},
			UpdateUserService: &user_services.UpdateUserService{
				Queries: queries,
			},
			SubscriptionService: &payment_service.SubscriptionService{
				Queries: queries,
			},
		},
		ContactController: &contact.ContactController{
			CreateContactService: &contact_service.CreateContactService{
				Queries:      queries,
				RabbitMqConn: rabbitMqConn,
			},
			GetContactService: &contact_service.GetContactService{
				Queries: queries,
			},
			DeleteContactService: &contact_service.DeleteContactService{
				Queries: queries,
			},
		},
		BuildingController: &building.BuildingController{
			CreateBuildingService: &building_service.CreateBuildingService{
				Queries:      queries,
				DB:           db,
				RabbitMqConn: rabbitMqConn,
			},
			GetBuildingService: &building_service.GetBuildingService{
				Queries: queries,
			},
			GetBuildingsStatisticsService: &building_service.GetBuildingsStatisticsService{
				Queries: queries,
			},
			UpdateBuildingService: &building_service.UpdateBuildingService{
				Queries:      queries,
				DB:           db,
				RabbitMqConn: rabbitMqConn,
			},
		},
		SmsController: &sms.SmsController{
			CreateSmsService: &sms_service.CreateSmsService{
				Queries:      queries,
				RabbitMqConn: rabbitMqConn,
			},
		},
		TaskController: &task.TaskController{
			CreateTaskService: &task_service.CreateTaskService{
				Queries: queries,
			},
			GetTasksService: &task_service.GetTasksService{
				Queries: queries,
			},
			UpdateTaskService: &task_service.UpdateTaskService{
				Queries: queries,
			},
		},
		ReportController: &report.ReportController{
			ReportService: &report_service.ReportService{
				Queries: queries,
			},
		},
		CalendarController: &calendar.CalendarController{
			CalendarService: &calendar_service.CalendarService{
				Queries: queries,
			},
		},
		RecommenderController: &recommendation.RecommenderController{
			Queries: queries,
		},
		DocumentController: &document.DocumentController{
			CreateDocumentService: &document_service.CreateDocumentService{
				Queries: queries,
			},
		},
		ShareableController: &shareable.ShareableController{
			Queries: queries,
		},
		SubscriptionController: &user.SubscriptionController{
			SubscriptionService: &payment_service.SubscriptionService{
				Queries: queries,
			},
			PaymentService: &paymentService,
		},
	}
}

func main() {
	config.LoadEnv()
	// time.Sleep(10 * time.Second)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	sqliteDb, err := sqlite.New("file", config.SqlitePath)
	if err != nil {
		panic(err)
	}
	defer sqliteDb.Close()
	queries := db.New(sqliteDb.GetDB())

	rabbitMqConn, err := amqp.Dial(config.RabbitMqHost)
	if err != nil {
		panic(err)
	}
	defer rabbitMqConn.Close()
	pclient, err := chargily.NewClient(config.ChargiliySecretKey, "test")
	if err != nil {
		panic(err)
	}

	templatesAbsPath, err := filepath.Abs("templates/")
	if err != nil {
		err = os.Mkdir("templates/", 0755)
		if err != nil {
			panic(err)
		}
	}
	storageAbsPath, err := filepath.Abs("storage/")
	if err != nil {
		err = os.Mkdir("storage/", 0755)
		if err != nil {
			panic(err)
		}
	}
	pdfGenerator := pdfservice.Generator{TemplatesDir: templatesAbsPath, Logger: logger, StorageDir: storageAbsPath}
	data := pdfservice.InvoiceTemplateData{
		Title:       "Monthly Payment Report",
		AgencyName:  "PulseTracker Agency",
		Description: "Summary of users who made payments in October",
		PaymentID:   "PAY-2025-10-001",
		Amount:      125000,
		Usernames: []db.User{
			{Fullname: "Alice Dupont", Email: "alice@example.com"},
			{Fullname: "Karim Benali", Email: "karim@example.com"},
		},
	}

	pdfGenerator.Generate("invoice.pdf", "email/invoice.html", data)

	paymentService := payment_service.PaymentService{Client: pclient, Queries: queries}
	controllers := InitServices(logger, sqliteDb.GetDB(), queries, rabbitMqConn, paymentService)
	server := web.NewServer("0.0.0.0:8085", logger, controllers)
	server.Run()
}
