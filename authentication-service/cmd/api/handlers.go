package main

import (
	"bytes"
	"common/rest"
	"encoding/json"
	"errors"
	"fmt"
	"common/constants"
	"net/http"
)

func (app *Config) Authenticate(writer http.ResponseWriter, request *http.Request) {

	var requestPayload rest.AuthPayload

	err := rest.ReadJson(writer, request, &requestPayload)

	if err != nil {
		rest.ErrorJson(writer, err, http.StatusBadRequest)
		return
	}

	// validate the user against the DB
	user, err := app.Repo.GetByEmail(requestPayload.Email)
	if err != nil {
		rest.ErrorJson(writer, errors.New(err.Error()), http.StatusBadRequest)
		return
	}

	valid, err := app.Repo.PasswordMatches(requestPayload.Password, *user)
	if err != nil || !valid {
		rest.ErrorJson(writer, errors.New("invalid password"), http.StatusUnauthorized)
		return
	}

	// Log authentication to the logger service
	err = app.logLogin("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		rest.ErrorJson(writer, err)
		return 
	}

	payload := rest.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	rest.WriteJson(writer, http.StatusAccepted, payload)

}

func (app *Config) logLogin(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest("POST", constants.LogServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	_, err = app.Client.Do(request)

	if err != nil {
		return err
	}

	return nil
}
