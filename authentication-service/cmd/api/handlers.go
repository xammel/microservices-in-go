package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"common/jsonhelpers"
)

func (app *Config) Authenticate(writer http.ResponseWriter, request *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := jsonhelpers.ReadJson(writer, request, &requestPayload)

	if err != nil {
		jsonhelpers.ErrorJson(writer, err, http.StatusBadRequest)
		return
	}

	// validate the user against the DB
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		jsonhelpers.ErrorJson(writer, errors.New(err.Error()), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		jsonhelpers.ErrorJson(writer, errors.New("invalid password"), http.StatusUnauthorized)
		return
	}

	// Log authentication to the logger service
	err = app.logLogin("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		jsonhelpers.ErrorJson(writer, err)
		return 
	}

	payload := jsonhelpers.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	jsonhelpers.WriteJson(writer, http.StatusAccepted, payload)

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
