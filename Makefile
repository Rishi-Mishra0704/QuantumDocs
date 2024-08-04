APP_NAME := quantumdocs
API_FILE := test_api.go
OUTPUT_FILE := api_docs.html

.PHONY: build run

build: 
	@templ generate
	@go build -o bin/ ./cmd/...

run: build
	@./bin/$(APP_NAME) -api $(API_FILE) -output $(OUTPUT_FILE)