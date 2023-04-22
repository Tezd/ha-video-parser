package main

import (
	"log"

	"ha-video-parser/cmd/server"
	"ha-video-parser/pkg/service"
)

func main() {
	handler := service.New()
	handler.RegisterHttpHandlers()

	if err := server.New(8080).ListenAndServe(); err != nil {
		log.Fatalf("failed to close server: [%v]", err)
	}
}
