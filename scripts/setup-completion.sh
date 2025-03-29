#!/bin/bash

# Main script to set up shell completion for togo

echo "togo Shell Completion Setup"
echo "=========================="
echo

# Make all scripts executable
chmod +x $(dirname "$0")/setup-*-completion.sh 2>/dev/null

# Determine current shell
CURRENT_SHELL=$(basename "$SHELL")

# Ask user which shell to setup
echo "Please select the shell to set up completion for:"
echo "1) bash"
echo "2) zsh"
echo "3) fish"
echo "4) powershell"
echo "5) Detect current shell ($CURRENT_SHELL)"
read -p "Enter choice [1-5]: " shell_choice

SCRIPT_TO_RUN=""

case $shell_choice in
    1)
        SCRIPT_TO_RUN="setup-bash-completion.sh"
        ;;
    2)
        SCRIPT_TO_RUN="setup-zsh-completion.sh"
        ;;
    3)
        SCRIPT_TO_RUN="setup-fish-completion.sh"
        ;;
    4)
        echo "To run PowerShell completion setup, you must be in PowerShell."
        echo "Please open PowerShell and run the following command:"
        echo "  ./scripts/setup-powershell-completion.ps1"
        exit 0
        ;;
    5|"")
        # Detect current shell
        case $CURRENT_SHELL in
            bash)
                SCRIPT_TO_RUN="setup-bash-completion.sh"
                ;;
            zsh)
                SCRIPT_TO_RUN="setup-zsh-completion.sh"
                ;;
            fish)
                SCRIPT_TO_RUN="setup-fish-completion.sh"
                ;;
            pwsh|powershell.exe)
                echo "You appear to be running PowerShell."
                echo "To run PowerShell completion setup, you must be in PowerShell."
                echo "Please open PowerShell and run the following command:"
                echo "  ./scripts/setup-powershell-completion.ps1"
                exit 0
                ;;
            *)
                echo "Your current shell ($CURRENT_SHELL) is not supported."
                echo "Please select a specific shell (options 1-4)."
                exit 1
                ;;
        esac
        ;;
    *)
        echo "Invalid choice. Please select a number between 1 and 5."
        exit 1
        ;;
esac

if [ -n "$SCRIPT_TO_RUN" ]; then
    SCRIPT_PATH="$(dirname "$0")/$SCRIPT_TO_RUN"
    
    if [ -f "$SCRIPT_PATH" ]; then
        echo "Running $SCRIPT_PATH..."
        "$SCRIPT_PATH"
    else
        echo "Error: Script not found at $SCRIPT_PATH"
        exit 1
    fi
fi 