package server

import (
	"context"
	"os"

	"github.com/Rishi-Mishra0704/QuantumDocs/models"
	templates "github.com/Rishi-Mishra0704/QuantumDocs/template"
)

func GenerateAPIDocs(apiDoc *models.APIDoc, outputPath string) error {
	component := templates.ApiDocTemplate(apiDoc)
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return component.Render(context.Background(), f)
}
