package database

import (
	"context"
	"database/sql"
	"errors"
	"path"

	"dev/taleroangel/epictetus/internal/env"
	"dev/taleroangel/epictetus/internal/security"
	"dev/taleroangel/epictetus/internal/storage"
	"dev/taleroangel/epictetus/internal/types"

	_ "embed"
)

// Custom error to represent absence of results
type NullDatabaseReference struct{}

func (e NullDatabaseReference) Error() string {
	return "database reference was null"
}

var (
	// Setup database filename
	DbFilename = path.Join(storage.StoragePath, "data.sqlite3")
)

// Scripts embeded queries
var (
	//go:embed scripts/schema.sql
	databaseSchemaSql string
	//go:embed scripts/create_user_stmt.sql
	createUserStmtSql string
	//go:embed scripts/user_by_username_query.sql
	userByUsernameQrySql string
)

// Create the setup database with initial provided credentials
func CreateDatabase(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return NullDatabaseReference{}
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

// * Queries * //

func QueryUserByUsername(ctx context.Context, username string) (*types.User, error) {
	// Get database from context
	db := ctx.Value(env.DatabaseContext).(*sql.DB)
	if db == nil {
		return nil, NullDatabaseReference{}
	}

	// Prepare the statement
	stmt, err := db.Prepare(userByUsernameQrySql)
	if err != nil {
		return nil, err
	}

	// Execute the statement
	row := stmt.QueryRow(username)
	if row.Err() != nil {
		// Returned without results
		return nil, row.Err()
	}

	// Get user from database
	var user types.User
	err = row.Scan(&user.User, &user.HashPass, &user.Name, &user.Sudo)
	if err != nil {
		return nil, err
	}

	// User is queried
	return &user, nil
}
