package main

import (
	"common/rest"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := rest.SetupMux()

	// 	Configure the "http.route" for the HTTP instrumentation.
	//	handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
	//	mux.Handle(pattern, handler)
	mux.Post("/", app.Broker)
	mux.Post("/handle", app.HandleSubmission)
	mux.Post("/log-grpc", app.LogViaGRPC)

	return mux
}
