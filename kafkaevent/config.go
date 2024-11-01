package kafkaevent

import (
	"log"
	"os"
)

// Config holds all the configuration variables
type Config struct {
	BROKER_URLS string
	TOPIC       string
	GROUP_ID    string
}

// ServiceConfig loads environment variables and returns a Config struct
func KafkaConfig() *Config {
	config := &Config{
		BROKER_URLS: os.Getenv("BROKER_URLS"),
		TOPIC:       os.Getenv("TOPIC"),
		GROUP_ID:    os.Getenv("GROUP_ID"),
	}

	log.Println("In config init, Zepto API URL:", config.BROKER_URLS)
	log.Println("In config init, Topic is:", config.TOPIC)

	return config
}
