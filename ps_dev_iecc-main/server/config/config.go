package config

import (
	"log"
	"os"
)

var GoogleClientID string

func LoadConfig() {
	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	if GoogleClientID == "" {
		log.Fatal("GOOGLE_CLIENT_ID not set")
	}
}
