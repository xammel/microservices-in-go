package main

import (
	"log"
	"log-service/data"
	"net/http"
)

type JsonPayload struct {
	Name string `json:"name"`
	Data string `"data"`
}

func (app *Config) WriteLog(writer http.ResponseWriter, request *http.Request) {
	// read json into var
	var requestPayload JsonPayload

	_ = app.readJson(writer, request, &requestPayload)

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJson(writer, err)
		return
	}

	log.Printf("Inserted %+v into Mongo DB \n", event)

	response := jsonResponse {
		Error: false,
		Message: "logged",
	}

	app.writeJson(writer, http.StatusAccepted, response)
}