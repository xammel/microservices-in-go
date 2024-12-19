package main

import (
	"net/http"
	"common/rest"
)

func (app *Config) routes() http.Handler {
	mux := rest.SetupMux()

	mux.Post("/", app.Broker)
	mux.Post("/handle", app.HandleSubmission)

	return mux
}
