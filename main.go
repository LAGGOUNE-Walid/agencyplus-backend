package main

import (
	"database/sql"
	"log/slog"
	"logispro/internal/config"
	"logispro/internal/db"
	"logispro/internal/services/building_service"
	"logispro/internal/services/contact_service"
	"logispro/internal/services/user_services"
	"logispro/internal/sqlite"
	"logispro/internal/web"
	"logispro/internal/web/controllers"
	"logispro/internal/web/controllers/building"
	"logispro/internal/web/controllers/contact"
	"logispro/internal/web/controllers/user"
	"os"
)

var JwtSecret = []byte("your-secret-key") // ⚠️ move to env/config

func InitServices(logger *slog.Logger, db *sql.DB, queries *db.Queries) controllers.Controller {
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
		},
		ContactController: &contact.ContactController{
			CreateContactService: &contact_service.CreateContactService{
				Queries: queries,
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
				Queries: queries,
				DB:      db,
			},
			GetBuildingService: &building_service.GetBuildingService{
				Queries: queries,
			},
			UpdateBuildingService: &building_service.UpdateBuildingService{
				Queries: queries,
			},
		},
	}
}

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	sqliteDb, err := sqlite.New("file", "db/database.sqlite")
	queries := db.New(sqliteDb.GetDB())
	if err != nil {
		panic(err)
	}
	config.LoadEnv()
	controllers := InitServices(logger, sqliteDb.GetDB(), queries)
	server := web.NewServer("0.0.0.0:8085", logger, controllers)
	server.Run()
}
