.PHONY: build clean install run help test

# Binary name
BINARY_NAME=togo
# Build directory
BUILD_DIR=./bin

# Main build target
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir $(BUILD_DIR) 2> /dev/null || true
	@go build -o ${BUILD_DIR}/$(BINARY_NAME)

# Clean built artifacts
clean:
	@echo "Cleaning..."
	@rm $(BUILD_DIR)/* 2> /dev/null || true

# Install to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@mv ${BUILD_DIR}/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Installed at $(GOPATH)/bin/$(BINARY_NAME)"

# Install to /usr/local/bin (requires sudo)
install-system: build
	@echo "Installing $(BINARY_NAME) to system..."
	@sudo mv ${BUILD_DIR}/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "Installed at /usr/local/bin/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@${BUILD_DIR}/$(BINARY_NAME)

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir $(BUILD_DIR) 2> /dev/null || true
	
	# Linux
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
	@echo "Built for Linux (amd64)"
	
	# Windows
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
	@echo "Built for Windows (amd64)"
	
	# macOS Intel
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	@echo "Built for macOS with Intel chips (amd64)"
	
	# macOS M-series
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64
	@echo "Built for macOS with M-series chips (arm64)"
	
	# ARM for Raspberry Pi etc.
	@GOOS=linux GOARCH=arm64 GOARM=7 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64
	@echo "Built for Linux (arm64)"

# Default help command
help:
	@echo "Make targets:"
	@echo "  build              - Build the binary"
	@echo "  clean              - Remove built files"
	@echo "  install            - Install to GOPATH/bin"
	@echo "  install-system     - Install to /usr/local/bin (requires sudo)"
	@echo "  test               - Run tests"
	@echo "  run                - Build and run locally"
	@echo "  build-all          - Build for multiple platforms"
	@echo "  help               - Show this help"

# Default target
default: build 