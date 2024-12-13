package main

import (
	"log"
	"net/http"
)

func (app *Config) SendMail(writer http.ResponseWriter, request *http.Request) {
	type mailMessage struct {
		From string `json:"from"`
		To string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := app.readJson(writer, request, &requestPayload)
	if err != nil {
		log.Println(err)
		app.errorJson(writer, err)
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
		app.errorJson(writer, err)
		return
	}

	payload := jsonResponse {
		Error: false, 
		Message: "sent to" + requestPayload.To,
	}

	app.writeJson(writer, http.StatusAccepted, payload)
}