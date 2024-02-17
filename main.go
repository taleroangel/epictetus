package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"dev/taleroangel/epictetus/api/middleware"
	"dev/taleroangel/epictetus/api/routes"
	"dev/taleroangel/epictetus/internal/config"
	"dev/taleroangel/epictetus/internal/database"
	"dev/taleroangel/epictetus/internal/handler"
	"dev/taleroangel/epictetus/internal/router"
	"dev/taleroangel/epictetus/internal/storage"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Read environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Missing .env file, please check documentation")
	}

	// Create configuration from environment
	cfg, err := config.NewCfgFromEnv()
	if err != nil {
		log.Fatalf("Environment is missing variables (%s)", err.Error())
	}

	// Create or ensure that storage path is created
	err = storage.CreateStoragePath()
	if err != nil {
		log.Fatalf("Cannot access application storage `%s` (%s)", storage.StoragePath, err.Error())
	}

	// Open the database
	db, err := sql.Open("sqlite3", database.DbFilename)
	if err != nil {
		log.Fatal("Failed to open the database")
	}

	// Close the database at the end of execution
	defer db.Close()

	// Check if database doest not exist and create it
	if _, err := os.Stat(database.DbFilename); err != nil {

		// Add variables to context
		initUsr := os.Getenv("SETUP_INITIAL_USER")
		initPass := os.Getenv("SETUP_INITIAL_PASS")

		// Check if env variables
		if len(initUsr) == 0 || len(initPass) == 0 {
			db.Close()
			os.Remove(database.DbFilename)
			log.Fatalf("Initial credential env variables `SETUP_INITIAL_USER` and `SETUP_INITIAL_PASS` were not found")
		}

		log.Print("Database does not exist, creating initial schema")
		err := database.CreateDbSchema(initUsr, initPass, db)

		// Delete database and exit program
		if err != nil {
			db.Close()
			os.Remove(database.DbFilename)
			log.Fatalf("Failed to create database: (%s)", err.Error())
		}
	}

	// Create the server env
	srvenv := config.SrvEnv{
		Database: db,
		Cfg:      cfg,
	}

	// Print ports
	lp := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	log.Printf("Server binded to `%s`", lp)

	// Configure the router
	router := &router.HttpEnvRouter{
		SrvEnv: srvenv,
		Routes: map[string]handler.HttpEnvHdlr{
			"/auth/signin": handler.HttpEnvHdlrFunc(routes.AuthLoginHandler),
			"/*": middleware.EnsureAuthenticated(router.SubRoute{Routes: map[string]handler.HttpEnvHdlr{
				"/auth/manage": router.RouteMethods{
					Methods: map[string]handler.HttpEnvHdlr{
						"GET": handler.HttpEnvHdlrFunc(routes.AuthGetUser),
					},
				},
			}}),
		},
	}

	// Start the HTTP server
	http.ListenAndServe(lp, router)
}
