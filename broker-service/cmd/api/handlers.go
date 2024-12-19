package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"common/rest"
	"common/rabbitmq"
	"common/constants"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := rest.JsonResponse {
		Error:   false,
		Message: "Hit the broker",
	}

	_ = rest.WriteJson(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(writer http.ResponseWriter, request *http.Request) {
	var requestPayload rest.RequestPayload

	error := rest.ReadJson(writer, request, &requestPayload)
	if error != nil {
		rest.ErrorJson(writer, error)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(writer, requestPayload.Auth)
	case "log":
		// Legacy
		//app.logItem(writer, requestPayload.Log)
		app.logEventViaRabbitMQ(writer, requestPayload.Log)
	case "mail":
		app.sendMail(writer, requestPayload.Mail)
	default:
		rest.ErrorJson(writer, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(writer http.ResponseWriter, authPayload rest.AuthPayload) {
	// create some json we'll sent to the auth microservice
	jsonData, _ := json.MarshalIndent(authPayload, "", "\t")

	// call the service
	request, error := http.NewRequest("POST", constants.AuthServiceURL, bytes.NewBuffer(jsonData))
	if error != nil {
		rest.ErrorJson(writer, error)
		return
	}

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		rest.ErrorJson(writer, error)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		rest.ErrorJson(writer, errors.New("invalid credentials"))
		return
	} else if response.StatusCode == http.StatusBadRequest {
		rest.ErrorJson(writer, errors.New("bad request made to postgres"))
		return
	}

	// create var we'll read response.Body into
	var jsonFromService rest.JsonResponse

	// decode json from auth
	error = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if error != nil {
		rest.ErrorJson(writer, error)
		return
	}

	if jsonFromService.Error {
		rest.ErrorJson(writer, error, http.StatusUnauthorized)
		return
	}

	var payload rest.JsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	rest.WriteJson(writer, http.StatusAccepted, payload)
}

/**
*** Legacy ***
**/
func (app *Config) logItem(writer http.ResponseWriter, logEntry rest.LogPayload) {
	// create some json we'll sent to the auth microservice
	jsonData, _ := json.MarshalIndent(logEntry, "", "\t")

	// call the service
	request, error := http.NewRequest("POST", constants.LogServiceURL, bytes.NewBuffer(jsonData))
	if error != nil {
		rest.ErrorJson(writer, error)
		return
	}

	// necessary? not necessary above...
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		rest.ErrorJson(writer, error)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		rest.ErrorJson(writer, error)
		return
	}

	var payload rest.JsonResponse
	payload.Error = false
	payload.Message = "logged"

	rest.WriteJson(writer, http.StatusAccepted, payload)
}

func (app *Config) sendMail(writer http.ResponseWriter, message rest.MailPayload) {
	jsonData, _ := json.MarshalIndent(message, "", "\t")

	// call the mail service
	request, err := http.NewRequest("POST", constants.MailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		rest.ErrorJson(writer, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		rest.ErrorJson(writer, err)
		return
	}

	defer response.Body.Close()

	// make sure we get the right status code back
	if response.StatusCode != http.StatusAccepted {
		rest.ErrorJson(writer, errors.New("error calling mail service"))
		return
	}

	// send back json
	var payload rest.JsonResponse
	payload.Error = false
	payload.Message = "Message was sent to: " + message.To

	rest.WriteJson(writer, http.StatusAccepted, payload)
}

func (app *Config) logEventViaRabbitMQ(writer http.ResponseWriter, logPayload rest.LogPayload) {
	
	// Push event to RabbitMQ
	err := app.pushToQueue(logPayload)
	if err != nil {
		rest.ErrorJson(writer, err)
		return
	}

	// Send success JSON back to caller
	var payload rest.JsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	rest.WriteJson(writer, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(payload rest.LogPayload) error {
	rabbitConn, err := rabbitmq.NewConnection(app.RabbitMQ)
	if err != nil {
		return err
	}

	jsonData, _ := json.MarshalIndent(&payload, "", "\t")
	err = event.Push(&rabbitConn, string(jsonData), "log.INFO")

	if err != nil {
		return err
	}

	return nil
}