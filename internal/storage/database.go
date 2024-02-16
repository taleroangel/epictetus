package storage

import (
	"context"
	"database/sql"
	"errors"

	"dev/taleroangel/epictetus/internal/env"
	"dev/taleroangel/epictetus/internal/security"
	"dev/taleroangel/epictetus/internal/types"

	_ "embed"
)

const (
	// Setup database filename
	DbFilename = (StoragePath + "data.sqlite3")
)

// Scripts embeded queries
var (
	//go:embed scripts/schema.sql
	databaseSchemaSql string
	//go:embed scripts/create_user_stmt.sql
	createUserStmtSql string
)

// Create the setup database with initial provided credentials
func CreateDatabase(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return errors.New("database reference is null")
	}

	// Hash password
	pass, err := security.HashPassword(ctx.Value(env.DatabaseInitialPass).(string))
	if err != nil {
		return err
	}

	// Fetch initial credentials
	dbadmin := types.User{
		User:     ctx.Value(env.DatabaseInitialUser).(string),
		Name:     "Administrator",
		HashPass: pass,
		Sudo:     true,
	}

	// Check if credentials are valid
	if (dbadmin.User == "") || (dbadmin.HashPass == "") {
		return errors.New("initial administrator credentials are missing")
	}

	// Create the table
	_, err = db.Exec(databaseSchemaSql)
	if err != nil {
		return err
	}

	// Create the user
	stmt, err := db.Prepare(createUserStmtSql)
	if err != nil {
		return err
	}

	// Execute the statement
	defer stmt.Close()
	_, err = stmt.Exec(
		dbadmin.User,
		dbadmin.HashPass,
		dbadmin.Name,
		dbadmin.Sudo,
	)

	// Return latest error
	return err
}
