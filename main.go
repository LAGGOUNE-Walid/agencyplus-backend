package main

import (
	"log/slog"
	"logispro/internal/config"
	"logispro/internal/db"
	"logispro/internal/services/user_services"
	"logispro/internal/sqlite"
	"logispro/internal/web"
	"logispro/internal/web/controllers"
	"logispro/internal/web/controllers/user"
	"os"
)

var JwtSecret = []byte("your-secret-key") // ⚠️ move to env/config

func InitServices(logger *slog.Logger, queries *db.Queries) controllers.Controller {
	return controllers.Controller{
		UserController: &user.UserController{
			CreateUserService: &user_services.CreateUserService{
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
	controllers := InitServices(logger, queries)
	server := web.NewServer("0.0.0.0:8085", logger, controllers)
	server.Run()
}
