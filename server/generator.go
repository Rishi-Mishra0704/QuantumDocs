package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Rishi-Mishra0704/QuantumDocs/models"
	"github.com/Rishi-Mishra0704/QuantumDocs/template"
)

// GenerateAPIDocs generates API documentation and writes it to the specified output file
func GenerateAPIDocs(apiDoc *models.APIDoc, outputDir, outputFile string) error {
	htmlContent := template.GenerateHTML(apiDoc)

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

	// Write the HTML content to the file
	_, err = f.WriteString(htmlContent)
	if err != nil {
		return fmt.Errorf("error writing to output file: %v", err)
	}

	// Update the in-memory HTML content
	template.HtmlMu.Lock()
	template.HtmlContent = htmlContent
	template.HtmlMu.Unlock()

	return nil
}

func GetHTML() string {
	template.HtmlMu.RLock()
	defer template.HtmlMu.RUnlock()
	return template.HtmlContent
}
