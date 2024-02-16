package routes

import (
	"encoding/json"
	stderr "errors"
	"net/http"

	"dev/taleroangel/epictetus/api/middleware"
	"dev/taleroangel/epictetus/internal/config"
	"dev/taleroangel/epictetus/internal/database"
	"dev/taleroangel/epictetus/internal/entities"
	"dev/taleroangel/epictetus/internal/errors"
	"dev/taleroangel/epictetus/internal/security"
)

func AuthLoginHandler(svrenv config.SrvEnv, r *http.Request) (any, error) {

	// Check that request is POST
	if r.Method != "POST" {
		return nil, errors.BadRequestHttpError()
	}

	// Get the body of the response
	var requestBody struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}

	// Decode the JSON body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		return nil, errors.BadRequestHttpError()
	}

	// Search for user in database
	usr, err := database.QueryUserByUsername(svrenv.Database, requestBody.User)
	if err != nil {
		return nil, errors.UnauthorizedHttpError()
	}

	// Check if password matches
	if !security.CheckPassword(requestBody.Pass, usr.HashPass) {
		return nil, errors.UnauthorizedHttpError()
	}

	// Generate the token
	tkn, err := security.GenerateToken(svrenv.SecretKey, *usr)
	if err != nil {
		return nil, err
	}

	// No error was found, return the token
	return struct {
		Token string `json:"token"`
	}{tkn}, nil
}

func AuthGetUser(svrenv config.SrvEnv, r *http.Request) (any, error) {

	// Check that request is GET
	if r.Method != "GET" {
		return nil, errors.BadRequestHttpError()
	}

	// Get the user from context
	user, present := r.Context().Value(middleware.UserContext).(*entities.User)
	if !present {
		return nil, errors.ForbiddenHttpError()
	}

	// Check if user is not null
	if user == nil {
		return nil, stderr.New("user cannot be null from an authentication context")
	}

	return user, nil
}
