package main

import (
	"common/constants"
	"common/rpc"
	"context"
	"log"
	"log-service/data"
	"time"
)

type LoggerRPCServer struct{}

func (server *LoggerRPCServer) LogInfo(payload rpc.RPCPayload, response *string) error {
	log.Println("Within LogInfo RPC receiver")
	collection := client.Database(constants.MongoDBName).Collection(constants.MongoDBCollectionName)
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("error writing to MongoDB", err)
		return err
	}

	*response = "Processed payload via RPC:" + payload.Name

	return nil
}
