package main

import (
	"common/jsonhelpers"
	"log"
	"log-service/data"
	"net/http"
)

type JsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(writer http.ResponseWriter, request *http.Request) {
	// read json into var
	var requestPayload JsonPayload

	_ = jsonhelpers.ReadJson(writer, request, &requestPayload)

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		jsonhelpers.ErrorJson(writer, err)
		return
	}

	log.Printf("Inserted %+v into Mongo DB \n", event)

	response := jsonhelpers.JsonResponse {
		Error: false,
		Message: "logged",
	}

	jsonhelpers.WriteJson(writer, http.StatusAccepted, response)
}