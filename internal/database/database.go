package database

import (
	"database/sql"
	"path"

	"dev/taleroangel/epictetus/internal/entities"
	"dev/taleroangel/epictetus/internal/security"
	"dev/taleroangel/epictetus/internal/storage"
)

var (
	// Setup database filename
	DbFilename = path.Join(storage.StoragePath, "data.sqlite3")
)

// Custom error to represent absence of database
type NilDbRefError struct{}

func (e NilDbRefError) Error() string {
	return "database reference was null"
}

// Create the setup database with initial provided credentials
func CreateDbSchema(initUsr, initPass string, db *sql.DB) error {

	if db == nil {
		return NilDbRefError{}
	}

	// Hash password
	pass, err := security.HashPassword(initPass)
	if err != nil {
		return err
	}

	// Fetch initial credentials
	dbadmin := entities.User{
		User:     initUsr,
		Name:     "Administrator",
		HashPass: pass,
		Sudo:     true,
	}

	// Create the schema
	_, err = db.Exec(createSchemaSql)
	if err != nil {
		return err
	}

	// Create the initial user
	stmt, err := db.Prepare(insertUserSql)
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

func QueryUserByUsername(db *sql.DB, username string) (*entities.User, error) {

	// Check if database is present
	if db == nil {
		return nil, NilDbRefError{}
	}

	// Prepare the statement
	stmt, err := db.Prepare(queryUserByUsernameSql)
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
	var user entities.User
	err = row.Scan(&user.User, &user.HashPass, &user.Name, &user.Sudo)
	if err != nil {
		return nil, err
	}

	// User is queried
	return &user, nil
}
