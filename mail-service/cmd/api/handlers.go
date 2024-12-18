package main

import (
	"log"
	"net/http"
	"common/rest"
)

func (app *Config) SendMail(writer http.ResponseWriter, request *http.Request) {
	type mailMessage struct {
		From string `json:"from"`
		To string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := rest.ReadJson(writer, request, &requestPayload)
	if err != nil {
		log.Println(err)
		rest.ErrorJson(writer, err)
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
		rest.ErrorJson(writer, err)
		return
	}

	payload := rest.JsonResponse {
		Error: false, 
		Message: "sent to" + requestPayload.To,
	}

	rest.WriteJson(writer, http.StatusAccepted, payload)
}