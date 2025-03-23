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

# Custom target to build the togo binary if needed
add_custom_target(build-togo
    COMMAND go build -o togo
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    COMMENT "Building togo binary"
)

# Setup ZSH completion
add_custom_target(setup-zsh-completion
    COMMAND mkdir -p ${COMPLETION_ZSH_DIR}
    COMMAND ${CMAKE_CURRENT_SOURCE_DIR}/togo completion zsh > ${COMPLETION_ZSH_DIR}/_togo
    COMMAND echo "fpath=(${COMPLETION_ZSH_DIR} \$$fpath)" >> $ENV{HOME}/.zshrc
    COMMAND echo "autoload -Uz compinit && compinit" >> $ENV{HOME}/.zshrc
    COMMAND echo "ZSH completion setup successful! Please restart your shell or run 'source ~/.zshrc'"
    DEPENDS build-togo
    COMMENT "Setting up ZSH completion"
)

# Setup Bash completion
add_custom_target(setup-bash-completion
    COMMAND mkdir -p $(dirname ${COMPLETION_BASH_FILE})
    COMMAND ${CMAKE_CURRENT_SOURCE_DIR}/togo completion bash > ${COMPLETION_BASH_FILE}
    COMMAND echo "[ -r ${COMPLETION_BASH_FILE} ] && . ${COMPLETION_BASH_FILE}" >> $ENV{HOME}/.bashrc
    COMMAND echo "Bash completion setup successful! Please restart your shell or run 'source ~/.bashrc'"
    DEPENDS build-togo
    COMMENT "Setting up Bash completion"
)

# Setup Fish completion
add_custom_target(setup-fish-completion
    COMMAND mkdir -p ${COMPLETION_FISH_DIR}
    COMMAND ${CMAKE_CURRENT_SOURCE_DIR}/togo completion fish > ${COMPLETION_FISH_DIR}/togo.fish
    COMMAND echo "Fish completion setup successful!"
    DEPENDS build-togo
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
    message(WARNING "Could not detect shell type. Please run one of the specific setup targets.")
endif()

# Add a custom target to setup for all shells
add_custom_target(setup-all-completions)
add_dependencies(setup-all-completions setup-zsh-completion setup-bash-completion setup-fish-completion)

# Instructions for the user
message(STATUS "----------------------------------------")
message(STATUS "Usage:")
message(STATUS "  cmake -B build && cmake --build build --target setup-completion")
message(STATUS "")
message(STATUS "Available targets:")
message(STATUS "  setup-completion      - Setup completion for detected default shell")
message(STATUS "  setup-zsh-completion  - Setup completion for ZSH")
message(STATUS "  setup-bash-completion - Setup completion for Bash")
message(STATUS "  setup-fish-completion - Setup completion for Fish")
message(STATUS "  setup-all-completions - Setup completion for all shells")
message(STATUS "----------------------------------------") 