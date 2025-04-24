<div align="center">
  <a href="https://github.com/prime-run/togo">
    <img src="https://github.com/SenaThenu/readme-forge/blob/main/src/assets/logo.svg?raw=true" alt="Repo Logo" height="100">
  </a>
</div>

<h1 align="center">ToGo</h1>


<div align="center">	A blazingly fast, simple and beautiful terminal-based to-do manager </div>

<h2 align="center">.</h2>

<div align="center">
  <img src="https://img.shields.io/badge/license-MIT-blue.svg?labelColor=003694&color=ffffff" alt="License">
  <img src="https://img.shields.io/github/contributors/prime-run/togo?labelColor=003694&color=ffffff" alt="GitHub contributors" >
  <img src="https://img.shields.io/github/stars/prime-run/togo.svg?labelColor=003694&color=ffffff" alt="Stars">
  <img src="https://img.shields.io/github/forks/prime-run/togo.svg?labelColor=003694&color=ffffff" alt="Forks">
  <img src="https://img.shields.io/github/issues/prime-run/togo.svg?labelColor=003694&color=ffffff" alt="Issues">
</div>
<div align="center">

<p align="center">  
<img src="https://github.com/user-attachments/assets/20743925-cc5e-4695-afee-8dba4467718b"
  alt="main-togo-screen-shot"
  width="633" height="353">
</p>
  
</div>

<summary><strong>Table of Contents 📜</strong></summary>

