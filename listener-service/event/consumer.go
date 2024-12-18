package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"common/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	connection *amqp.Connection
	queueName  string
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(connection *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		connection: connection,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return event.DeclareExchange(channel)
}

func (consumer *Consumer) Listen(topics []string) error {
	channel, err := consumer.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	queue, err := event.DeclareRandomQueue(channel)
	if err != nil {
		return err
	}

	for _, s := range topics {
		log.Printf("Binding channel to queue with params: name: %s, key: %s, exchange: %s", queue.Name, s, "logs_topic")
		err = channel.QueueBind(
			queue.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", queue.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	log.Printf("Received payload: %+v \n", payload)
	switch payload.Name {
	case "log", "event":
		// log it all
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	case "auth":
		// authenticate

	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	// create some json we'll sent to the auth microservice
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// necessary? not necessary above...
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
