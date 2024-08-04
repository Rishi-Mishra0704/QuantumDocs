package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Rishi-Mishra0704/QuantumDocs/server"
)

func main() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	// Define flags for the API file path and output path
	apiFilePath := flag.String("api", "test_api.go", "Path to the API file")
	outputPath := flag.String("output", "api_docs.html", "Path for the output HTML file")
	flag.Parse()

	// Construct absolute paths
	absApiPath := filepath.Join(cwd, *apiFilePath)
	absOutputPath := filepath.Join(cwd, *outputPath)

	apiDoc, err := server.ParseAPIDoc(absApiPath)
	if err != nil {
		log.Fatalf("Error parsing API documentation: %v", err)
	}

	apiDoc.Title = "QuantumDocs API"
	apiDoc.Description = "This is the API documentation for QuantumDocs"
	apiDoc.Version = "1.0.0"

	err = server.GenerateAPIDocs(apiDoc, absOutputPath)
	if err != nil {
		log.Fatalf("Error generating API documentation: %v", err)
	}

	fmt.Println("API documentation generated successfully!")
}
