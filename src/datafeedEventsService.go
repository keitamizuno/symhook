package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type datafeedID struct {
	ID string `json:"id"`
}

type datafeedContents []struct {
	ID        string `json:"id"`
	MessageID string `json:"messageId"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"`
	Initiator struct {
		User struct {
			UserID      int64  `json:"userId"`
			FirstName   string `json:"firstName"`
			LastName    string `json:"lastName"`
			DisplayName string `json:"displayName"`
			Email       string `json:"email"`
			Username    string `json:"username"`
		} `json:"user"`
	} `json:"initiator"`
	Payload struct {
		MessageSent struct {
			Message struct {
				MessageID string `json:"messageId"`
				Timestamp int64  `json:"timestamp"`
				Message   string `json:"message"`
				Data      string `json:"data"`
				User      struct {
					UserID      int64  `json:"userId"`
					FirstName   string `json:"firstName"`
					LastName    string `json:"lastName"`
					DisplayName string `json:"displayName"`
					Email       string `json:"email"`
					Username    string `json:"username"`
				} `json:"user"`
				Stream struct {
					StreamID   string `json:"streamId"`
					StreamType string `json:"streamType"`
				} `json:"stream"`
				ExternalRecipients bool   `json:"externalRecipients"`
				UserAgent          string `json:"userAgent"`
				OriginalFormat     string `json:"originalFormat"`
			} `json:"message"`
		} `json:"messageSent"`
	} `json:"payload"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func initDatafeedEventsService(botClient *botClient) {

	//createDatafeedEvent
	datafeedID := createDatafeedEventsService(botClient)

	readDatafeedEventService(botClient, datafeedID)

}

func createDatafeedEventsService(botClient *botClient) string {

	req, _ := http.NewRequest(
		"POST",
		"https://"+botClient.config.AgentHost+"/agent/v4/datafeed/create",
		nil,
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("sessionToken", botClient.sessionToken)
	req.Header.Set("keyManagerToken", botClient.kmToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	jsonBytes := ([]byte)(byteArray)
	data := new(datafeedID)
	if err := json.Unmarshal(jsonBytes, data); err != nil {
		log.Fatal(err.Error())
	}

	return data.ID

}

func readDatafeedEventService(botClient *botClient, datafeedID string) {

	req, _ := http.NewRequest(
		"GET",
		"https://"+botClient.config.AgentHost+"/agent/v4/datafeed/"+datafeedID+"/read",
		nil,
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("sessionToken", botClient.sessionToken)
	req.Header.Set("keyManagerToken", botClient.kmToken)

	client := &http.Client{}

	log.Info("datafeedEvent is started.")

	// datafeed loop
	for {
		resp, err := client.Do(req)
		if err != nil {
			log.Error(err.Error())
		}

		defer resp.Body.Close()

		data := new(datafeedContents)

		switch resp.StatusCode {
		case 401:
			authenticate(botClient)
		case 204:
			// log.Info("No content found")
		case 200:
			if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
				log.Warn(err.Error())
			} else {
				for _, d := range *data {
					if d.Type == "MESSAGESENT" {
						onMessage(botClient, d)
					}
				}
			}
		default:
			// dumpResp, _ := httputil.DumpResponse(resp, true)
			log.Error(resp)

		}

	}

}
