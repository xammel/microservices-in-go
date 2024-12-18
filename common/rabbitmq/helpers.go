package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type RabbitConnection struct {
	Connection *amqp.Connection
}

func NewConnection(connection *amqp.Connection) (RabbitConnection, error) {
	rabbitConn := RabbitConnection{
		Connection: connection,
	}

	err := rabbitConn.setup()
	if err != nil {
		return RabbitConnection{}, err
	}

	return rabbitConn, nil
}

func (rabbitConn *RabbitConnection) setup() error {
	channel, err := rabbitConn.Connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return DeclareExchange(channel)
}