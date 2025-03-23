.PHONY: build clean install run help test

# Binary name
BINARY_NAME=togo
# Build directory
BUILD_DIR=./bin

# Main build target
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .

# Clean built artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)

# Install to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@mv $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Installed at $(GOPATH)/bin/$(BINARY_NAME)"

# Install to /usr/local/bin (requires sudo)
install-system: build
	@echo "Installing $(BINARY_NAME) to system..."
	@sudo mv $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "Installed at /usr/local/bin/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
	@echo "Built for Linux (amd64)"
	
	# Windows
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
	@echo "Built for Windows (amd64)"
	
	# macOS
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	@echo "Built for macOS (amd64)"
	
	# ARM for Raspberry Pi etc.
	@GOOS=linux GOARCH=arm GOARM=7 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm
	@echo "Built for Linux (ARM)"

# Default help command
help:
	@echo "Make targets:"
	@echo "  build         - Build the binary"
	@echo "  clean         - Remove built files"
	@echo "  install       - Install to GOPATH/bin"
	@echo "  install-system - Install to /usr/local/bin (requires sudo)"
	@echo "  test          - Run tests"
	@echo "  run           - Build and run locally"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  help          - Show this help"

# Default target
default: build 