package server

import (
	"net/http"

	"github.com/Rishi-Mishra0704/QuantumDocs/models"
	"github.com/Rishi-Mishra0704/QuantumDocs/template"
)

var apiDoc *models.APIDoc

// InitializeAPIDoc sets up the API documentation
func InitializeAPIDoc(doc *models.APIDoc) {
	apiDoc = doc
}

// ServeAPIDocs handles the HTTP request and renders the API documentation
func ServeAPIDocs(w http.ResponseWriter, r *http.Request) {
	if apiDoc == nil {
		http.Error(w, "API documentation not initialized", http.StatusInternalServerError)
		return
	}

	// Render the HTML directly to the response writer
	err := template.GenerateHtml(apiDoc).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering API documentation", http.StatusInternalServerError)
	}
}

// SetupRoutes configures the HTTP routes for the server
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/quantumdocs", ServeAPIDocs)
	return mux
}

// StartServer starts the HTTP server
func StartServer(addr string) error {
	return http.ListenAndServe(addr, SetupRoutes())
}
