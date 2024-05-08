package main

import (
	"broker/logs"
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload

	if err := app.readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, err)
		return
	}

	conn, err := grpc.Dial(
		"logger-service:50001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: payload.Log.Name,
			Data: payload.Log.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, jsonRes{
		Error:   false,
		Message: "logged",
		Data:    res.Result,
	})
}
