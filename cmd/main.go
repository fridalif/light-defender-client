package main

import (
	"light-defender-client/internal/connector"
	"light-defender-client/pkg/config"
	"log"
)

func main() {
	appConfig, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
	server := connector.NewConnector(appConfig)

	err = server.Run()
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
