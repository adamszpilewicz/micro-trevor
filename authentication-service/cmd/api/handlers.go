package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Auth(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errJSON(w, err, http.StatusBadRequest)
		return
	}

	log.Println(payload)

	user, err := app.models.GetByEmail(context.Background(), payload.Email)
	if err != nil {
		app.errJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	log.Println(user)

	valid, err := user.PasswordMatches(payload.Password)

	log.Println(valid)

	if err != nil || !valid {
		app.errJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("logged in user: %s", user.Email))
	if err != nil {
		app.errJSON(w, err)
	}

	payloadToSend := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged user %s", user.Email),
		Data:    user,
	}

	log.Println(payloadToSend)

	app.writeJSON(w, payloadToSend, http.StatusAccepted)

}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://log-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	return nil
}
