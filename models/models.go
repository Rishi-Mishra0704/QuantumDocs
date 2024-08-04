package models

type Endpoint struct {
	Method         string
	Path           string
	Description    string
	Parameters     []Parameter
	RequestSchema  string
	ResponseSchema string
}

type Parameter struct {
	Name        string
	Type        string
	Description string
	Required    bool
}

type APIDoc struct {
	Title       string
	Description string
	Version     string
	Endpoints   []Endpoint
}

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
