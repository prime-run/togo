<div align="center">
  <a href="https://github.com/prime-run/togo">
    <img src="https://github.com/SenaThenu/readme-forge/blob/main/src/assets/logo.svg?raw=true" alt="Repo Logo" height="100">
  </a>
</div>

<h1 align="center">ToGo</h1>

<div align="center"> A fast and simple terminal-based to-do manager. </div>

<div align="center">
<p align="center">  
<img src="https://github.com/user-attachments/assets/7c013d6a-4c3e-48a2-88f0-6ed65bc61ff5"
  alt="main-togo-screen-shot"
  width="633" height="353">
</p>
</div>

<div align="center">
  <img src="https://img.shields.io/badge/license-MIT-blue.svg?labelColor=003694&color=ffffff" alt="License">
  <img src="https://img.shields.io/github/contributors/prime-run/togo?labelColor=003694&color=ffffff" alt="GitHub contributors" >
  <img src="https://img.shields.io/github/stars/prime-run/togo.svg?labelColor=003694&color=ffffff" alt="Stars">
  <img src="https://img.shields.io/github/forks/prime-run/togo.svg?labelColor=003694&color=ffffff" alt="Forks">
  <img src="https://img.shields.io/github/issues/prime-run/togo.svg?labelColor=003694&color=ffffff" alt="Issues">
</div>

- [Features](#features)
- [Installation](#installation)
  - [ArchLinux](#archlinux)
  - [pre-built binaries (recommended)](#pre-built-binaries-recommended)
  - [Linux and Mac (x86_64)](#linux-and-mac-x86_64)
  - [Go Cli](#go-cli)
  - [make](#make)
- [Usage](#usage)
  - [Managing Your Tasks](#managing-your-tasks)
    - [1. Interactive Mode](#1-interactive-mode)
    - [2. Command-Line Operations](#2-command-line-operations)
      - [a) Direct selection by partial name](#a-direct-selection-by-partial-name)
      - [b) Interactive selection list](#b-interactive-selection-list)
      - [c) Shell completion integration](#c-shell-completion-integration)
  - [Available Commands](#available-commands)
  - [Additional Options](#additional-options)
- [Features In Depth](#features-in-depth)
  - [Shell Completion](#shell-completion)

---

## Features

- **CLI & TUI Interfaces**: cli for single task operations and an interactive TUI for bulk operations and list view.
- **Vim Keybinds**: a handful of vim actions built-in the TUI (a few vim actions will be added soon).
- **Fuzzy Search & Filtering**: Find tasks quickly with partial name matching.
- **Tab Completion**: Shell Completion with fuzz support built-in the completion script.

## Installation

### ArchLinux

togo is pushed to arch [AUR](https://aur.archlinux.org/packages/togo) with **zero** dependencies.

```bash

paru -S togo

```

### pre-built binaries (recommended)

Download the latest pre-built binaries from the [Releases](https://github.com/prime-run/togo/releases) page.

Linux(x86_64)

```bash
wget https://github.com/prime-run/togo/releases/download/v1.0.5/togo_1.0.5_linux_amd64.tar.gz
mkdir -p ~/.local/bin
tar -xzf togo_*.tar.gz -C ~/.local/bin/togo
```

macOS (Apple Silicon arm64):

```bash
wget https://github.com/prime-run/togo/releases/download/v1.0.5/togo_1.0.5_darwin_arm64.tar.gz
mkdir -p ~/.local/bin
tar -xzf togo_*.tar.gz -C ~/.local/bin/togo
```

> [!NOTE]
> Don't forget to set PATH environment variable.

```bash
export PATH="$HOME/.local/bin:$PATH"
```

### Go Cli

```bash
go install github.com/prime-run/togo@latest
```

> Make sure `$GOPATH/bin` is in your PATH to access the installed binary.

### build from source

```bash
# Clone the repository
git clone https://github.com/prime-run/togo.git
cd togo
make
# And follow the prompts
```

All Make installation methods should setup shell completion out of the box.
In case your shell didn't get detected, you can run `togo completion --help`

## Usage

To add a task:

```bash
togo add "Task description"
# or without quotes its stdin!
togo add Call the client about project scope
```

### Managing Your Tasks

Togo provides two primary modes of operation:

#### 1. Interactive Mode

The `TUI` to work with your todos visually and allows for bulk actions:

```bash
togo
# or
togo list           # Active todos only
togo list --all     # All todos
togo list --archived # Archived todos only
```

#### 2. Command-Line Operations

Togo offers flexible command syntax with three usage patterns:

##### a) Direct selection by partial name

```bash
togo toggle meeting
```

If only one task contains "meeting", it executes immediately - no selection needed. If multiple tasks match (e.g., "team meeting" and "client meeting"), Togo automatically opens the selection list so you can choose which one you meant.

##### b) Interactive selection list

```bash
togo toggle
```

Opens a selection list where you can choose from available tasks:

<p align="center">
<img src="https://github.com/user-attachments/assets/aa0c3005-af4c-4f2e-bf4c-df6681050ad6"
  alt="main-togo-screen-shot"
  width="738">
  
</p>

As you type, Togo searches through your tasks and filters the results.

##### c) Shell completion integration

If you've configured shell completion (see below), you can use tab-completion:

```bash
togo toggle [TAB]
```

Your shell will present available tasks. Type a few letters to filter by name:

```bash
togo toggle me[TAB]
```

Shell will show only tasks containing "me" - perfect for quick selection.

> [!TIP]
> This really speeds task managment up since the `fuzz` is supported by the completion script.
> e.g.take `Auto ***snapshot*** /boot ...` as task that needs to be set as `completed`,
> Just running `togo toggle snap` would toggle it! and if `snap` matches more than one task,
> you'd be prompted to select from matches.

### Available Commands

- `togo add "Task description"` - Add a new task
- `togo toggle [task]` - Toggle completion status
- `togo archive [task]` - Archive a completed task
- `togo unarchive [task]` - Restore an archived task
- `togo delete [task]` - Remove a task permanently
- `togo list [flags]` - View tasks (--all, --archived)

### Features In Depth

### Shell Completion

Enabling shell completion allows for tab-completion of commands and task names, improving efficiency.

#### Zsh

```bash
# 1. Create completion directory
mkdir -p ~/.zsh/completion
echo "fpath=(~/.zsh/completion \$fpath)" >> ~/.zshrc

# 2. Enable completions
echo "autoload -Uz compinit && compinit" >> ~/.zshrc

# 3. apply Togo completion
togo completion zsh > ~/.zsh/completion/_togo
source ~/.zshrc
```

#### Bash

```bash

# 1. Ensure completion is sourced
echo "[ -r /usr/share/bash-completion/bash_completion ] && . /usr/share/bash-completion/bash_completion" >> ~/.bashrc
source ~/.bashrc

# 2. Install Togo completion
togo completion bash > ~/.bash_completion
source ~/.bash_completion
```

#### Fish

```bash
mkdir -p ~/.config/fish/completions
togo completion fish > ~/.config/fish/completions/togo.fish
```

> All the taks are stored in a JSON file under at `~/.togo/todos.json`.

## Built With 🔧

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Cobra](https://img.shields.io/badge/Cobra-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://github.com/spf13/cobra)
[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-FF75B7?style=for-the-badge&logo=go&logoColor=white)](https://github.com/charmbracelet/bubbletea)
[![Lip Gloss](https://img.shields.io/badge/Lip%20Gloss-FFABE7?style=for-the-badge&logo=go&logoColor=white)](https://github.com/charmbracelet/lipgloss)

## License

This project is licensed under MIT - see the [LICENSE](LICENSE) file for details.
