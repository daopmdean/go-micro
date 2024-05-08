package main

import "net/http"

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, jsonRes{
		Error:   false,
		Message: "logged via RabbitMQ",
	})
}
