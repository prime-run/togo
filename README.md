<!-- <div align="center"> -->
  <!-- <a href="https://github.com/prime-run/togo"> -->
  <!--   <img src="https://github.com/SenaThenu/readme-forge/blob/main/src/assets/logo.svg?raw=true" alt="Repo Logo" height="100"> -->
  <!-- </a> -->
<!-- </div> -->

# Togo

togo is a command line task/todo management utility designed to be simple, fast, and easy to use.

<div align="center">
<p align="center">  
<img src="https://github.com/user-attachments/assets/7c013d6a-4c3e-48a2-88f0-6ed65bc61ff5"
  alt="main-togo-screen-shot"
  width="633" height="353">
</p>
</div>

- [Features](#features)
- [Installation](#installation)
  - [Arch Linux](#arch-linux)
  - [Pre-built Binaries (Recommended)](#pre-built-binaries-recommended)
  - [Go CLI](#go-cli)
  - [Build from Source](#build-from-source)
- [Usage](#usage)
  - [1. Interactive Mode](#1-interactive-mode)
  - [2. Command-Line Operations](#2-command-line-operations)
    - [a) Direct selection by partial name](#a-direct-selection-by-partial-name)
    - [b) Interactive selection list](#b-interactive-selection-list)
    - [c) Shell completion integration](#c-shell-completion-integration)
  - [Available Commands](#available-commands)
- [Features in Depth](#features-in-depth)
  - [Shell Completion](#shell-completion)

---

## Features

- **CLI & TUI Interfaces**: A CLI for single-task operations and an interactive TUI for bulk operations and list view.
- **Vim Keybinds**: A handful of Vim actions are built into the TUI (more will be added soon).
- **Fuzzy Search & Filtering**: Find tasks quickly with partial name matching.
- **Tab Completion**: Shell completion with fuzzy matching support built into the completion script.
- **Project/Global Sources**: Load/save from the closest `.togo` file in your project tree, with a global fallback. Force with `--source project|global`.

## Installation

### Arch Linux

`togo` is available on the [Arch User Repository (AUR)](https://aur.archlinux.org/packages/togo-bin)

```bash
paru -S togo-bin
```

### Pre-built Binaries (Recommended)

Download the latest pre-built binaries from the [Releases](https://github.com/prime-run/togo/releases) page.

**Linux (x86_64)**

```bash
wget https://github.com/prime-run/togo/releases/download/v1.0.5/togo_1.0.5_linux_amd64.tar.gz
mkdir -p ~/.local/bin
tar -xzf togo_*.tar.gz -C ~/.local/bin/togo
```

**macOS (Apple Silicon arm64)**

```bash
wget https://github.com/prime-run/togo/releases/download/v1.0.5/togo_1.0.5_darwin_arm64.tar.gz
mkdir -p ~/.local/bin
tar -xzf togo_*.tar.gz -C ~/.local/bin/togo
```

> [!NOTE]
> Don't forget to set the `PATH` environment variable.

```bash
export PATH="$HOME/.local/bin:$PATH"
```

### Go CLI

```bash
go install github.com/prime-run/togo@latest
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/prime-run/togo.git
cd togo
make
# And follow the prompts
```

All `make` installation methods should set up shell completion out of the box.
In case your shell isn't detected, you can run `togo completion --help`.

## Usage

To add a task:

```bash
togo add "Task description"
# or without quotes, it's treated as stdin!
togo add Call the client about project scope
```

### Storage: project vs global (new)

Togo can load and save your tasks from two sources:

- "project" (default): uses the closest `.togo` file found by walking up from the current directory. If none is found, it falls back to the global file.
- "global": always uses the global file in your user config directory.

Control the source with the persistent flag on any command:

```bash
# project (default)
togo list
togo add "Write docs" -s project

# force global
togo list --source global
```

Initialize project-local storage in your repo:

```bash
togo init   # creates ./\.togo with an empty JSON list
```

In the TUI, the header shows the active source as `source: project` or `source: global`.

### Managing Your Tasks

Togo provides two primary modes of operation:

#### 1. Interactive Mode (TUI)

The TUI lets you work with your todos visually and allows for bulk actions:

```bash
togo
# or
togo list           # Active todos only
togo list --all     # All todos
togo list --archived # Archived todos only
```

#### 2. Command-Line Operations

Togo offers flexible command syntax with three usage patterns:

> All commands accept the persistent `--source|-s` flag: `project` (default) or `global`.

##### a) Direct selection by partial name

```bash
togo toggle meeting
# force using global source
togo toggle meeting --source global
```

If only one task contains "meeting," it executes immediately—no selection needed. If multiple tasks match (e.g., "team meeting" and "client meeting"), Togo automatically opens the selection list so you can choose the one you meant.

##### b) Interactive selection list

```bash
togo toggle
# against project-local source explicitly
togo toggle -s project
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

The shell will show only tasks containing "me"—perfect for quick selection.

`--source` also supports completion to `project` and `global`.

> [!TIP]
> This really speeds up task management since fuzzy matching is supported by the completion script.
> For example, to mark the task `Auto ***snapshot*** /boot ...` as `completed`,
> just running `togo toggle snap` would toggle it. If `snap` matches more than one task,
> you'll be prompted to select from the matches.

### Available Commands

- `togo add "Task description"` - Add a new task
- `togo toggle [task]` - Toggle completion status
- `togo archive [task]` - Archive a completed task
- `togo unarchive [task]` - Restore an archived task
- `togo delete [task]` - Remove a task permanently
- `togo list [flags]` - View tasks (`--all`, `--archived`)
- `togo init` - Create an empty `.togo` file in the current directory (enable project-local storage)

Notes:

- All commands accept `--source|-s {project|global}` to control where tasks are read/written.

### Features in Depth

### Shell Completion

Enabling shell completion allows for tab-completion of commands and task names, improving efficiency.

#### Zsh

```bash
# 1. Create completion directory
mkdir -p ~/.zsh/completion
echo "fpath=(~/.zsh/completion \$fpath)" >> ~/.zshrc

# 2. Enable completions
echo "autoload -Uz compinit && compinit" >> ~/.zshrc

# 3. Apply Togo completion
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

> Storage locations:
>
> - Project: nearest `./.togo` file (JSON).
> - Global: `$XDG_CONFIG_HOME/togo/todos.json` (or `$HOME/.config/togo/todos.json`).

## Built With 🔧

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Cobra](https://img.shields.io/badge/Cobra-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://github.com/spf13/cobra)
[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-FF75B7?style=for-the-badge&logo=go&logoColor=white)](https://github.com/charmbracelet/bubbletea)
[![Lip Gloss](https://img.shields.io/badge/Lip%20Gloss-FFABE7?style=for-the-badge&logo=go&logoColor=white)](https://github.com/charmbracelet/lipgloss)

## License

This project is licensed under MIT - see the [LICENSE](LICENSE) file for details.
