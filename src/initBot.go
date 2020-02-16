package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type config struct {
	SessionAuthHost        string `json:"sessionAuthHost"`
	SessionAuthPort        int    `json:"sessionAuthPort"`
	KeyAuthHost            string `json:"keyAuthHost"`
	KeyAuthPort            int    `json:"keyAuthPort"`
	PodHost                string `json:"podHost"`
	PodPort                int    `json:"podPort"`
	AgentHost              string `json:"agentHost"`
	AgentPort              int    `json:"agentPort"`
	AuthType               string `json:"authType"`
	BotCertPath            string `json:"botCertPath"`
	BotCertName            string `json:"botCertName"`
	BotCertPassword        string `json:"botCertPassword"`
	BotPrivateKeyPath      string `json:"botPrivateKeyPath"`
	BotPrivateKeyName      string `json:"botPrivateKeyName"`
	BotUsername            string `json:"botUsername"`
	BotEmailAddress        string `json:"botEmailAddress"`
	AppCertPath            string `json:"appCertPath"`
	AppCertName            string `json:"appCertName"`
	AppCertPassword        string `json:"appCertPassword"`
	ProxyURL               string `json:"proxyURL"`
	ProxyUsername          string `json:"proxyUsername"`
	ProxyPassword          string `json:"proxyPassword"`
	AuthTokenRefreshPeriod string `json:"authTokenRefreshPeriod"`
}

type authenticateResp struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func initBot(configPath string) (botClientInit *botClient) {

	botClientInit = new(botClient)

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := json.Unmarshal(content, &botClientInit.config); err != nil {
		log.Fatal(err.Error())
	}

	authenticate(botClientInit)

	return botClientInit

}

func authenticate(botClient *botClient) {

	keyData, err := ioutil.ReadFile(botClient.config.BotPrivateKeyPath + botClient.config.BotPrivateKeyName)
	if err != nil {
		log.Fatal(err.Error())
	}

	privateKeyBlock, _ := pem.Decode(keyData)

	// parsing the privateKey from decoded keyData
	secretKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		log.Fatal(err.Error())
	}

	// jwt payload info
	claims := jwt.StandardClaims{
		Subject:   botClient.config.BotUsername,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Unix(time.Now().Unix(), 0).Add(5 * time.Minute).Unix(),
	}
	unSignedToken := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	token, err := unSignedToken.SignedString(secretKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	tokenJSON := `{"token":"` + token + `"}`

	sessionAuthenticate(botClient, tokenJSON)
	kmAuthenticate(botClient, tokenJSON)

	log.Info("the bot is authenticated successfully!")

}

func sessionAuthenticate(botClient *botClient, tokenJSON string) {

	req, _ := http.NewRequest(
		"POST",
		"https://"+botClient.config.SessionAuthHost+":"+strconv.Itoa(botClient.config.SessionAuthPort)+"/login/pubkey/authenticate",
		bytes.NewBuffer([]byte(tokenJSON)),
	)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.WithFields(log.Fields{
			"result": "cannot connect to sessionAuth API. \ncheck if sessionAuthHost or sessionAuthPort in config.json is correct.",
		}).Fatal(err.Error())
	}

	switch resp.StatusCode {
	case 401:
		log.Fatal("Cannot authenticate user, please check if your RSA key pair is correct.")
	case 200:
		log.Info("received session token successfully.")
	default:
		log.Error(resp)
	}

	defer resp.Body.Close()

	data := new(authenticateResp)
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err.Error())
	}

	resp.Body.Close()

	botClient.sessionToken = data.Token
}

func kmAuthenticate(botClient *botClient, tokenJSON string) {

	req, _ := http.NewRequest(
		"POST",
		"https://"+botClient.config.KeyAuthHost+":"+strconv.Itoa(botClient.config.KeyAuthPort)+"/relay/pubkey/authenticate",
		bytes.NewBuffer([]byte(tokenJSON)),
	)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.WithFields(log.Fields{
			"result": "cannot connect to keyManagerAuth API. \ncheck if keyAuthHost or keyAuthPort in config.json is correct.",
		}).Fatal(err.Error())
	}

	switch resp.StatusCode {
	case 401:
		log.Fatal("Cannot authenticate user, please check if your RSA key pair is correct.")
	case 200:
		log.Info("received keymanager token successfully.")
	default:
		log.Error(resp)
	}

	defer resp.Body.Close()

	data := new(authenticateResp)
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err.Error())
	}

	botClient.kmToken = data.Token
}
