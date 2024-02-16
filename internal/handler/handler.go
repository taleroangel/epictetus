package handler

import (
	"net/http"

	"dev/taleroangel/epictetus/internal/config"
)

// Handler interface
type HttpEnvHdlr interface {
	Serve(env config.SrvEnv, r *http.Request) (any, error)
}

// HTTP handler with config.SrvEnv support and HTTP error return support
type HttpEnvHdlrFunc func(env config.SrvEnv, r *http.Request) (any, error)

func (f HttpEnvHdlrFunc) Serve(env config.SrvEnv, r *http.Request) (any, error) {
	return f(env, r)
}

// HTTP Middleware function
type HttpEnvMiddleware func(next HttpEnvHdlr) HttpEnvHdlr
