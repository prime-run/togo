#!/bin/bash

echo "Setting up togo fish completion..."

TOGO_PATH=$(which togo 2>/dev/null)
if [ -z "$TOGO_PATH" ]; then
    TOGO_PATH="$PWD/togo"
    if [ ! -f "$TOGO_PATH" ]; then
        echo "Error: togo binary not found."
        exit 1
    fi
fi

if ! command -v fish >/dev/null 2>&1; then
    echo "fish is not installed."
    echo "Please install fish using your package manager."
    exit 1
fi

SYSTEM_COMPLETION_DIR="/usr/share/fish/vendor_completions.d"
USER_COMPLETION_DIR="$HOME/.config/fish/completions"

mkdir -p "$USER_COMPLETION_DIR"

echo "Do you want to install completion for all users (requires sudo) or just for you?"
echo "1) All users (sudo required)"
echo "2) Just for me (no sudo required)"
read -p "Enter your choice (1 or 2): " choice

if [ "$choice" = "1" ]; then
    echo "Installing system-wide completion (sudo required)..."
    
    if [ ! -d "$SYSTEM_COMPLETION_DIR" ]; then
        echo "System completion directory doesn't exist, creating it (requires sudo)..."
        sudo mkdir -p "$SYSTEM_COMPLETION_DIR"
    fi
    
    "$TOGO_PATH" completion fish | sudo tee "$SYSTEM_COMPLETION_DIR/togo.fish" > /dev/null
    
    if [ $? -eq 0 ]; then
        echo "fish completion installed successfully for all users."
        echo "Please restart your shell for the changes to take effect."
    else
        echo "Failed to install system-wide completion. Trying user-local installation..."
        "$TOGO_PATH" completion fish > "$USER_COMPLETION_DIR/togo.fish"
        
        echo "fish completion installed for current user only."
        echo "Please restart your shell for the changes to take effect."
    fi
else
    echo "Installing user-local completion..."
    "$TOGO_PATH" completion fish > "$USER_COMPLETION_DIR/togo.fish"
    
    echo "fish completion installed for current user only."
    echo "Please restart your shell for the changes to take effect."
fi

echo "To use completion in current shell without restarting, run:"
echo "  $TOGO_PATH completion fish | source"

echo "Setup complete!" 