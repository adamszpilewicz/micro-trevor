package main

import (
	"log"
	"log-service/db"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	//	insert data
	event := db.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	log.Printf("---> trying to write log entry: %s", event)
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errJSON(w, err)
		return
	}

	res := jsonResponse{
		Error:   false,
		Message: "logged message",
	}

	app.writeJSON(w, res, http.StatusAccepted)
}
