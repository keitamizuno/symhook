package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitBot(t *testing.T) {

	botClient := initBot(configPATH)

	assert.True(t, botClient != nil)

}
