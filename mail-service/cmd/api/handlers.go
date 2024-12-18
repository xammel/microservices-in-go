package main

import (
	"log"
	"net/http"
	"common/jsonhelpers"
)

func (app *Config) SendMail(writer http.ResponseWriter, request *http.Request) {
	type mailMessage struct {
		From string `json:"from"`
		To string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := jsonhelpers.ReadJson(writer, request, &requestPayload)
	if err != nil {
		log.Println(err)
		jsonhelpers.ErrorJson(writer, err)
		return 
	}

	msg := Message {
		From: requestPayload.From,
		To: requestPayload.To, 
		Subject: requestPayload.Subject,
		Data: requestPayload.Message,
	}

	err = app.Mailer.sendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		jsonhelpers.ErrorJson(writer, err)
		return
	}

	payload := jsonhelpers.JsonResponse {
		Error: false, 
		Message: "sent to" + requestPayload.To,
	}

	jsonhelpers.WriteJson(writer, http.StatusAccepted, payload)
}