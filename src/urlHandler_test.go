package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURLHandler(t *testing.T) {

	//http
	httpMessage := `Hi, this is handling HTTP TAG test. this is a mrkdwn http tag. <http://twitter.com/KeitaMizuno2|this is my twitter account> And this is a normal url. http://twitter.com/KeitaMizuno2`

	expectHTTPMLformat := `Hi, this is handling HTTP TAG test. this is a mrkdwn http tag. <a href='http://twitter.com/KeitaMizuno2'>this is my twitter account</a> And this is a normal url. <a href='http://twitter.com/KeitaMizuno2'>http://twitter.com/KeitaMizuno2</a>`

	ParsedHTTPMessage := urlHandler(httpMessage)
	assert.Equal(t, expectHTTPMLformat, ParsedHTTPMessage)

	//https
	HTTPSMessage := `Hi, this is handling HTTPS TAG test. this is a mrkdwn http tag. <https://twitter.com/KeitaMizuno2> And this is a normal url. https://twitter.com/KeitaMizuno2`

	expectHTTPSMLformat := `Hi, this is handling HTTPS TAG test. this is a mrkdwn http tag. <a href='https://twitter.com/KeitaMizuno2'>https://twitter.com/KeitaMizuno2</a> And this is a normal url. <a href='https://twitter.com/KeitaMizuno2'>https://twitter.com/KeitaMizuno2</a>`

	parsedHTTPSMessage := urlHandler(HTTPSMessage)
	assert.Equal(t, expectHTTPSMLformat, parsedHTTPSMessage)

	//MLFormatMessage
	//In this case, urlHandler() does nothing even though a text includes 'http:'.
	orifinalMlFormatMessage := `Hi, this is MLFormat test. this is a mrkdwn http tag. <a href='http://twitter.com/KeitaMizuno2'>http://twitter.com/KeitaMizuno2</a>`

	expectMLformat := `Hi, this is MLFormat test. this is a mrkdwn http tag. <a href='http://twitter.com/KeitaMizuno2'>http://twitter.com/KeitaMizuno2</a>`

	notParsedMessage := urlHandler(orifinalMlFormatMessage)
	assert.Equal(t, expectMLformat, notParsedMessage)

	//mixed
	mixMessage := `Hi, this is handling mixed test. this is a mrkdwn http tag. <http://twitter.com/KeitaMizuno2> And this is a normal url. https://twitter.com/KeitaMizuno2 <a href='http://twitter.com/KeitaMizuno2'>http://twitter.com/KeitaMizuno2</a>`

	expectMixMLformat := `Hi, this is handling mixed test. this is a mrkdwn http tag. <a href='http://twitter.com/KeitaMizuno2'>http://twitter.com/KeitaMizuno2</a> And this is a normal url. <a href='https://twitter.com/KeitaMizuno2'>https://twitter.com/KeitaMizuno2</a> <a href='http://twitter.com/KeitaMizuno2'>http://twitter.com/KeitaMizuno2</a>`

	parsedMixMessage := urlHandler(mixMessage)
	assert.Equal(t, expectMixMLformat, parsedMixMessage)

}
