package main

import (
	"common/rest"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := rest.SetupMux()

	mux.Post("/authenticate", app.Authenticate)

	return mux
}
