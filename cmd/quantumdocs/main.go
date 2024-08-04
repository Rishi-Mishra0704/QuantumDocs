package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Rishi-Mishra0704/QuantumDocs/server"
	"github.com/fsnotify/fsnotify"
)

// Config struct to hold the configuration values
type Config struct {
	APIFilePath string     `json:"apiFilePath"`
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

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}

	defer watcher.Close()

	done := make(chan bool)

	// Watch API source file and config file
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("File modified:", event.Name)
					// Reload config if the config file was modified
					if event.Name == "quantumdocs.json" {
						newConfig, err := loadConfig("quantumdocs.json")
						if err != nil {
							log.Printf("Error reloading config: %v", err)
						} else {
							config = newConfig
						}
					}
					err := generateAPIDocumentation(config)
					if err != nil {
						log.Printf("Error regenerating API documentation: %v", err)
					} else {
						fmt.Println("API documentation updated successfully!")
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(config.APIFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Add("quantumdocs.json")
	if err != nil {
		log.Fatal(err)
	}

	// Initial documentation generation
	err = generateAPIDocumentation(config)
	if err != nil {
		log.Fatalf("Error generating API documentation: %v", err)
	}

	// Serve the documentation
	http.HandleFunc("/quantumdocs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(server.GetHTML()))
	})

	fmt.Println("Serving documentation at http://localhost:8080/quantumdocs")
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	<-done
}

// generateAPIDocumentation processes the API file and updates the in-memory HTML content
func generateAPIDocumentation(config *Config) error {
	apiDoc, err := server.ParseAPIDoc(config.APIFilePath)
	if err != nil {
		return fmt.Errorf("error parsing API documentation: %w", err)
	}

	// Set API documentation metadata from the config
	apiDoc.Title = config.APIDoc.Title
	apiDoc.Description = config.APIDoc.Description
	apiDoc.Version = config.APIDoc.Version

	err = server.GenerateAPIDocs(apiDoc, "quantumdocs", "index.html")
	if err != nil {
		return fmt.Errorf("error generating API documentation: %w", err)
	}

	return nil
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
