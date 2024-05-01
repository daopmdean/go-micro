package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonRes{
		Error:   false,
		Message: "Hit the broker",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
