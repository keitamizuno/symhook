package main

import (
	"fmt"
	"image/color"
	"os"
	"reflect"

	log "github.com/sirupsen/logrus"
)

const breakTag = "<br/>"

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func parseSlackMessage(slackMessage slackAttachment) (string, error) {

	var mlFormatMessage string

	if slackMessage.Pretext != "" {
		mlFormatMessage += slackMessage.Pretext + breakTag
	}

	if slackMessage.Text != "" || len(slackMessage.Fields) > 0 {

		cardColor, err := parseCardColor(slackMessage.Color)
		if err != nil {
			log.Info(err.Error())
			return "", err
		}

		mlFormatMessage += `<card accent=\"tempo-bg-color--` + cardColor + `\">`

		if slackMessage.Text != "" {
			mlFormatMessage += slackMessage.Text + breakTag
		}

		for _, value := range slackMessage.Fields {
			if value.Title != "" {
				mlFormatMessage += "<b>" + value.Title + "</b>" + breakTag
			}
			if value.Value != "" {
				mlFormatMessage += reflect.TypeOf(value.Value).String() + breakTag
			}
		}
		mlFormatMessage += `</card>`
	}

	return mlFormatMessage, nil

}

func parseCardColor(hexColorCode string) (string, error) {

	colorRGB, err := parseHexColor(hexColorCode)
	if err != nil {
		return "", err
	}

	var isRedStrong bool
	var isGreenStrong bool
	var isBlueStrong bool

	if colorRGB.R > 120 {
		isRedStrong = true
	}
	if colorRGB.G > 120 {
		isGreenStrong = true
	}
	if colorRGB.B > 120 {
		isBlueStrong = true
	}

	var cardColor string

	switch {
	case isRedStrong == true && isGreenStrong == true && isBlueStrong == true:
		cardColor = "white"

	case isRedStrong == true && isGreenStrong == true && isBlueStrong == false:
		cardColor = "yellow"

	case isRedStrong == true && isGreenStrong == false && isBlueStrong == true:
		cardColor = "purple"

	case isRedStrong == true && isGreenStrong == false && isBlueStrong == false:
		cardColor = "red"

	case isRedStrong == false && isGreenStrong == true && isBlueStrong == false:
		cardColor = "green"

	case isRedStrong == false && isGreenStrong == false && isBlueStrong == true:
		cardColor = "blue"

	default:
		cardColor = "black"

	}

	return cardColor, nil
}

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("HexColor is invalid length, must be 7 or 4")

	}
	return
}
