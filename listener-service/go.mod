module listener

go 1.23.4

// require github.com/rabbitmq/amqp091-go v1.10.0 // indirect
// require common v0.0.0
require (
	common v0.0.0
	github.com/rabbitmq/amqp091-go v1.10.0
)

replace common => ../common
