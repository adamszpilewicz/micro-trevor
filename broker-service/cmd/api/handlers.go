package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "hit the broker",
	}

	app.writeJSON(w, payload, http.StatusOK)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errJSON(w, errors.New("unknow action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	// create json data to be sent
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// call the service
	logService := "http://log-service/log"
	request, err := http.NewRequest("POST", logService, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"
	app.writeJSON(w, payload, http.StatusAccepted)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	//	 create some json data to be sent to auth service
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	log.Println(string(jsonData))

	//	call the service
	request, err := http.NewRequest("POST", "http://authentication-service/auth", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errJSON(w, err)
		return
	}
	defer response.Body.Close()

	log.Println(response.Body)

	//	make sure we get the correct response code
	if response.StatusCode == http.StatusUnauthorized {
		app.errJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errJSON(w, errors.New("error calling auth service"))
		//app.errJSON(w, errors.New(string(response.Body))
		return
	}

	//	 create variable we'll read response.Body into
	var jsonFromResponse jsonResponse

	//	decode the json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromResponse)
	if err != nil {
		app.errJSON(w, err)
		return
	}

	if jsonFromResponse.Error {
		app.errJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromResponse.Data

	app.writeJSON(w, payload, http.StatusAccepted)

}
