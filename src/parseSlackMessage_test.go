package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSlackMessage(t *testing.T) {

	// got json message formatted slack attachment.
	slackAttachmentJSON := `
		   {
			  "fallback":"fallback Test",
			  "pretext":"attachments Test",
			  "color":"#0000FF",
			  "fields":[
				 {
					"title":"attachment01",
					"value":"This is attachment"
				 }
			  ]
		   }`

	symphonyMessage :=
		`attachments Test<br/><card accent=\"tempo-bg-color--blue\"><b>attachment01</b><br/>string<br/></card>`

	var slackAttachment slackAttachment

	err := json.Unmarshal([]byte(slackAttachmentJSON), &slackAttachment)
	if err != nil {
		fmt.Println(err)
		return
	}

	//
	mlFormatMessage, err := parseSlackMessage(slackAttachment)

	assert.Equal(t, symphonyMessage, mlFormatMessage)

}

func TestParceCardColor(t *testing.T) {

	hexColorCodeWhite := "#FFFFFF"
	hexColorCodeBlack := "#000000"
	hexColorCodeBule := "#0000FF"
	hexColorCodeYellow := "#FFFF00"
	hexColorCodeGreen := "#008000"
	hexColorCodePurple := "#800080"
	hexColorCodeRed := "#FF0000"

	color, _ := parseCardColor(hexColorCodeWhite)
	assert.Equal(t, "white", color)

	color, _ = parseCardColor(hexColorCodeBlack)
	assert.Equal(t, "black", color)

	color, _ = parseCardColor(hexColorCodeBule)
	assert.Equal(t, "blue", color)

	color, _ = parseCardColor(hexColorCodeYellow)
	assert.Equal(t, "yellow", color)

	color, _ = parseCardColor(hexColorCodeGreen)
	assert.Equal(t, "green", color)

	color, _ = parseCardColor(hexColorCodePurple)
	assert.Equal(t, "purple", color)

	color, _ = parseCardColor(hexColorCodeRed)
	assert.Equal(t, "red", color)

}
