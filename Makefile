BINARY=devctl-tempo
BIN_DIR=bin
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

.PHONY: test build build-all

test:
	go test ./...

build:
	@mkdir -p $(BIN_DIR)
	GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build -o $(BIN_DIR)/$(BINARY) ./main.go
