package middleware

import (
	"context"
	"dev/taleroangel/epictetus/internal/config"
	"dev/taleroangel/epictetus/internal/errors"
	"dev/taleroangel/epictetus/internal/handler"
	"dev/taleroangel/epictetus/internal/security"
	"net/http"
)

// Authentication type for context.Context variables
type AuthContext int

const (
	// Grab the user from the AuthContext
	UserContext AuthContext = iota
)

func EnsureAuthenticated(next handler.HttpEnvHdlr) handler.HttpEnvHdlr {
	return handler.HttpEnvHdlrFunc(
		func(env config.SrvEnv, r *http.Request) (any, error) {
			// Get the authorization header and check Bearer shcema
			head := r.Header.Get("Authorization")
			if len(head) > len("Bearer ") && head[:len("Bearer ")] != "Bearer " {
				return nil, errors.UnauthorizedHttpError()
			}

			// Get the token without the bearer schema
			tok := head[len("Bearer "):]

			// Get the user data
			usr, err := security.ValidateToken(env.SecretKey, tok)
			if err != nil {
				return nil, err
			}

			// Grab the context and append user to it
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserContext, usr)

			// Call next middleware/handler in chain
			return next.Serve(env, r.WithContext(ctx))
		},
	)
}
