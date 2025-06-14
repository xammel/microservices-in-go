module broker

go 1.23.4

require (
	common v0.0.0
	github.com/rabbitmq/amqp091-go v1.10.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0
	go.opentelemetry.io/otel v1.36.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.12.2
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.36.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.36.0
	go.opentelemetry.io/otel/log v0.12.2
	go.opentelemetry.io/otel/sdk v1.36.0
	go.opentelemetry.io/otel/sdk/log v0.12.2
	go.opentelemetry.io/otel/sdk/metric v1.36.0
	golang.org/x/net v0.40.0
	google.golang.org/grpc v1.72.1
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-chi/chi/v5 v5.2.0 // indirect
	github.com/go-chi/cors v1.2.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250519155744-55703ea1f237 // indirect
)

replace common => ../common
