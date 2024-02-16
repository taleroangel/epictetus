package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"dev/taleroangel/epictetus/api/handler"
	"dev/taleroangel/epictetus/internal/database"
	"dev/taleroangel/epictetus/internal/env"
	"dev/taleroangel/epictetus/internal/storage"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// * --- Database Creation --- * //

	// Read environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Missing .env file, please check documentation")
	}

	// Create the context
	ctx := context.Background()

	// Append variables to context
	ctx = context.WithValue(ctx, env.SecretKey, os.Getenv("SECRET_KEY"))
	ctx = context.WithValue(ctx, env.TokenTTL, os.Getenv("TOKEN_TTL_HRS"))
	ctx = context.WithValue(ctx, env.ServerPort, os.Getenv("SERVER_PORT"))

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
		ctx = context.WithValue(ctx, env.DatabaseInitialUser, os.Getenv("SETUP_INITIAL_USER"))
		ctx = context.WithValue(ctx, env.DatabaseInitialPass, os.Getenv("SETUP_INITIAL_PASS"))

		log.Print("Database does not exist, creating initial schema")
		err := database.CreateDatabase(ctx, db)

		// Delete database and exit program
		if err != nil {
			db.Close()
			os.Remove(database.DbFilename)
			log.Fatalf("Failed to create database: (%s)", err.Error())
		}
	}

	// Store database in context
	ctx = context.WithValue(ctx, env.DatabaseContext, db)

	// * --- Create the HTTP server --- * //

	// Print ports
	lp := fmt.Sprintf("localhost:%s", ctx.Value(env.ServerPort))
	log.Printf("Server binded to `%s`", lp)

	// Bind handlers
	http.Handle("/auth/signin", handler.AuthSignInHandler(ctx))

	// Start the HTTP server
	http.ListenAndServe(lp, nil)
}
