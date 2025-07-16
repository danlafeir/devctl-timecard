APP_NAME=devctl-tempo
STANDALONE_NAME=tempo
BUILD_DIR=bin
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Default values (can be overridden)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

.PHONY: all build build-all build-standalone build-all-standalone clean test

# Usage:
#   make build           # builds for your current system, output: bin/devctl-tempo-<os>-<arch>-<hash>
#   make build-all       # builds for all supported systems
#   make build-standalone # builds standalone binary named 'tempo' for your current system
#   make build-all-standalone # builds standalone binaries for all supported systems
#   GOOS=linux GOARCH=amd64 make build  # cross-compiles for linux/amd64

all: build

build:
	@echo "Building $(APP_NAME) for $(GOOS)/$(GOARCH)..."
	@rm -rf $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)
	GIT_HASH=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(APP_NAME)-$(GOOS)-$(GOARCH)-$$GIT_HASH ./main.go

build-all:
	@echo "Building $(APP_NAME) for all supported OS/ARCH combinations..."
	@rm -rf $(BUILD_DIR)/$(APP_NAME)*
	@mkdir -p $(BUILD_DIR)
	GIT_HASH=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64-$$GIT_HASH ./main.go; \
	GOOS=linux GOARCH=arm64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64-$$GIT_HASH ./main.go; \
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64-$$GIT_HASH ./main.go; \
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64-$$GIT_HASH ./main.go

build-standalone:
	@echo "Building standalone tempo binary for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GIT_HASH=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/tempo ./main.go

build-all-standalone:
	@echo "Building standalone tempo binaries for all supported OS/ARCH combinations..."
	@rm -rf $(BUILD_DIR)/$(STANDALONE_NAME)*
	@mkdir -p $(BUILD_DIR)
	GIT_HASH=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(STANDALONE_NAME)-linux-amd64 ./main.go; \
	GOOS=linux GOARCH=arm64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(STANDALONE_NAME)-linux-arm64 ./main.go; \
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(STANDALONE_NAME)-darwin-amd64 ./main.go; \
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(STANDALONE_NAME)-darwin-arm64 ./main.go

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...
