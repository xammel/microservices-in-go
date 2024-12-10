package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(writer http.ResponseWriter, request *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(writer, request, &requestPayload)

	if err != nil {
		app.errorJson(writer, err, http.StatusBadRequest)
		return
	}

	// validate the user against the DB
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJson(writer, errors.New(err.Error()), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJson(writer, errors.New("invalid password"), http.StatusUnauthorized)
		return
	}

	// Log authentication to the logger service
	err = app.logLogin("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJson(writer, err)
		return 
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	app.writeJson(writer, http.StatusAccepted, payload)

}

func (app *Config) logLogin(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)

	if err != nil {
		return err
	}

	return nil
}
