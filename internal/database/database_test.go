package database

import (
	"database/sql"
	"path"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestUserQueries(t *testing.T) {

	// Database file, delete after test
	var dbFile = path.Join(t.TempDir(), "dbtest.sqlite3")
	t.Logf("Database stored in: %s", dbFile)

	// Open the database
	db, err := sql.Open("sqlite3", dbFile+"?cache=shared&mode=memory")
	if err != nil {
		t.Fatal("Failed to open the database")
	}

	// Close the database after test
	t.Cleanup(func() {
		db.Close()
	})

	// Create database
	CreateDbSchema("root", "root", db)

	// Test queries
	usr, err := QueryUserByUsername(db, "root")
	if err != nil {
		t.Error(err)
	} else if usr.User != "root" {
		t.Error(usr)
	}
}
