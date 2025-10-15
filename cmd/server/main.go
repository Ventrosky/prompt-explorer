package main

import (
	"log"

	"github.com/ventrosky/prompt-explorer/app"
)

func main() {
	// Create and start the application
	application := app.New()

	if err := application.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
