package main

import (
	"common/constants"
	commonrpc "common/rpc"
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoUrl = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {

	// Connect to mongodb
	mongoClient, err := connectToMongoDB()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// Create context to be able to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	log.Println("Starting RPC service on port:", rpcPort)
	err = rpc.Register(new(LoggerRPCServer))
	go commonrpc.RPCListen(rpcPort)

	log.Println("Starting gRPC service on port:", constants.GrpcPort)
	go app.gRPCListen()

	log.Println("Starting logger service on port:", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func connectToMongoDB() (*mongo.Client, error) {
	// Create connection options
	clientOptions := options.Client().ApplyURI(mongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//Connect
	client, error := mongo.Connect(context.TODO(), clientOptions)
	if error != nil {
		log.Println("Error connecting: ", error)
		return nil, error
	}

	log.Println("Connected to MongoDB!")

	return client, nil
}
