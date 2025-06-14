package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const webPort = "8080"

type Config struct {
	RabbitMQ *amqp.Connection
}

func main() {

	// try to connect to rabbit mq
	rabbitConnection, err := connectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConnection.Close()

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	app := Config{
		RabbitMQ: rabbitConnection,
	}

	log.Printf("Starting broker service on port %s \n", webPort)

	var handler http.Handler = app.routes()
	handler = otelhttp.NewHandler(handler, "/")

	// Define http server
	serve := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: handler,
	}

	// Start the server
	err = serve.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}

func connectToRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbitMQ is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		fmt.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
