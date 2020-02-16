package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type botClient struct {
	sessionToken string
	kmToken      string
	config
}

type incomingWebhookRequest struct {
	Text        string             `json:"text"`
	Username    string             `json:"username"`
	IconURL     string             `json:"icon_url"`
	ChannelName string             `json:"channel"`
	Props       stringInterface    `json:"props"`
	Attachments []*slackAttachment `json:"attachments"`
	Type        string             `json:"type"`
	IconEmoji   string             `json:"icon_emoji"`
}

type stringInterface map[string]interface{}

type slackAttachment struct {
	ID         int64                   `json:"id"`
	Fallback   string                  `json:"fallback"`
	Color      string                  `json:"color"`
	Pretext    string                  `json:"pretext"`
	AuthorName string                  `json:"author_name"`
	AuthorLink string                  `json:"author_link"`
	AuthorIcon string                  `json:"author_icon"`
	Title      string                  `json:"title"`
	TitleLink  string                  `json:"title_link"`
	Text       string                  `json:"text"`
	Fields     []*slackAttachmentField `json:"fields"`
	ImageURL   string                  `json:"image_url"`
	ThumbURL   string                  `json:"thumb_url"`
	Footer     string                  `json:"footer"`
	FooterIcon string                  `json:"footer_icon"`
	Timestamp  interface{}             `json:"ts"` // This is either a string or an int64
	Actions    []*postAction           `json:"actions,omitempty"`
}

type slackAttachmentField struct {
	Title string              `json:"title"`
	Value interface{}         `json:"value"`
	Short slackCompatibleBool `json:"short"`
}

type slackCompatibleBool bool

type postAction struct {
	// A unique Action ID. If not set, generated automatically.
	ID string `json:"id,omitempty"`

	// The type of the interactive element. Currently supported are
	// "select" and "button".
	Type string `json:"type,omitempty"`

	// The text on the button, or in the select placeholder.
	Name string `json:"name,omitempty"`

	// If the action is disabled.
	Disabled bool `json:"disabled,omitempty"`

	// DataSource indicates the data source for the select action. If left
	// empty, the select is populated from Options. Other supported values
	// are "users" and "channels".
	DataSource string `json:"data_source,omitempty"`

	// Options contains either the buttons that will be displayed on the post
	// or the values listed in a select dropdowon on the post.
	Options []*postActionOptions `json:"options,omitempty"`

	// DefaultOption contains the option, if any, that will appear as the
	// default selection in a select box. It has no effect when used with
	// other types of actions.
	DefaultOption string `json:"default_option,omitempty"`

	// Defines the interaction with the backend upon a user action.
	// Integration contains Context, which is private plugin data;
	// Integrations are stripped from Posts when they are sent to the
	// client, or are encrypted in a Cookie.
	Integration *postActionIntegration `json:"integration,omitempty"`
	Cookie      string                 `json:"cookie,omitempty" db:"-"`
}

type postActionOptions struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

type postActionIntegration struct {
	URL     string                 `json:"url,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func (botClient *botClient) postWebhook(w http.ResponseWriter, r *http.Request) {

	clientIP := getIP(r)
	log.Info("webhook got a new request from " + clientIP)

	params := mux.Vars(r)
	var incomingWebhookRequest incomingWebhookRequest

	contentType := r.Header.Get("Content-Type")
	switch strings.Split(contentType, "; ")[0] {
	case "application/x-www-form-urlencoded":
		payload := strings.NewReader(r.FormValue("payload"))
		log.Info("a original message  : " + r.FormValue("payload"))
		err := json.NewDecoder(payload).Decode(&incomingWebhookRequest)
		if err != nil {
			log.Warn(err.Error())
		}
	case "application/json":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Warn(err.Error())
		}
		log.Info("a original message  : " + string(body))
		err = json.Unmarshal(body, &incomingWebhookRequest)
		if err != nil {
			log.Warn(err.Error())
		}
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"an error occured in webhook": "Content-Type should be application/x-www-form-urlencoded or application/json"}`)
	}

	var sendMessageText string

	if incomingWebhookRequest.Text != "" {
		sendMessageText += urlHandler(incomingWebhookRequest.Text)
	}
	if len(incomingWebhookRequest.Attachments) > 0 {
		if sendMessageText != "" {
			sendMessageText += breakTag + breakTag
		}
		for _, value := range incomingWebhookRequest.Attachments {
			mlFormatMessage, err := parseSlackMessage(*value)
			if err != nil {
				webhookResopse(w, http.StatusBadRequest, err.Error())
				return
			}
			sendMessageText += urlHandler(mlFormatMessage)
		}
	}

	log.Info("the message converted to MLformat : " + sendMessageText)

	responseBody := new(responseBody)
	responseBody, err := sendMessage(botClient, params["streamId"], sendMessageText)
	if err != nil {
		webhookResopse(w, http.StatusBadRequest, err.Error())
	}

	if responseBody.StatusCode != http.StatusOK {
		webhookResopse(w, http.StatusBadRequest, responseBody.Message)
	} else {
		webhookResopse(w, http.StatusOK, responseBody.Message)
	}

}

// getIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func getIP(r *http.Request) string {

	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func webhookResopse(w http.ResponseWriter, respStatus int, respMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(respStatus)
	io.WriteString(w, `{"MessageFromSymphony": "`+respMessage+`"}`)
}

func initWebhook(botClient *botClient) {
	router := mux.NewRouter()
	router.HandleFunc("/symphony-hooks/{streamId}", botClient.postWebhook).Methods("POST")
	log.Info("webhook is listening on port 8445.")
	http.ListenAndServe(":8445", router)
}
