package main

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	// Fetch initial data
	err := weatherInit()
	if err != nil {
		log.Fatalf("failed to fetch initial data: %s", err)
	}

	// Connect to MQTT
	err = mqttInit()
	if err != nil {
		log.Fatalf("failed to connect to mqtt: %s", err)
	}

	// Run HTTP server
	runHTTP()

	// Wait for exit
	ch := make(chan bool)
	<-ch
}
