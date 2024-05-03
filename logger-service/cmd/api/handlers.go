package main

import (
	"log-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var payload JSONPayload
	if err := app.readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, err)
		return
	}

	event := data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, jsonRes{
		Error:   false,
		Message: "logged",
		Data:    event,
	})
}

func (app *Config) GetLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := app.Models.LogEntry.All()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, jsonRes{
		Error:   false,
		Message: "get logs successful",
		Data:    logs,
	})
}
