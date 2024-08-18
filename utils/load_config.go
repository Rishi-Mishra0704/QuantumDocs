package utils

import (
	"encoding/json"
	"os"

	"github.com/Rishi-Mishra0704/QuantumDocs/models"
)

// LoadConfig reads and parses the JSON configuration file
func LoadConfig() (*models.Config, error) {
	// Use the relative path directly if the file is in the root directory
	filename := "quantumdocs.json"

	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON into the config struct
	var config models.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
