package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	user, err := app.Repo.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	valid, err := app.Repo.PasswordMatches(requestPayload.Password, *user)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		log.Println("log authen failed", err)
	}

	app.writeJSON(w, http.StatusOK, jsonRes{
		Error:   false,
		Message: "auth successful",
		Data:    user,
	})
}

func (app *Config) logRequest(name, data string) error {
	entry := struct {
		Name string
		Data string
	}{
		Name: name,
		Data: data,
	}

	jsonData, _ := json.Marshal(entry)

	logServiceUrl := "http://logger-service/log"
	req, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("request logging failed")
	}

	return nil
}
