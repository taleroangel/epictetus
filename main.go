package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"dev/taleroangel/epictetus/internal/env"
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

	// Create or ensure that storage path is created
	err = storage.CreateStoragePath()
	if err != nil {
		log.Fatalf("Cannot access application storage `%s` (%s)", storage.StoragePath, err.Error())
	}

	// Create the context
	ctx := context.Background()

	// Open the database
	db, err := sql.Open("sqlite3", storage.DbFilename)
	if err != nil {
		log.Fatal("Failed to open the database")
	}

	// Close the database at the end of execution
	defer db.Close()

	// Check if database doest not exist
	if _, err := os.Stat(storage.DbFilename); err != nil {
		// Add variables to context
		ctx = context.WithValue(ctx, env.DatabaseInitialUser, os.Getenv("SETUP_INITIAL_USER"))
		ctx = context.WithValue(ctx, env.DatabaseInitialPass, os.Getenv("SETUP_INITIAL_PASS"))

		log.Print("Database does not exist, creating initial schema")
		err := storage.CreateDatabase(ctx, db)

		// Delete database and exit program
		if err != nil {
			db.Close()
			os.Remove(storage.DbFilename)
			log.Fatalf("Failed to create database: (%s)", err.Error())
		}
	}
}
