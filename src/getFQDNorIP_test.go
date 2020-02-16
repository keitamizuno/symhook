package main

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFQDN(t *testing.T) {

	// get IP addr for constract a webhook url
	webhookIP := getFQDNorIP()

	// webhookIP is IP
	rep := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	isIP := rep.MatchString(webhookIP)

	assert.True(t, isIP)

	// set IP address to environment variable
	os.Setenv("FQDN_IP", "192.11.0.1")
	webhookIP = getFQDNorIP()

	rep = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	isIP = rep.MatchString(webhookIP)

	// clear the environment variable
	os.Unsetenv("FQDN_IP")

	assert.True(t, isIP)

	// set FQDN address to environment variable
	os.Setenv("FQDN_IP", "test-symphony-webhook.com")

	webhookFQDN := getFQDNorIP()

	rep = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$`)
	isFQDN := rep.MatchString(webhookFQDN)

	// clear the environment variable
	os.Unsetenv("FQDN_IP")

	assert.True(t, isFQDN)
}
