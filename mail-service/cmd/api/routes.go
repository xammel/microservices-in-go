package main

import (
	"common/rest"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := rest.SetupMux()

	mux.Post("/send", app.SendMail)

	return mux
}
