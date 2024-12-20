package main

import (
	"common/rest"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := rest.SetupMux()

	mux.Post("/", app.Broker)
	mux.Post("/handle", app.HandleSubmission)
	mux.Post("/log-grpc", app.LogViaGRPC)

	return mux
}
