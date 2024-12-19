package main

import (
	"log"
	"net/http"
	"common/rest"
)

func (app *Config) SendMail(writer http.ResponseWriter, request *http.Request) {

	var requestPayload rest.MailPayload

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