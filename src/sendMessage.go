package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type responseBody struct {
	StatusCode int
	Message    string
}

type messageFormSymphony struct {
	StatusCode int
	Message    string
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func sendMessage(botClient *botClient, streamID string, message string) (*responseBody, error) {

	jsonStr := `{"message": "<messageML>` + message + `</messageML>"}`

	req, _ := http.NewRequest(
		"POST",
		"https://"+botClient.config.AgentHost+"/agent/v4/stream/"+streamID+"/message/create",
		bytes.NewBuffer([]byte(jsonStr)),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("sessionToken", botClient.sessionToken)
	req.Header.Set("keyManagerToken", botClient.kmToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Warn(err.Error())
	}

	defer resp.Body.Close()

	var symphonyResp responseBody

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warn(err.Error())
	}
	err = json.Unmarshal(body, &symphonyResp)
	if err != nil {
		log.Warn(err.Error())
	}

	//check status
	switch resp.StatusCode {
	case 200:
		symphonyResp.StatusCode = resp.StatusCode

	// case of Unauthorized
	case 401:
		// reauthenticate
		authenticate(botClient)
		sendMessage(botClient, streamID, message)
	default:
		symphonyResp.StatusCode = resp.StatusCode

	}

	return &symphonyResp, nil

}
