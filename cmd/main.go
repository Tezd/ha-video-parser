package main

import (
	"log"

	"ha-video-parser/cmd/server"
)

func main() {
	if err := server.New(8080).ListenAndServe(); err != nil {
		log.Fatalf("failed to close server: [%v]", err)
	}
}
