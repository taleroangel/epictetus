package handler

import (
	"context"
	"dev/taleroangel/epictetus/internal/database"
	"dev/taleroangel/epictetus/internal/security"
	"encoding/json"
	"errors"
	"net/http"
)

// Either username or password is incorrect
type BadCredentialsError struct{}

func (bce BadCredentialsError) Error() string {
	return "bad credentials"
}

// Authenticate user and return token
func AuthSignInHandler(appCtx context.Context) http.Handler {
	// Private types
	type AuthRequest struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}
	type AuthResponse struct {
		Token string `json:"token"`
	}

	// Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get method
		if r.Method != "POST" {
			ErrorHandler(errors.New("only POST method is allowed"), http.StatusBadRequest, w, r)
			return
		}

		// Where to store the authentication request
		var authRequest AuthRequest

		// Get request body
		err := json.NewDecoder(r.Body).Decode(&authRequest)
		if err != nil {
			ErrorHandler(err, http.StatusBadRequest, w, r)
			return
		}

		// Validate in database
		usr, err := database.QueryUserByUsername(appCtx, authRequest.User)
		if err != nil {
			ErrorHandler(BadCredentialsError{}, http.StatusForbidden, w, r)
			return
		}

		// Check password
		if !security.CheckPassword(authRequest.Pass, usr.HashPass) {
			ErrorHandler(BadCredentialsError{}, http.StatusForbidden, w, r)
			return
		}

		// Create the token
		tok, err := security.GenerateToken(appCtx, *usr)
		if err != nil {
			ErrorHandler(err, http.StatusInternalServerError, w, r)
			return
		}

		// Create the response
		resp := AuthResponse{
			Token: tok,
		}

		// Return the response
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			ErrorHandler(err, http.StatusInternalServerError, w, r)
		}
	})
}