- [Features ✨](#features-)
- [Installation 📥](#installation-)
  - [ArchLinux](#archlinux)
  - [pre-built binaries (recommended)](#pre-built-binaries-recommended)
  - [Linux and Mac (x86_64):](#linux-and-mac-x86_64)
  - [Go Cli](#go-cli)
  - [make](#make)
- [Usage 🛠️](#usage-)
  - [Managing Your Tasks](#managing-your-tasks)
    - [1. Interactive Mode](#1-interactive-mode)
    - [2. Command-Line Operations](#2-command-line-operations)
      - [a) Direct selection by partial name](#a-direct-selection-by-partial-name)
      - [b) Interactive selection list](#b-interactive-selection-list)
      - [c) Shell completion integration](#c-shell-completion-integration)
  - [Available Commands](#available-commands)
  - [Additional Options 📌](#additional-options-)
- [Features In Depth 🧠](#features-in-depth-)
  - [Shell Completion](#shell-completion)
    - [Zsh](#zsh)
    - [Bash](#bash)
    - [Fish](#fish)
    - [PowerShell](#powershell)
  - [Data Storage](#data-storage)

---

<p align="center">
<img src="https://github.com/user-attachments/assets/14aa2924-e310-4c46-a7e2-9af9effe89f0"
  alt="main-togo-screen-shot"
  width="738">
</p>

## Features ✨

- **Zero-friction capture**: Add ideas directly from your terminal without interrupting your flow
- **Beautiful terminal UI**: Interactive interface for managing todos when you're ready to organize
- **VIM keybinds**: HJKL motions support
- **Multiple management methods**: Use either interactive mode or command-line operations to manage
- **Flexible organization**: Toggle completion, archive finished tasks, delete what's no longer needed
- **Search/filtering**: Find tasks quickly in lists or through partial start-matching
- **Shell integration**: Tab completion for workflow integration

## Installation 📥

### ArchLinux

togo is pushed to arch [AUR](https://aur.archlinux.org/packages/togo) with **zero** dependencies.

```bash
yay -Sy togo
#or
paru -Sy togo
```

### pre-built binaries (recommended)

Download the latest pre-built binaries for your operating system from the [Releases](https://github.com/prime-run/togo/releases) page.
After downloading and extracting, ensure that ~/.local/bin is in your system's PATH environment variable. You can usually do this by adding the following line to your shell's configuration file (e.g., `.bashrc` , `.zshrc` ):

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Then, reload your shell configuration:

```bash
source ~/.bashrc  # For Bash
# or
source ~/.zshrc  # For Zsh
```

Now you should be able to run togo in your terminal.

### Linux and Mac (x86_64):

```bash
wget https://github.com/prime-run/togo/releases/download/v1.0.2/togo_1.0.2_linux_amd64.tar.gz
mkdir -p ~/.local/bin
tar -xzf togo_*.tar.gz -C ~/.local/bin/togo
```

macOS (Apple Silicon arm64):

```bash
wget https://github.com/prime-run/togo/releases/download/v1.0.2/togo_1.0.2_darwin_arm64.tar.gz
mkdir -p ~/.local/bin
tar -xzf togo_*.tar.gz -C ~/.local/bin/togo
```

### Go Cli

The simplest way to install Togo:

> [!CAUTION]
> go version > 1.24 is required

```bash
go install github.com/prime-run/togo@latest
```

Make sure `$GOPATH/bin` is in your PATH to access the installed binary.

### make

```bash
# Clone the repository
git clone https://github.com/prime-run/togo.git
cd togo
make
# And follow the prompts
```

All Make installation methods include automatic shell completion setup out of the box, so you can immediately use tab completion for commands and task names. In case your shell didn't get detected, you can run `togo completion --help`

## Usage 🛠️

Add your first task:

```bash
togo add should I use "s in my shell std inputs?"
# or dont even add quotes
togo add Call the client about project scope
```

### Managing Your Tasks

Togo offers two primary ways to manage your tasks:

#### 1. Interactive Mode

Open the interactive UI to work with your todos visually:

```bash
togo
# or
togo list           # Active todos only
togo list --all     # All todos
togo list --archived # Archived todos only
```

The interactive mode shows helpful keyboard shortcuts right in the interface.

<p align="center">
<img src="https://github.com/user-attachments/assets/e75cb61e-00f5-4c5b-ae44-66727521d2c4"
  alt="main-togo-screen-shot"
  width="686" height="289">
</p>

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
<img src="https://github.com/user-attachments/assets/fbcf3408-d568-4fd9-ad1e-71e31c64090a"
  alt="main-togo-screen-shot"
  width="738">
</p>

As you type, Togo searches through your tasks and filters the results.

##### c) Shell completion integration

If you've installed shell completion (see below), you can use:

```bash
togo toggle [TAB]
```

Your shell will present available tasks. Type a few letters to filter by name:

```bash
togo toggle me[TAB]
```

Shell will show only tasks containing "me" - perfect for quick selection.

### Available Commands

- `togo add "Task description"` - Add a new task
- `togo toggle [task]` - Toggle completion status
- `togo archive [task]` - Archive a completed task
- `togo unarchive [task]` - Restore an archived task
- `togo delete [task]` - Remove a task permanently
- `togo list [flags]` - View tasks (--all, --archived)

### Additional Options 📌

Every command supports `-h` or `--help` flags to display detailed usage information:

```bash
togo toggle --help
togo add -h
```

## Features In Depth 🧠

### Shell Completion

Setting up shell completion makes Togo even more efficient by enabling tab completion for commands and tasks.

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

#### PowerShell

```powershell
yeah, i don't think you PS guys need this tool :)
```

### Data Storage

Togo stores all your data in a simple JSON file at `~/.togo/todos.json`.

## Built With 🔧

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Cobra](https://img.shields.io/badge/Cobra-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://github.com/spf13/cobra)
[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-FF75B7?style=for-the-badge&logo=go&logoColor=white)](https://github.com/charmbracelet/bubbletea)
[![Lip Gloss](https://img.shields.io/badge/Lip%20Gloss-FFABE7?style=for-the-badge&logo=go&logoColor=white)](https://github.com/charmbracelet/lipgloss)

## License

This project is licensed under MIT - see the [LICENSE](LICENSE) file for details.
