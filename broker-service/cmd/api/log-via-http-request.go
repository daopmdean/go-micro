package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

func (app *Config) logViaHttpRequest(w http.ResponseWriter, l LogPayload) {
	jsonData, _ := json.Marshal(l)

	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("error calling log service"))
		return
	}

	app.writeJSON(w, http.StatusOK, jsonRes{
		Error:   false,
		Message: "logged",
	})
}
