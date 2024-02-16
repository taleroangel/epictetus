package handler

import (
	"encoding/json"
	"net/http"

	"dev/taleroangel/epictetus/api/response"
)

func ErrorHandler(err error, code int, w http.ResponseWriter, r *http.Request) {
	// Parse error into json response
	b, err := json.Marshal(response.NewReqResFromError(err))

	// If parsing failed return server error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Return the error
	http.Error(w, string(b), code)
}
