package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"
	"common/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbit mq
	rabbitConnection, err := connectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConnection.Close()

	// create consumer
	rabbitConn, err := rabbitmq.NewConnection(rabbitConnection)
	if err != nil {
		panic(err)
	}

	// start listening for messages (watch the Q and consume events)
	log.Println("Listening for and consuming RabbitMQ messages...")
	err = event.Listen(&rabbitConn, []string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
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
