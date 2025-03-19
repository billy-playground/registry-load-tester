# Variables
APP_NAME := test
CMD_DIR := ./cmd/test
ASSETS_DIR := ./assets
BUILD_DIR := ./bin
GO_FILES := $(shell find . -type f -name '*.go')

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build: $(BUILD_DIR)/$(APP_NAME)

$(BUILD_DIR)/$(APP_NAME): $(GO_FILES)
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)
	@cp -r $(ASSETS_DIR) $(BUILD_DIR)/$(ASSETS_DIR)

# Run the application with 3 instances
.PHONY: run
run: build
	$(BUILD_DIR)/$(APP_NAME) 3

# Clean up build artifacts
.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)

# Run tests
.PHONY: test
test:
	go test ./...

# Lint the code
.PHONY: lint
lint:
	golangci-lint run

# Format the code
.PHONY: fmt
fmt:
	go fmt ./...

# Check for outdated dependencies
.PHONY: deps
deps:
	go list -u -m all

# Install dependencies
.PHONY: install
install:
	go mod tidy