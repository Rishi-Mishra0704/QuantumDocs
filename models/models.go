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
