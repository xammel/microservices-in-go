package event

import (
	"log"
	"common/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Push(rabbitConn *rabbitmq.RabbitConnection, event string, severity string) error {
	channel, err := rabbitConn.Connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Printf("Pushing following event to RabbitMQ: %+v", event)

	err = channel.Publish(
		"logs_topic", 
		severity, 
		false, 
		false, 
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}