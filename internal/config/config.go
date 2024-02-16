package config

import (
	"fmt"
	"os"
	"strconv"
)

// Error when a variable is not found in env
type UnsetVariableError struct {
	variable string
}

func (e UnsetVariableError) Error() string {
	return fmt.Sprintf("environment variable `%s` was not set", e.variable)
}

// Server configurations
type Cfg struct {
	// Serve port
	Port int
	// Security key
	SecretKey []byte
}

// Crete the configuration from environment variables
func NewCfgFromEnv() (*Cfg, error) {
	// Obtain the port as string
	pstr := os.Getenv("SERVER_PORT")
	if len(pstr) == 0 {
		return nil, UnsetVariableError{"SERVER_PORT"}
	}

	// Parse the port
	port, err := strconv.Atoi(pstr)
	if err != nil {
		return nil, err
	}

	// Retrieve the secret key
	skey := os.Getenv("SECRET_KEY")
	if len(skey) == 0 {
		return nil, UnsetVariableError{"SECRET_KEY"}
	}

	// Return the configuration
	return &Cfg{
		Port:      port,
		SecretKey: []byte(skey),
	}, nil
}
