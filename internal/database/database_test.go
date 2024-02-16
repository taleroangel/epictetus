package database

import (
	"context"
	"database/sql"
	"path"
	"testing"

	"dev/taleroangel/epictetus/internal/env"

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

	// Create the context and append DB to context
	ctx := context.Background()
	ctx = context.WithValue(ctx, env.DatabaseInitialUser, "root")
	ctx = context.WithValue(ctx, env.DatabaseInitialPass, "testing")
	ctx = context.WithValue(ctx, env.DatabaseContext, db)

	// Create database
	CreateDatabase(ctx, db)

	// Test queries
	usr, err := QueryUserByUsername(ctx, ctx.Value(env.DatabaseInitialUser).(string))
	if err != nil {
		t.Error(err)
	} else if usr.User != ctx.Value(env.DatabaseInitialUser).(string) {
		t.Error(usr)
	}
}
