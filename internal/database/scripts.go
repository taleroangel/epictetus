package database

import _ "embed"

// Scripts embeded queries
var (
	//go:embed scripts/schema/create_schema.sql
	createSchemaSql string
	//go:embed scripts/insert/insert_user.sql
	insertUserSql string
	//go:embed scripts/queries/query_user_by_username.sql
	queryUserByUsernameSql string
)
