package server

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Rishi-Mishra0704/QuantumDocs/models"
	templates "github.com/Rishi-Mishra0704/QuantumDocs/template"
)

// GenerateAPIDocs generates API documentation and writes it to the specified output file
func GenerateAPIDocs(apiDoc *models.APIDoc, outputDir, outputFile string) error {
	component := templates.ApiDocTemplate(apiDoc)

	// Ensure the output directory exists
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	outputPath := filepath.Join(outputDir, outputFile)

	// Create or overwrite the output file
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer f.Close()

	// Render the documentation into the file
	return component.Render(context.Background(), f)
}
