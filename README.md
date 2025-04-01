# Togo

<p align="center">
  A command-line todo application built in Go for developers who need to capture ideas without breaking their workflow.
</p>

---

## Why Togo?

Ever been programming and had a brilliant idea or remembered an important task? You know the struggle—interrupting your flow means losing focus, but ignoring it risks forgetting something important. This is where Togo shines, especially for developers and those with ADHD tendencies.

Togo lets you capture those thoughts instantly without breaking your workflow. Your terminal is always just a keystroke away, so just dump your todo with `togo add` and get back to what you were doing. Your brain can relax knowing the idea is safely stored somewhere, and you can maintain your precious focus.

**The core philosophy:** Add it now, manage it later.

---

## Features

- **Zero-friction capture**: Add ideas directly from your terminal without interrupting your flow.
- **Interactive terminal UI**: Manage todos visually with an intuitive interface.
- **VIM keybinds**: Navigate the UI with HJKL motions.
- **Flexible management**: Use either interactive mode or command-line operations to manage tasks.
- **Powerful filtering**: View active, archived, or all todos with simple commands.
- **Search and partial matching**: Quickly find tasks by name or partial matches.
- **Shell integration**: Tab completion for commands and task names.
- **Cross-platform**: Works on Linux, macOS, and Windows.

---

## Installation

### Option 1: Install via Go

The simplest way to install Togo:

```bash
go install github.com/prime-run/togo@latest
```

Make sure `$GOPATH/bin` is in your `PATH` to access the installed binary.

### Option 2: Manual Build

```bash
# Clone the repository
git clone https://github.com/prime-run/togo.git
cd togo

# Build the binary
go build -o togo
```

### Option 3: Using Make (Recommended)

```bash
# Clone the repository
git clone https://github.com/prime-run/togo.git
cd togo

# Install to GOPATH/bin (includes automatic shell completion)
make install

# OR install system-wide (requires sudo)
make install-system
```

---

## Usage

### Capturing Ideas

When inspiration strikes or a task pops into your head, just:

```bash
togo add "Call the client about project scope"
togo add "Research Go concurrency patterns"
```

Then get right back to what you were doing. No more mental juggling or lost ideas.

### Managing Your Tasks

Togo offers two primary ways to manage your tasks:

#### 1. Interactive Mode

Open the interactive UI to work with your todos visually:

```bash
togo
# or
togo list           # Show active todos
togo list --all     # Show all todos
togo list --archived # Show archived todos only
```

The interactive mode shows helpful keyboard shortcuts right in the interface.

#### 2. Command-Line Operations

Togo offers flexible command syntax with three usage patterns:

##### a) Direct selection by partial name

```bash
togo toggle meeting
```

If only one task contains "meeting", it executes immediately. If multiple tasks match, Togo opens a selection list so you can choose.

##### b) Interactive selection list

```bash
togo toggle
```

Opens a selection list where you can choose from available tasks.

##### c) Shell completion integration

If you've installed shell completion, you can use:

```bash
togo toggle [TAB]
```

Your shell will present available tasks. Type a few letters to filter by name.

---

## Available Commands

- `togo add "Task description"` - Add a new task.
- `togo toggle [task]` - Toggle completion status.
- `togo archive [task]` - Archive a completed task.
- `togo unarchive [task]` - Restore an archived task.
- `togo delete [task]` - Remove a task permanently.
- `togo list [flags]` - View tasks (`--all`, `--archived`).

Every command supports `-h` or `--help` flags to display detailed usage information:

```bash
togo toggle --help
togo add -h
```

---

## Shell Completion

Setting up shell completion makes Togo even more efficient by enabling tab completion for commands and tasks.

### Zsh

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

### Bash

```bash
# 1. Ensure completion is sourced
echo "[ -r /usr/share/bash-completion/bash_completion ] && . /usr/share/bash-completion/bash_completion" >> ~/.bashrc
source ~/.bashrc

# 2. Install Togo completion
togo completion bash > ~/.bash_completion
source ~/.bash_completion
```

### Fish

```bash
mkdir -p ~/.config/fish/completions
togo completion fish > ~/.config/fish/completions/togo.fish
```

---

## Data Storage

Togo stores all your data in a simple JSON file at `~/.togo/todos.json`.

---

## Contributing

Contributions are welcome! If you have ideas for new features or improvements, feel free to open an issue or submit a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
