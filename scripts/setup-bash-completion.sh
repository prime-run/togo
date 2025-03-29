#!/bin/bash

echo "Setting up togo bash completion..."

if ! command -v bash-completion >/dev/null 2>&1; then
    echo "bash-completion is not installed."
    echo "Please install it using your package manager."
    echo "For example:"
    echo "  Ubuntu/Debian: sudo apt-get install bash-completion"
    echo "  Fedora: sudo dnf install bash-completion"
    echo "  Arch Linux: sudo pacman -S bash-completion"
    echo "  macOS: brew install bash-completion"
    exit 1
fi

TOGO_PATH=$(which togo 2>/dev/null)
if [ -z "$TOGO_PATH" ]; then
    TOGO_PATH="$PWD/togo"
    if [ ! -f "$TOGO_PATH" ]; then
        echo "Error: togo binary not found."
        exit 1
    fi
fi

if [[ "$OSTYPE" == "darwin"* ]]; then
    COMPLETION_DIR=$(brew --prefix 2>/dev/null)/etc/bash_completion.d
    if [ -z "$COMPLETION_DIR" ] || [ ! -d "$COMPLETION_DIR" ]; then
        echo "Error: Homebrew bash-completion directory not found."
        echo "Please install homebrew and bash-completion first."
        exit 1
    fi
    
    echo "Installing completion for macOS..."
    echo "Generating completion script for bash..."
    "$TOGO_PATH" completion bash > "$COMPLETION_DIR/togo"
    
    echo "Bash completion installed successfully."
    echo "Please restart your shell or source the completion file:"
    echo "  source '$COMPLETION_DIR/togo'"
else
    SYSTEM_COMPLETION_DIR="/etc/bash_completion.d"
    USER_COMPLETION_DIR="$HOME/.local/share/bash-completion/completions"
    
    mkdir -p "$USER_COMPLETION_DIR"
    
    echo "Do you want to install completion for all users (requires sudo) or just for you?"
    echo "1) All users (sudo required)"
    echo "2) Just for me (no sudo required)"
    read -p "Enter your choice (1 or 2): " choice
    
    if [ "$choice" = "1" ]; then
        echo "Installing system-wide completion (sudo required)..."
        "$TOGO_PATH" completion bash | sudo tee "$SYSTEM_COMPLETION_DIR/togo" > /dev/null
        
        if [ $? -eq 0 ]; then
            echo "Bash completion installed successfully for all users."
            echo "Please restart your shell or source the completion file:"
            echo "  source '$SYSTEM_COMPLETION_DIR/togo'"
        else
            echo "Failed to install system-wide completion. Trying user-local installation..."
            "$TOGO_PATH" completion bash > "$USER_COMPLETION_DIR/togo"
            
            echo "Adding source command to your .bashrc..."
            echo "source $USER_COMPLETION_DIR/togo" >> "$HOME/.bashrc"
            
            echo "Bash completion installed for current user only."
            echo "Please restart your shell or run:"
            echo "  source '$USER_COMPLETION_DIR/togo'"
        fi
    else
        echo "Installing user-local completion..."
        "$TOGO_PATH" completion bash > "$USER_COMPLETION_DIR/togo"
        
        if ! grep -q "source.*$USER_COMPLETION_DIR/togo" "$HOME/.bashrc"; then
            echo "Adding source command to your .bashrc..."
            echo "source $USER_COMPLETION_DIR/togo" >> "$HOME/.bashrc"
        fi
        
        echo "Bash completion installed for current user only."
        echo "Please restart your shell or run:"
        echo "  source '$USER_COMPLETION_DIR/togo'"
    fi
fi

echo "Setup complete!" 