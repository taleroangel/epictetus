// Provide types for context.Context
package env

type EnvEnum int

const (
	SecretKey EnvEnum = iota
	TokenValidFor
	DatabaseInitialUser
	DatabaseInitialPass
)
