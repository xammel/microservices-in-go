module listener

go 1.23.4

// require github.com/rabbitmq/amqp091-go v1.10.0 // indirect
// require common v0.0.0
require (
	common v0.0.0
	github.com/rabbitmq/amqp091-go v1.10.0
)

require (
	github.com/go-chi/chi/v5 v5.2.0 // indirect
	github.com/go-chi/cors v1.2.1 // indirect
)

replace common => ../common
