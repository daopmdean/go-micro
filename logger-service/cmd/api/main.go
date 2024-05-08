package main

import (
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

// mongodb://admin:password@localhost:27018/logs?authSource=admin&readPreference=primary

const (
	port     = "80"
	rpcPort  = "5001"
	mongoUrl = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	mongoClient, err := connectMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(mongoClient),
	}

	rpc.Register(new(RPCServer))
	go app.rpcListen()

	go app.gRPCListen()

	if err := app.serve(); err != nil {
		log.Println(err)
	}
}

func (app *Config) serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	log.Printf("Logger service started on %s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func connectMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = c.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("---mongo connected")

	return c, nil
}
