package main

import (
	log "github.com/sirupsen/logrus"
	"net"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

// getFQDNorIP  Get preferred outbound FQDN or IP
func getFQDNorIP() string {

	FQDNorIP := os.Getenv("FQDN_IP")
	if FQDNorIP != "" {
		return FQDNorIP
	}

	// in case that FQDN_IP is not set, the private ip will be returned.
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()

}
