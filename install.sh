#!/bin/bash

# Build the project using Makefile
echo "Building project..."
if ! make build; then
    echo "Build failed. Exiting."
    exit 1
fi

# Installation type selection
echo "Select installation type:"
echo "1) Local user install (~/.local/bin)"
echo "2) System-wide install (requires sudo) (/usr/local/bin)"
read -p "Enter choice [1-2]: " install_choice

# Validate input
if [[ ! "$install_choice" =~ ^[12]$ ]]; then
    echo "Invalid selection. Exiting."
    exit 1
fi

# Install paths
BINARY_NAME="togo"
BUILD_PATH="./bin/$BINARY_NAME"

if [ ! -f "$BUILD_PATH" ]; then
    echo "Error: Built binary not found at $BUILD_PATH"
    exit 1
fi

# Determine OS
case "$(uname -s)" in
    Linux*)     OS=Linux;;
    Darwin*)    OS=macOS;;
    *)          OS="UNKNOWN";;
esac

if [[ "$OS" == "UNKNOWN" ]]; then
    echo "Unsupported OS. Proceeding with basic installation only."
fi

# Handle installation
case $install_choice in
    1)
        TARGET_DIR="$HOME/.local/bin"
        mkdir -p "$TARGET_DIR"
        if cp "$BUILD_PATH" "$TARGET_DIR/"; then
            echo "Installed to $TARGET_DIR/$BINARY_NAME"
        else
            echo "Local installation failed"
            exit 1
        fi
        ;;
    2)
        TARGET_DIR="/usr/local/bin"
        echo "Installing system-wide (requires sudo)..."
        if sudo cp "$BUILD_PATH" "$TARGET_DIR/"; then
            echo "Installed to $TARGET_DIR/$BINARY_NAME"
        else
            echo "System-wide installation failed"
            exit 1
        fi
        ;;
esac

# Detect and display current shell
CURRENT_SHELL=$(basename "$SHELL")
echo "Detected shell: $CURRENT_SHELL"

INSTALLED_BIN="${TARGET_DIR}/${BINARY_NAME}"

# Handle system-wide and local installations differently
if [[ "$install_choice" == "2" ]]; then
    # System-wide completion setup
    if [[ "$CURRENT_SHELL" == "zsh" ]]; then
        if [[ "$OS" == "Linux" ]]; then
            echo "Setting up system-wide zsh completion..."
            echo "Running: sudo $INSTALLED_BIN completion zsh > \"\${fpath[1]}/_togo\""
            sudo bash -c "$INSTALLED_BIN completion zsh > \$(echo \${fpath[1]}/_togo)"
        else
            echo "For macOS system-wide zsh completion, follow these instructions:"
            $INSTALLED_BIN completion zsh --help
        fi
    elif [[ "$CURRENT_SHELL" == "bash" ]]; then
        if [[ "$OS" == "Linux" ]]; then
            echo "Setting up system-wide bash completion..."
            sudo mkdir -p /etc/bash_completion.d
            sudo bash -c "$INSTALLED_BIN completion bash > /etc/bash_completion.d/togo"
            echo "Bash completion installed to /etc/bash_completion.d/togo"
        else
            echo "For macOS system-wide bash completion, follow these instructions:"
            $INSTALLED_BIN completion bash --help
        fi
    elif [[ "$CURRENT_SHELL" == "fish" ]]; then
        echo "For fish system-wide completion, follow these instructions:"
        $INSTALLED_BIN completion fish --help
    else
        echo "Unsupported shell for autocompletion. Please use built-in completion or see README.md"
        echo "Run 'togo completion -h' for more information."
    fi
else
    # Local user completion setup
    if [[ "$CURRENT_SHELL" != "zsh" && "$CURRENT_SHELL" != "bash" && "$CURRENT_SHELL" != "fish" ]]; then
        echo "Unsupported shell for autocompletion. Please use built-in completion or see README.md"
        echo "Run 'togo completion -h' for more information."
        exit 0
    fi

    # Shell-specific completion setup
    if [[ "$CURRENT_SHELL" == "zsh" ]]; then
        echo "Configuring zsh autocompletion..."
        ZSHRC_PATH="$HOME/.zshrc"
        if [ ! -f "$ZSHRC_PATH" ]; then
            touch "$ZSHRC_PATH" || { 
                echo "Failed to create ~/.zshrc"; 
                echo "Add this line manually to your shell config:";
                echo "eval \"\$($INSTALLED_BIN completion zsh)\"";
                exit 1; 
            }
        fi

        if ! grep -q "eval \"\$(${INSTALLED_BIN} completion zsh)\"" "$ZSHRC_PATH"; then
            if echo "eval \"\$(${INSTALLED_BIN} completion zsh)\"" >> "$ZSHRC_PATH"; then
                echo "Added zsh completion to ~/.zshrc"
                source "$ZSHRC_PATH" 2>/dev/null || echo "Please restart your shell or run: source ~/.zshrc"
            else
                echo "Failed to add completion to ~/.zshrc"
                echo "Run 'togo completion -h' for manual setup instructions."
            fi
        else
            echo "Zsh completion already configured in ~/.zshrc"
        fi
    elif [[ "$CURRENT_SHELL" == "bash" ]]; then
        echo "Configuring bash autocompletion..."
        BASHRC_PATH="$HOME/.bashrc"
        if [ ! -f "$BASHRC_PATH" ]; then
            touch "$BASHRC_PATH" || { 
                echo "Failed to create ~/.bashrc"; 
                echo "Add this line manually to your shell config:";
                echo "eval \"\$(togo completion bash)\"";
                exit 1; 
            }
        fi

        if ! grep -q "eval \"\$(togo completion bash)\"" "$BASHRC_PATH"; then
            if echo "eval \"\$(togo completion bash)\"" >> "$BASHRC_PATH"; then
                echo "Added bash completion to ~/.bashrc"
                source "$BASHRC_PATH" 2>/dev/null || echo "Please restart your shell or run: source ~/.bashrc"
            else
                echo "Failed to add completion to ~/.bashrc"
                echo "Run 'togo completion -h' for manual setup instructions."
            fi
        else
            echo "Bash completion already configured in ~/.bashrc"
        fi
    elif [[ "$CURRENT_SHELL" == "fish" ]]; then
        echo "Configuring fish autocompletion..."
        FISH_CONFIG_PATH="$HOME/.config/fish/config.fish"
        mkdir -p "$(dirname "$FISH_CONFIG_PATH")"
        
        if [ ! -f "$FISH_CONFIG_PATH" ]; then
            touch "$FISH_CONFIG_PATH" || { 
                echo "Failed to create fish config"; 
                echo "Add this line manually to your shell config:";
                echo "togo completion fish | source";
                exit 1; 
            }
        fi
        
        if ! grep -q "togo completion fish | source" "$FISH_CONFIG_PATH" 2>/dev/null; then
            if echo "togo completion fish | source" >> "$FISH_CONFIG_PATH"; then
                echo "Added fish completion to ~/.config/fish/config.fish"
                echo "Please restart your shell or run: source ~/.config/fish/config.fish"
            else
                echo "Failed to add completion to ~/.config/fish/config.fish"
                echo "Run 'togo completion -h' for manual setup instructions."
            fi
        else
            echo "Fish completion already configured in ~/.config/fish/config.fish"
        fi
    fi
fi
