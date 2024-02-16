// Provide types for context.Context
package env

// Environment variables in context
type EnvEnum int

const (
	SecretKey EnvEnum = iota
	TokenTTL
	ServerPort
	DatabaseInitialUser
	DatabaseInitialPass
)

// Other objects stored in context
type ContextEnum int

const (
	DatabaseContext ContextEnum = iota
)
