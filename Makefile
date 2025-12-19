.PHONY: build install clean fmt format

BINARY_NAME := togo
BUILD_DIR := ./bin
BUILD_PATH := $(BUILD_DIR)/$(BINARY_NAME)

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_PATH)

install: build
	@echo "Starting installation process..."
	@if [ ! -f "$(BUILD_PATH)" ]; then \
		echo "❌ Error: Binary not found at $(BUILD_PATH)"; \
		exit 1; \
	fi
	@echo "Please select your preferred installation method:"
	@echo "1) Local user install (~/.local/bin) - Recommended for single user"
	@echo "2) System-wide install (/usr/local/bin) - Requires sudo, for all users"
	@read -p "Enter your choice [1-2]: " install_choice; \
	if [[ ! "$$install_choice" =~ ^[12]$$ ]]; then \
		echo "❌ Invalid selection. Please choose 1 or 2."; \
		exit 1; \
	fi; \
	case "$$(uname -s)" in \
		Linux*) OS=Linux ;; \
		Darwin*) OS=macOS ;; \
		*) OS="UNKNOWN" ;; \
	esac; \
	if [[ "$$OS" == "UNKNOWN" ]]; then \
		echo "⚠️  Unsupported operating system detected. Proceeding with basic installation."; \
	fi; \
	case $$install_choice in \
		1) \
			TARGET_DIR="$$HOME/.local/bin"; \
			mkdir -p "$$TARGET_DIR"; \
			if cp "$(BUILD_PATH)" "$$TARGET_DIR/"; then \
				echo "✅ Successfully installed to $$TARGET_DIR/$(BINARY_NAME)"; \
			else \
				echo "❌ Local installation failed. Please check permissions."; \
				exit 1; \
			fi; \
			;; \
		2) \
			TARGET_DIR="/usr/local/bin"; \
			echo "Installing system-wide (requires sudo privileges)..."; \
			if sudo cp "$(BUILD_PATH)" "$$TARGET_DIR/"; then \
				echo "✅ Successfully installed to $$TARGET_DIR/$(BINARY_NAME)"; \
			else \
				echo "❌ System-wide installation failed. Please check sudo permissions."; \
				exit 1; \
			fi; \
			;; \
	esac; \
	CURRENT_SHELL=$$(basename "$$SHELL"); \
	echo "Detected shell: $$CURRENT_SHELL"; \
	INSTALLED_BIN="$$TARGET_DIR/$(BINARY_NAME)"; \
	if [[ "$$install_choice" == "2" ]]; then \
		if [[ "$$CURRENT_SHELL" == "zsh" ]]; then \
			echo "Setting up zsh completion for all users..."; \
			echo "Installing completion script to system-wide location..."; \
			sudo bash -c "$$INSTALLED_BIN completion zsh > \$$(echo \$${fpath[1]}/_togo)"; \
			ROOT_ZSHRC="/root/.zshrc"; \
			if [ ! -f "$$ROOT_ZSHRC" ]; then \
				sudo touch "$$ROOT_ZSHRC"; \
			fi; \
			if ! sudo grep -q "eval \"\$$($$INSTALLED_BIN completion zsh)\"" "$$ROOT_ZSHRC" 2>/dev/null; then \
				echo "eval \"\$$($$INSTALLED_BIN completion zsh)\"" | sudo tee -a "$$ROOT_ZSHRC" >/dev/null; \
			fi; \
			USER_ZSHRC="$$HOME/.zshrc"; \
			if [ ! -f "$$USER_ZSHRC" ]; then \
				touch "$$USER_ZSHRC"; \
			fi; \
			if ! grep -q "eval \"\$$($$INSTALLED_BIN completion zsh)\"" "$$USER_ZSHRC" 2>/dev/null; then \
				echo "eval \"\$$($$INSTALLED_BIN completion zsh)\"" >>"$$USER_ZSHRC"; \
			fi; \
			echo "✅ Zsh completion configured for both root and current user"; \
		elif [[ "$$CURRENT_SHELL" == "bash" ]]; then \
			echo "Setting up bash completion for all users..."; \
			sudo mkdir -p /etc/bash_completion.d; \
			sudo bash -c "$$INSTALLED_BIN completion bash > /etc/bash_completion.d/togo"; \
			echo "✅ Bash completion installed to /etc/bash_completion.d/togo"; \
			ROOT_BASHRC="/root/.bashrc"; \
			if [ ! -f "$$ROOT_BASHRC" ]; then \
				sudo touch "$$ROOT_BASHRC"; \
			fi; \
			if ! sudo grep -q "eval \"\$$($$INSTALLED_BIN completion bash)\"" "$$ROOT_BASHRC" 2>/dev/null; then \
				echo "eval \"\$$($$INSTALLED_BIN completion bash)\"" | sudo tee -a "$$ROOT_BASHRC" >/dev/null; \
			fi; \
			USER_BASHRC="$$HOME/.bashrc"; \
			if [ ! -f "$$USER_BASHRC" ]; then \
				touch "$$USER_BASHRC"; \
			fi; \
			if ! grep -q "eval \"\$$($$INSTALLED_BIN completion bash)\"" "$$USER_BASHRC" 2>/dev/null; then \
				echo "eval \"\$$($$INSTALLED_BIN completion bash)\"" >>"$$USER_BASHRC"; \
			fi; \
			echo "✅ Bash completion configured for both root and current user"; \
		elif [[ "$$CURRENT_SHELL" == "fish" ]]; then \
			echo "Setting up fish completion for all users..."; \
			ROOT_FISH_COMPLETIONS="/root/.config/fish/completions"; \
			USER_FISH_COMPLETIONS="$$HOME/.config/fish/completions"; \
			if [[ "$$OS" == "Linux" ]]; then \
				echo "Installing fish completion for root user..."; \
				sudo mkdir -p "$$ROOT_FISH_COMPLETIONS"; \
				echo "togo completion fish | source" | sudo tee "$$ROOT_FISH_COMPLETIONS/togo.fish" >/dev/null; \
				echo "✅ Root fish completion installed to $$ROOT_FISH_COMPLETIONS/togo.fish"; \
			fi; \
			echo "Installing fish completion for current user..."; \
			mkdir -p "$$USER_FISH_COMPLETIONS"; \
			echo "togo completion fish | source" >"$$USER_FISH_COMPLETIONS/togo.fish"; \
			echo "✅ User fish completion installed to $$USER_FISH_COMPLETIONS/togo.fish"; \
		else \
			echo "⚠️  Unsupported shell for autocompletion. Please see README.md for manual setup."; \
			echo "Run 'togo completion -h' for more information."; \
		fi; \
	else \
		if [[ "$$CURRENT_SHELL" != "zsh" && "$$CURRENT_SHELL" != "bash" && "$$CURRENT_SHELL" != "fish" ]]; then \
			echo "⚠️  Unsupported shell for autocompletion. Please see README.md for manual setup."; \
			echo "Run 'togo completion -h' for more information."; \
			exit 0; \
		fi; \
		if [[ "$$CURRENT_SHELL" == "zsh" ]]; then \
			echo "Setting up zsh completion for current user..."; \
			ZSHRC_PATH="$$HOME/.zshrc"; \
			if [ ! -f "$$ZSHRC_PATH" ]; then \
				touch "$$ZSHRC_PATH" || { \
					echo "❌ Failed to create ~/.zshrc"; \
					echo "Please add this line manually to your shell config:"; \
					echo "eval \"\$$($$INSTALLED_BIN completion zsh)\""; \
					exit 1; \
				}; \
			fi; \
			if echo "eval \"\$$($$INSTALLED_BIN completion zsh)\"" >>"$$ZSHRC_PATH"; then \
				echo "✅ Added zsh completion to ~/.zshrc"; \
				source "$$ZSHRC_PATH" 2>/dev/null || echo "Please restart your shell or run: source ~/.zshrc"; \
			else \
				echo "❌ Failed to add completion to ~/.zshrc"; \
				echo "Run 'togo completion -h' for manual setup instructions."; \
			fi; \
		elif [[ "$$CURRENT_SHELL" == "bash" ]]; then \
			echo "Setting up bash completion for current user..."; \
			BASHRC_PATH="$$HOME/.bashrc"; \
			if [ ! -f "$$BASHRC_PATH" ]; then \
				touch "$$BASHRC_PATH" || { \
					echo "❌ Failed to create ~/.bashrc"; \
					echo "Please add this line manually to your shell config:"; \
					echo "eval \"\$$(togo completion bash)\""; \
					exit 1; \
				}; \
			fi; \
			if echo "eval \"\$$(togo completion bash)\"" >>"$$BASHRC_PATH"; then \
				echo "✅ Added bash completion to ~/.bashrc"; \
				source "$$BASHRC_PATH" 2>/dev/null || echo "Please restart your shell or run: source ~/.bashrc"; \
			else \
				echo "❌ Failed to add completion to ~/.bashrc"; \
				echo "Run 'togo completion -h' for manual setup instructions."; \
			fi; \
		elif [[ "$$CURRENT_SHELL" == "fish" ]]; then \
			echo "Setting up fish completion for current user..."; \
			USER_FISH_COMPLETIONS="$$HOME/.config/fish/completions"; \
			mkdir -p "$$USER_FISH_COMPLETIONS"; \
			if echo "togo completion fish | source" >"$$USER_FISH_COMPLETIONS/togo.fish"; then \
				echo "✅ Added fish completion to $$USER_FISH_COMPLETIONS/togo.fish"; \
				echo "Please restart your shell or run: source $$USER_FISH_COMPLETIONS/togo.fish"; \
			else \
				echo "❌ Failed to add fish completion"; \
				echo "Please add this line manually to your fish completions:"; \
				echo "togo completion fish | source"; \
			fi; \
		fi; \
	fi

clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)

fmt:
	@echo "Fixing import issues with goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@goimports -w .
	@echo "Fixing code formatting with gofmt..."
	@gofmt -s -w .

format: fmt

.DEFAULT_GOAL := install

