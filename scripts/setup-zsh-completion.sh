#!/bin/bash

echo "Setting up togo zsh completion..."

TOGO_PATH=$(which togo 2>/dev/null)
if [ -z "$TOGO_PATH" ]; then
    TOGO_PATH="$PWD/togo"
    if [ ! -f "$TOGO_PATH" ]; then
        echo "Error: togo binary not found."
        exit 1
    fi
fi

if ! command -v zsh >/dev/null 2>&1; then
    echo "zsh is not installed."
    echo "Please install zsh using your package manager."
    exit 1
fi

if [ -f "$HOME/.zshrc" ]; then
    if ! grep -q "compinit" "$HOME/.zshrc"; then
        echo "Adding compinit to .zshrc..."
        echo "autoload -U compinit; compinit" >> "$HOME/.zshrc"
        echo "Added compinit to .zshrc"
    fi
else
    echo "Creating .zshrc with compinit..."
    echo "autoload -U compinit; compinit" > "$HOME/.zshrc"
fi

if [[ "$OSTYPE" == "darwin"* ]]; then
    if ! command -v brew >/dev/null 2>&1; then
        echo "Homebrew not found. Cannot determine zsh site-functions path."
        echo "Please install homebrew or use the manual installation method."
        
        USER_COMPLETION_DIR="$HOME/.zsh/completion"
        mkdir -p "$USER_COMPLETION_DIR"
        
        echo "Installing to user-local directory instead..."
        "$TOGO_PATH" completion zsh > "$USER_COMPLETION_DIR/_togo"
        
        if ! grep -q "$USER_COMPLETION_DIR" "$HOME/.zshrc"; then
            echo "Adding completion directory to .zshrc..."
            echo "fpath=($USER_COMPLETION_DIR \$fpath)" >> "$HOME/.zshrc"
        fi
        
        echo "Completion installed to $USER_COMPLETION_DIR/_togo"
        echo "Please restart your shell or source your .zshrc:"
        echo "  source $HOME/.zshrc"
    else
        COMPLETION_DIR=$(brew --prefix)/share/zsh/site-functions
        
        echo "Installing completion for macOS..."
        "$TOGO_PATH" completion zsh > "$COMPLETION_DIR/_togo"
        
        echo "zsh completion installed successfully."
        echo "Please restart your shell or run:"
        echo "  autoload -U compinit; compinit"
    fi
else
    SYSTEM_COMPLETION_DIR="/usr/share/zsh/site-functions"
    USER_COMPLETION_DIR="$HOME/.zsh/completion"
    
    mkdir -p "$USER_COMPLETION_DIR"
    
    if ! grep -q "$USER_COMPLETION_DIR" "$HOME/.zshrc"; then
        echo "Adding completion directory to .zshrc..."
        echo "fpath=($USER_COMPLETION_DIR \$fpath)" >> "$HOME/.zshrc"
    fi
    
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
        
        "$TOGO_PATH" completion zsh | sudo tee "$SYSTEM_COMPLETION_DIR/_togo" > /dev/null
        
        if [ $? -eq 0 ]; then
            echo "zsh completion installed successfully for all users."
            echo "Please restart your shell or run:"
            echo "  autoload -U compinit; compinit"
        else
            echo "Failed to install system-wide completion. Trying user-local installation..."
            "$TOGO_PATH" completion zsh > "$USER_COMPLETION_DIR/_togo"
            
            echo "zsh completion installed for current user only."
            echo "Please restart your shell or source your .zshrc:"
            echo "  source $HOME/.zshrc"
        fi
    else
        echo "Installing user-local completion..."
        "$TOGO_PATH" completion zsh > "$USER_COMPLETION_DIR/_togo"
        
        echo "zsh completion installed for current user only."
        echo "Please restart your shell or source your .zshrc:"
        echo "  source $HOME/.zshrc"
    fi
fi

echo "Setup complete!" 