package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct {

}

func main() {

	app := Config{}

	log.Printf("Starting broker service on port %s \n", webPort)

	// Define http server
	serve := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// Start the server

	err := serve.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}