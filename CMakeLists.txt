cmake_minimum_required(VERSION 3.10)
project(togo)

# Variables for installation paths
set(COMPLETION_ZSH_DIR "$ENV{HOME}/.zsh/completion" CACHE PATH "ZSH completion directory")
set(COMPLETION_BASH_FILE "$ENV{HOME}/.bash_completion" CACHE PATH "Bash completion file")
set(COMPLETION_FISH_DIR "$ENV{HOME}/.config/fish/completions" CACHE PATH "Fish completion directory")

# Detect shell
find_program(SHELL_ZSH zsh)
find_program(SHELL_BASH bash)
find_program(SHELL_FISH fish)

# Check if togo binary exists
if(NOT EXISTS "${CMAKE_CURRENT_SOURCE_DIR}/togo")
    message(STATUS "No togo binary found. Skipping completion setup.")
    return()
endif()

# Setup ZSH completion
add_custom_target(setup-zsh-completion
    COMMAND ${CMAKE_COMMAND} -E make_directory ${COMPLETION_ZSH_DIR}
    COMMAND ${CMAKE_CURRENT_SOURCE_DIR}/togo completion zsh > ${COMPLETION_ZSH_DIR}/_togo || (echo "Error generating ZSH completion" && exit 0)
    COMMAND grep -q "fpath=(${COMPLETION_ZSH_DIR}" $ENV{HOME}/.zshrc || echo "fpath=(${COMPLETION_ZSH_DIR} \$$fpath)" >> $ENV{HOME}/.zshrc || true
    COMMAND grep -q "compinit" $ENV{HOME}/.zshrc || echo "autoload -Uz compinit && compinit" >> $ENV{HOME}/.zshrc || true
    COMMAND echo "ZSH completion setup successful! Please restart your shell or run 'source ~/.zshrc'"
    COMMENT "Setting up ZSH completion"
)

# Setup Bash completion
add_custom_target(setup-bash-completion
    COMMAND ${CMAKE_COMMAND} -E make_directory $(dirname ${COMPLETION_BASH_FILE})
    COMMAND ${CMAKE_CURRENT_SOURCE_DIR}/togo completion bash > ${COMPLETION_BASH_FILE} || (echo "Error generating Bash completion" && exit 0)
    COMMAND grep -q "bash_completion" $ENV{HOME}/.bashrc || echo "[ -r ${COMPLETION_BASH_FILE} ] && . ${COMPLETION_BASH_FILE}" >> $ENV{HOME}/.bashrc || true
    COMMAND echo "Bash completion setup successful! Please restart your shell or run 'source ~/.bashrc'"
    COMMENT "Setting up Bash completion"
)

# Setup Fish completion
add_custom_target(setup-fish-completion
    COMMAND ${CMAKE_COMMAND} -E make_directory ${COMPLETION_FISH_DIR}
    COMMAND ${CMAKE_CURRENT_SOURCE_DIR}/togo completion fish > ${COMPLETION_FISH_DIR}/togo.fish || (echo "Error generating Fish completion" && exit 0)
    COMMAND echo "Fish completion setup successful!"
    COMMENT "Setting up Fish completion"
)

# Detect default shell and set the default target
add_custom_target(setup-completion)

if(SHELL_ZSH)
    add_dependencies(setup-completion setup-zsh-completion)
    message(STATUS "Detected ZSH as default shell")
elseif(SHELL_BASH)
    add_dependencies(setup-completion setup-bash-completion)
    message(STATUS "Detected Bash as default shell")
elseif(SHELL_FISH)
    add_dependencies(setup-completion setup-fish-completion)
    message(STATUS "Detected Fish as default shell")
else()
    message(STATUS "Could not detect shell type. No completion will be set up automatically.")
endif()

# Add a custom target to setup for all shells
add_custom_target(setup-all-completions)
add_dependencies(setup-all-completions setup-zsh-completion setup-bash-completion setup-fish-completion)

# Instructions for the user
message(STATUS "----------------------------------------")
message(STATUS "Togo shell completion setup")
message(STATUS "----------------------------------------") 