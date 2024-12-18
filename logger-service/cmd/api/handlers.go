package main

import (
	"common/rest"
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

	_ = rest.ReadJson(writer, request, &requestPayload)

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		rest.ErrorJson(writer, err)
		return
	}

	log.Printf("Inserted %+v into Mongo DB \n", event)

	response := rest.JsonResponse {
		Error: false,
		Message: "logged",
	}

	rest.WriteJson(writer, http.StatusAccepted, response)
}