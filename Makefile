.DEFAULT_GOAL := build
BIN_FILE=crud

build:
	@go build -o "cmd/crud/crud" crud/cmd/crud

run:
	@cd cmd/crud; ./crud

clean:
	@cd cmd/crud; go clean;

test:
	go test ./...