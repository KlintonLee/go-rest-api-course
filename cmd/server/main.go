package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/KlintonLee/go-rest-api-course/internal/database"
	transportHTTP "github.com/KlintonLee/go-rest-api-course/internal/transport/http"
	"github.com/joho/godotenv"
)

// App - the struct which contains things like
// pointers to database connections
type App struct{}

// Run - handles the startup of our application
func (app *App) Run() error {
	var err error
	fmt.Println("Setting up our application")

	err = godotenv.Load(os.ExpandEnv("$GOPATH/src/go-rest-api-course/.env"))
	if err != nil {
		return err
	}

	_, err = database.NewDatabase()
	if err != nil {
		return err
	}

	handler := transportHTTP.NewHandler()
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		fmt.Println("Failed to set up server")
		return err
	}

	return nil
}

// Our main entrypoint for the application
func main() {
	app := App{}
	if err := app.Run(); err != nil {
		fmt.Println("Error starting up!")
		fmt.Println(err)
	}
}
