package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Rishi-Mishra0704/QuantumDocs/models"
	"github.com/Rishi-Mishra0704/QuantumDocs/server"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
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
	clients := make(map[*websocket.Conn]bool)
	broadcast := make(chan string)

	// WebSocket server to handle client connections and broadcasting messages
	go func() {
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println("Error while connecting:", err)
				return
			}
			clients[conn] = true
			defer conn.Close()
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					// Handle WebSocket close codes gracefully
					if websocket.IsCloseError(err, websocket.CloseGoingAway) {
						log.Println("WebSocket client disconnected")
						return
					}
					delete(clients, conn)
					log.Println("WebSocket client disconnected due to error:", err)
					break
				}
			}
		})

		for {
			msg := <-broadcast
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					client.Close()
					delete(clients, client)
					log.Printf("Error sending message to client: %v", err)
				} else {
					return
				}
			}
		}
	}()

	// Watch API source file and config file
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Printf("File modified: %s\n", event.Name)
					if event.Name == "quantumdocs.json" {
						newConfig, err := loadConfig("quantumdocs.json")
						if err != nil {
							return
						} else {
							config = newConfig
						}
					} else if event.Name == config.APIFilePath {
						err := generateAPIDocumentation(config)
						if err != nil {
							return
						} else {
							broadcast <- "reload"
						}
					} else {
						fmt.Printf("Ignored file modification: %s\n", event.Name)
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
		log.Fatalf("Error adding API file to watcher: %v", err)
	}

	err = watcher.Add("quantumdocs.json")
	if err != nil {
		log.Fatalf("Error adding config file to watcher: %v", err)
	}

	// Initial documentation generation
	err = generateAPIDocumentation(config)
	if err != nil {
		log.Fatalf("Error generating API documentation: %v", err)
	}

	// Serve the documentation and WebSocket server
	http.HandleFunc("/quantumdocs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(server.GetHTML()))
	})

	fmt.Printf("Serving documentation at %s/quantumdocs", config.BaseUrl)
	go func() {
		log.Fatal(http.ListenAndServe(config.Port, nil))
	}()

	<-done
}

// generateAPIDocumentation processes the API file and updates the in-memory HTML content
func generateAPIDocumentation(config *models.Config) error {
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
func loadConfig(filename string) (*models.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config models.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
