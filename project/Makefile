FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerApp
MAILER_BINARY=mailerApp
LISTENER_BINARY=listenerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
build_and_docker_compose: compile
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

stop_swarm: 
	@echo "Stopping swarm"
	docker stack rm myapp
	@echo "Stopped swarm!"

build_swarm_images: 
	@echo "Building swarm docker images..."
	cd ../broker-service && docker build -f broker-service.dockerfile -t maxmel/broker-service:1.0.0 .
	cd ../authentication-service && docker build -f authentication-service.dockerfile -t maxmel/authentication-service:1.0.0 . 
	cd ../logger-service && docker build -f logger-service.dockerfile -t maxmel/logger-service:1.0.0 . 
	cd ../mail-service && docker build -f mail-service.dockerfile -t maxmel/mail-service:1.0.0 . 
	cd ../listener-service && docker build -f listener-service.dockerfile -t maxmel/listener-service:1.0.0 . 
	cd ../front-end && docker build -f front-end.dockerfile -t maxmel/front-end:1.0.0 .
	cd ../project && docker build -f caddy.dockerfile -t maxmel/micro-caddy:1.0.0 .

start_swarm: 
	@echo "Starting swarm..."
	docker stack deploy -c swarm.yml myapp
	@echo "Started swarm!"

build_and_start_swarm: stop_swarm compile build_swarm_images start_swarm

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building authentication binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

build_mailer:
	@echo "Building mailer binary..."
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAILER_BINARY} ./cmd/api
	@echo "Done!"

build_listener:
	@echo "Building listener binary..."
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Done!"

## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env GOOS=linux CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## build_front_local: builds the front end library for mac silicon
build_front_local:
	@echo "Building front end binary..."
	cd ../front-end && env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front_local
	@echo "Starting front end"
	cd ../front-end && export BROKER_URL="http://localhost:8080" && ./${FRONT_END_BINARY}

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

compile: build_broker build_auth build_logger build_mailer build_listener build_front
test:
	@echo "Running tests..."
	cd ../broker-service && go test -v ./... 
	cd ../authentication-service && go test -v ./... 
	cd ../logger-service && go test -v ./... 
	cd ../mail-service && go test -v ./... 
	cd ../listener-service && go test -v ./... 
	@echo "Tests completed!"
start_all: up_build start
kill_all: down stop
