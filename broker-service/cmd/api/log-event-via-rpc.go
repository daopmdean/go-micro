package main

import (
	"net/http"
	"net/rpc"
)

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logEventViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := RPCPayload(l)

	var result string
	err = client.Call("RPCServer.LogInfo", payload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, jsonRes{
		Error:   false,
		Message: "logged via RPC",
		Data:    result,
	})
}
