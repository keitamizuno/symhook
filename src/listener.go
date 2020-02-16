package main

import (
	"strings"
)

type datafeedContent struct {
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

const (
	trimLeft  string = "<div data-format=\"PresentationML\" data-version=\"2.0\" class=\"wysiwyg\"><p>"
	trimRight string = `</p></div>`
)

// onMessage is called when datafeedEventsService get a new message from IMs or chat rooms
func onMessage(botClient *botClient, datafeed datafeedContent) {

	// extract text
	messageText := datafeed.Payload.MessageSent.Message.Message
	messageText = strings.TrimLeft(messageText, trimLeft)
	messageText = strings.TrimRight(messageText, trimRight)

	if messageText == "/webhook" {

		FQDNorIP := getFQDNorIP()
		streamID := datafeed.Payload.MessageSent.Message.Stream.StreamID
		webhookURL := "http://" + FQDNorIP + ":8445/symphony-hooks/" + streamID

		sendMessage(botClient, streamID, webhookURL)

	}
}
