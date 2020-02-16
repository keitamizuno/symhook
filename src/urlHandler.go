package main

import (
	"strings"
)

func urlHandler(messageText string) string {

	// FIXME: error handling

	messageText = parseURL(messageText, "http:")
	MLformatMessage := parseURL(messageText, "https:")

	return MLformatMessage

}

func parseURL(messageText string, urlHead string) string {

	var MLformatMessage string
	var targetText string
	var displayText string
	var URLHeadPosition int
	var URLEndPosition int
	var totalPosition int

	for strings.Contains(messageText[totalPosition:], urlHead) {

		URLHeadPosition = strings.Index(messageText[totalPosition:], urlHead)
		targetText = messageText[totalPosition+URLHeadPosition:]

		// case of mrkdwn http(s) tag
		// eg.) <http(s):xxxx.com|this is a link>
		if messageText[totalPosition+URLHeadPosition-1:totalPosition+URLHeadPosition] == "<" {
			MLformatMessage += messageText[totalPosition:][:URLHeadPosition-1]
			URLEndPosition := strings.Index(targetText, ">")
			targetText = targetText[:URLEndPosition]
			totalPosition += URLHeadPosition + URLEndPosition + 1

			if vbarPosition := strings.Index(targetText, "|"); vbarPosition != -1 {
				displayText = targetText[vbarPosition+1:]
				targetText = targetText[:vbarPosition]
			} else {
				displayText = targetText
			}

			MLformatMessage += "<a href='" + targetText + "'>" + displayText + "</a>"

			// case of MLformat (a tag
			// eg.) <a href="http(s):xxxx.com">xxxx.com</a>
		} else if strings.Contains(messageText[totalPosition+URLHeadPosition-7:totalPosition+URLHeadPosition], "href=") {
			URLEndPosition = strings.Index(targetText, "</a>")
			MLformatMessage += messageText[totalPosition:][:URLHeadPosition+URLEndPosition]
			totalPosition += URLHeadPosition + URLEndPosition

			// case of normal url
			// eg.) http(s):xxxx.com
		} else {
			MLformatMessage += messageText[totalPosition:][:URLHeadPosition]
			if URLEndPosition = strings.Index(targetText, " "); URLEndPosition != -1 {
				targetText = targetText[:URLEndPosition]
				totalPosition += URLHeadPosition + URLEndPosition
			} else {
				totalPosition = len(messageText)
			}

			MLformatMessage += "<a href='" + targetText + "'>" + targetText + "</a>"
		}
	}

	MLformatMessage += messageText[totalPosition:]

	if MLformatMessage == "" {
		MLformatMessage = messageText
	}

	return MLformatMessage
}
