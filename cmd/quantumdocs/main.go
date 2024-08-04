package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Rishi-Mishra0704/QuantumDocs/server"
)

// Config struct to hold the configuration values
type Config struct {
	APIFilePath string     `json:"apiFilePath"`
	OutputPath  string     `json:"outputPath"`
	APIDoc      APIDocMeta `json:"apiDoc"`
}

// APIDocMeta struct to hold API documentation metadata
type APIDocMeta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

func main() {
	// Load configuration from JSON file
	config, err := loadConfig("quantumdocs.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Construct absolute paths
	absApiPath := config.APIFilePath

	apiDoc, err := server.ParseAPIDoc(absApiPath)
	if err != nil {
		log.Fatalf("Error parsing API documentation: %v", err)
	}

	// Set API documentation metadata from the config
	apiDoc.Title = config.APIDoc.Title
	apiDoc.Description = config.APIDoc.Description
	apiDoc.Version = config.APIDoc.Version

	err = server.GenerateAPIDocs(apiDoc, "quantumdocs", "index.html")
	if err != nil {
		log.Fatalf("Error generating API documentation: %v", err)
	}

	fmt.Println("API documentation generated successfully!")
}

// loadConfig reads and parses the JSON configuration file
func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
