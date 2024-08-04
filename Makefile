APP_NAME := quantumdocs
.PHONY: build run

build: 
	@templ generate
	@go build -o bin/ ./cmd/...

run: build
	@./bin/$(APP_NAME)