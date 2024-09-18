package main

import (
	"log"
	"os"
	"tx-parser/internal/app"
)

func main() {
	// Load the config path from an environment variable or use a default path
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	// Initialize the application with the config path
	application, err := app.NewApp(configPath)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Run the application
	if err := application.Run(); err != nil {
		log.Fatalf("Application encountered an error: %v", err)
	}
}
