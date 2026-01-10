APP_NAME=timecard
DEVCTL_PLUGIN_NAME=devctl-timecard
BUILD_DIR=bin
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Default values (can be overridden)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

.PHONY: all build build-all build-devctl build-all-devctl clean test

# Usage:
#   make build           				# builds for your current system, output: bin/timecard-<os>-<arch>-<hash>
#   make build-devctl 					# builds devctl plugin binary named 'timecard' for your current system, output: bin/timecard-<os>-<arch>-<hash>
#   make build-all      				# builds for all supported systems, output: bin/timecard-<os>-<arch>-<hash>
#   make build-all-devctl 				# builds devctl plugin binaries for all supported systems
#   GOOS=linux GOARCH=amd64 make build  # cross-compiles for linux/amd64

all: build-all build-all-devctl

build:
	@echo "Building $(APP_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(APP_NAME) ./main.go
	@echo "Built $(APP_NAME) successfully to target /$(BUILD_DIR)"

build-devctl:
	@echo "Building timecard binary for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/timecard ./main.go
	@echo "Built $(DEVCTL_PLUGIN_NAME) successfully to target /$(BUILD_DIR)"

build-all:
	@echo "Building $(APP_NAME) binaries for all supported OS/ARCH combinations..."
	@rm -rf $(BUILD_DIR)/$(APP_NAME)*
	@mkdir -p $(BUILD_DIR)
	@GIT_HASH=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64-$$GIT_HASH ./main.go; \
	GOOS=linux GOARCH=arm64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64-$$GIT_HASH ./main.go; \
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64-$$GIT_HASH ./main.go; \
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64-$$GIT_HASH ./main.go
	@echo "Built $(APP_NAME) successfully to target /$(BUILD_DIR)"

build-all-devctl:
	@echo "Building $(DEVCTL_PLUGIN_NAME) binaries for all supported OS/ARCH combinations..."
	@rm -rf $(BUILD_DIR)/$(DEVCTL_PLUGIN_NAME)*
	@mkdir -p $(BUILD_DIR)
	@GIT_HASH=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(DEVCTL_PLUGIN_NAME)-linux-amd64 ./main.go; \
	GOOS=linux GOARCH=arm64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(DEVCTL_PLUGIN_NAME)-linux-arm64 ./main.go; \
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(DEVCTL_PLUGIN_NAME)-darwin-amd64 ./main.go; \
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.BuildGitHash=$$GIT_HASH' -X 'main.BuildLatestHash=$$GIT_HASH'" -o $(BUILD_DIR)/$(DEVCTL_PLUGIN_NAME)-darwin-arm64 ./main.go
	@echo "Built $(DEVCTL_PLUGIN_NAME) successfully to target /$(BUILD_DIR)"

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...
