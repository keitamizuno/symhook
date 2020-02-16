package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDataFeed(t *testing.T) {

	//Given: successfuly get session token and km token
	botClient := initBot("../config/config.json")

	//When: run createDatafeedEventsService() in main.go
	datafeedID := createDatafeedEventsService(botClient)

	//Then: datafeedID is not nil
	assert.NotNil(t, datafeedID)

}
