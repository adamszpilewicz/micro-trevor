package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct {
}

func (app *Config) sendMail(w http.ResponseWriter, mail MailPayload) {
	jsonData, _ := json.MarshalIndent(mail, "", "\t")

	//	call the service
	mailServiceURL := "http://mail-service/send"

	//	POST mail service
	req, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errJSON(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.errJSON(w, err)
		return
	}
	defer res.Body.Close()

	//	check the correct response
	if res.StatusCode != http.StatusAccepted {
		app.errJSON(w, errors.New("error calling mail service"))
		return
	}

	//	send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to: " + mail.To
	app.writeJSON(w, payload, http.StatusAccepted)
}

func main() {
	app := Config{}

	log.Printf("starting broker at the port %s", webPort)

	// starting server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	log.Fatal(err)
}
