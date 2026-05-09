package main

import (
	"log"

	"github.com/bigdann09/notifications/internal/config"
	"github.com/bigdann09/notifications/internal/server"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config", err)
	}

	srv := server.NewServer(config)
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal("failed to start server", err)
		}
	}()

	srv.Shutdown()
}
