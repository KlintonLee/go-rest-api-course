package main

import (
	"net/http"
	"os"

	"github.com/KlintonLee/go-rest-api-course/internal/database"
	"github.com/KlintonLee/go-rest-api-course/internal/database/repositories"
	transportHTTP "github.com/KlintonLee/go-rest-api-course/internal/transport/http"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// App - contains application information
type App struct {
	Name    string
	Version string
}

// Run - handles the startup of our application
func (app *App) Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(
		log.Fields{
			"AppName":    app.Name,
			"AppVersion": app.Version,
		}).Info("Setting up our application")

	var err error

	err = godotenv.Load(os.ExpandEnv("$GOPATH/src/go-rest-api-course/.env"))
	if err != nil {
		log.Error("Failed to load enviroments variables")
		return err
	}

	db, err := database.NewDatabase()
	if err != nil {
		log.Error("Failed to setup database")
		return err
	}

	err = database.MigrateDB(db)
	if err != nil {
		log.Error("Failed to migrate database")
		return err
	}

	commentsRepositoryDb := repositories.CommentsRepositoryDB{Db: db}

	handler := transportHTTP.NewHandler(&commentsRepositoryDb)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		log.Error("Failed to set up server")
		return err
	}

	return nil
}

// Our main entrypoint for the application
func main() {
	app := App{
		Name:    "Comment API",
		Version: "1.0",
	}
	if err := app.Run(); err != nil {
		log.Error("Error starting up our REST API")
		log.Fatal(err)
	}
}
