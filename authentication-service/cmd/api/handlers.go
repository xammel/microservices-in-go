package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(writer http.ResponseWriter, request *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	error := app.readJson(writer, request, &requestPayload)

	if error != nil {
		app.errorJson(writer, error, http.StatusBadRequest)
		return
	}

	// validate the user against the DB
	user, error := app.Models.User.GetByEmail(requestPayload.Email)
	if error != nil {
		app.errorJson(writer, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, error := user.PasswordMatches(requestPayload.Password)
	if error != nil || !valid {
		app.errorJson(writer, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	app.writeJson(writer, http.StatusAccepted, payload)

}
