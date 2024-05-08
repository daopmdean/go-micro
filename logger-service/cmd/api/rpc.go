package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	coll := client.Database("logs").Collection("logs")
	_, err := coll.InsertOne(context.TODO(), data.LogEntry{
		Name:      "RPC:" + payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting log entry: ", err)
		return err
	}

	*resp = "Log entry created via RPC: " + payload.Name
	return nil
}
