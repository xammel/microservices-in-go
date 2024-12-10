package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(writer http.ResponseWriter, request *http.Request) {
	var requestPayload RequestPayload

	error := app.readJson(writer, request, &requestPayload)
	if error != nil {
		app.errorJson(writer, error)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(writer, requestPayload.Auth)
	case "log":
		app.logItem(writer, requestPayload.Log)
	default:
		app.errorJson(writer, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(writer http.ResponseWriter, authPayload AuthPayload) {
	// create some json we'll sent to the auth microservice
	jsonData, _ := json.MarshalIndent(authPayload, "", "\t")

	// call the service
	request, error := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if error != nil {
		app.errorJson(writer, error)
		return
	}

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		app.errorJson(writer, error)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(writer, errors.New("invalid credentials"))
		return
	} else if response.StatusCode == http.StatusBadRequest {
		app.errorJson(writer, errors.New("bad request made to postgres"))
		return
	}

	// create var we'll read response.Body into
	var jsonFromService jsonResponse

	// decode json from auth
	error = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if error != nil {
		app.errorJson(writer, error)
		return
	}

	if jsonFromService.Error {
		app.errorJson(writer, error, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJson(writer, http.StatusAccepted, payload)
}

func (app *Config) logItem(writer http.ResponseWriter, logEntry LogPayload) {
	// create some json we'll sent to the auth microservice
	jsonData, _ := json.MarshalIndent(logEntry, "", "\t")

	// call the service
	request, error := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if error != nil {
		app.errorJson(writer, error)
		return
	}

	// necessary? not necessary above...
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		app.errorJson(writer, error)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJson(writer, error)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJson(writer, http.StatusAccepted, payload)
}
