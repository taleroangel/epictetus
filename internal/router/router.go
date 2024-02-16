package router

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"dev/taleroangel/epictetus/internal/config"
	"dev/taleroangel/epictetus/internal/errors"
	"dev/taleroangel/epictetus/internal/handler"
)

type HttpEnvRouter struct {
	// Server environment
	config.SrvEnv
	Routes map[string]handler.HttpEnvHdlr
}

func (router HttpEnvRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Grab the URI
	uri := r.URL.RequestURI()
	// Create the response encoder
	enc := json.NewEncoder(w)

	// Compare every route regex with the URI
	for rr, handler := range router.Routes {
		// Check if regex matches
		if mat, _ := regexp.MatchString(rr, uri); mat {
			// Call the handler
			resp, err := handler.Serve(router.SrvEnv, r)

			// Write the error if found
			if err != nil {
				switch herr := err.(type) {
				// HTTP error is directly sent
				case errors.HttpStatusError:
					// 404 is managed by its own
					if herr.Code == http.StatusNotFound {
						http.NotFound(w, r)
						return
					}
					// Other http errors are managed
					enc.Encode(err)

				// Other errors are first encoded into HttpStatusError
				default:
					log.Printf("Error while processing route `%s` caused and INTERNAL_SERVER_ERROR (%s)", rr, err.Error())
					enc.Encode(&errors.HttpStatusError{
						Code: http.StatusInternalServerError,
						Err:  err.Error(),
					})
				}
			} else if resp != nil {
				// Return the response
				enc.Encode(resp)
			} else {
				// If no response is returned, then return Ok
				w.WriteHeader(http.StatusOK)
			}

			// Exit the loop, route was already found
			return
		}
	}

	// Handle 404
	http.NotFound(w, r)
}

type SubRoute struct {
	Routes map[string]handler.HttpEnvHdlr
}

func (router SubRoute) Serve(env config.SrvEnv, r *http.Request) (any, error) {
	// Grab the URI
	uri := r.URL.RequestURI()

	// Compare every route regex with the URI
	for rr, handler := range router.Routes {
		// Check if regex matches
		if mat, _ := regexp.MatchString(rr, uri); mat {
			// Call the handler
			return handler.Serve(env, r)
		}
	}

	// Subroute not found
	return nil, errors.NotFoundHttpError()
}
