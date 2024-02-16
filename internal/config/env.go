package config

import (
	"database/sql"
)

// Server environment with dabatase pool and configuration
type SrvEnv struct {
	// Database pool
	Database *sql.DB
	// Server configurations
	*Cfg
}
