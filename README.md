# Togo

A command-line todo application built in Go, featuring a terminal UI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea) and structured with [Cobra](https://github.com/spf13/cobra).

![togo CLI Screenshot, (3 different features of it ) ](https://github.com/user-attachments/assets/7907d938-06ae-418a-b44c-96581e3edb1c)

[![togo CLI Demo Video, (demo vid) ](https://github.com/user-attachments/assets/14afdab1-2f6b-419c-9ace-958d8c167646)](https://github.com/user-attachments/assets/14afdab1-2f6b-419c-9ace-958d8c167646)

## Why Togo?

Sometimes you remember something or want to add a todo but you're busy working on other stuff and you shouldn't get distracted. Nothing is ever closer to you than your terminal, so just dump your todo with `togo add` and decide what to do with it later! Stay focused and never forget those important tasks.

Togo helps you:
- Capture thoughts instantly without breaking your workflow
- Manage your tasks from where you already work - the terminal
- Keep a clear mind by getting todos out of your head and into a system
- Quickly organize, prioritize, and track your tasks with minimal friction

## Features

- **Terminal-based UI**: Beautiful interactive interface for managing todos
- **Todo Management**:
  - Add new tasks quickly from the command line
  - Toggle completion status with a keystroke
  - Archive completed tasks to keep your active list clean
  - Unarchive tasks when needed
  - Delete tasks permanently
- **Flexible Viewing Options**:
  - View active todos (default)
  - View archived todos
  - View all todos together
- **Data Persistence**: All data saved locally to `~/.togo/todos.json`
- **Time Tracking**: Created time with relative time display (e.g., "2h ago")
- **Single json db**: ```bash ~/.togo/todos.json```
  ![250322_21h23m36s_screenshot](https://github.com/user-attachments/assets/7edd1331-9ae2-4362-87f5-e51e0bf1089c)

## Usage

### Basic Commands

#### Running without arguments (Interactive Mode)

```bash
togo
```
This launches the interactive table view of your active todos, where you can:
- Navigate with arrow keys
- Press ENTER to toggle completion status
- Press 'd' to delete a todo
- Press 'a' to archive a todo
- Press 'q' to quit and save changes

#### Adding Tasks

```bash
togo add "Buy groceries"
togo add "Finish project report by Friday"
```

#### Listing Tasks

List active todos (default):
```bash
togo list
```

List archived todos:
```bash
togo list --archived
# or
togo list -a
```

List all todos (both active and archived):
```bash
togo list --all
```

#### Toggle Completion Status

Toggle a task as complete/incomplete by ID:
```bash
togo toggle 1
```

You can also toggle by typing part of the task name:
```bash
togo toggle meeting
```
If there's only one task containing "meeting", it will toggle that task immediately without showing a selection list.

Or simply type `togo toggle` and press ENTER to get an interactive selection list:

![Small selection list for task selection](./pics/small-list.png)

As you type, Togo will search through your tasks and show matching results. If you have shell completion installed (see Shell Completion section below), you can press TAB instead of ENTER to get intelligent suggestions directly from your shell.

#### Archive/Unarchive Tasks

Archive a completed task by ID:
```bash
togo archive 1
```

You can archive using part of the task name:
```bash
togo archive report
```
If "report" uniquely identifies a task, it will be archived immediately.

Or use `togo archive` and press ENTER for an interactive selection list of active todos.

Unarchive a task by ID:
```bash
togo unarchive 1
```

The same partial matching works for unarchiving:
```bash
togo unarchive grocery
```

Or use `togo unarchive` and press ENTER for an interactive selection list of archived todos.

#### Delete Tasks

Delete a task by ID:
```bash
togo delete 1
```

Delete using part of the task name:
```bash
togo delete vim
```
If there's only one task containing "vim", it will be deleted immediately.

Or use `togo delete` and press ENTER for an interactive selection.

### Keyboard Shortcuts (Interactive Mode)

| Key       | Action                    |
|-----------|---------------------------|
| ↑/↓       | Navigate between todos    |
| Enter     | Toggle completion status  |
| a         | Archive selected todo     |
| u         | Unarchive selected todo   |
| d         | Delete selected todo      |
| q         | Quit and save changes     |

### Data Storage

Togo stores your todos in a JSON file located at `~/.togo/todos.json`. This allows your todos to persist between sessions and be accessible from anywhere on your system.

## Installation

### Via Go Install

The simplest way to install Togo:

```bash
go install github.com/ashkansamadiyan/togo@latest
```

### From Source

1. Clone this repository:
   ```bash
   git clone https://github.com/ashkansamadiyan/togo.git
   cd togo
   ```

2. Install using Make (recommended):
   ```bash
   # Install to GOPATH/bin
   make install
   
   # Or install system-wide (requires sudo)
   make install-system
   ```

3. Or build manually:
   ```bash
   # Build the executable
   go build
   
   # Copy to a directory in your PATH
   sudo cp togo /usr/local/bin/
   ```

### Shell Completion

Togo supports shell completion scripts for bash, zsh, fish, and PowerShell. Here's how to set it up for your shell:

#### Bash

1. First, ensure bash-completion is installed:
   ```bash
   # For Arch Linux
   sudo pacman -S bash-completion
   
   # For Ubuntu/Debian
   sudo apt install bash-completion
   
   # For Fedora
   sudo dnf install bash-completion
   ```

2. Make sure it's sourced in your ~/.bashrc:
   ```bash
   echo "[ -r /usr/share/bash-completion/bash_completion ] && . /usr/share/bash-completion/bash_completion" >> ~/.bashrc
   source ~/.bashrc
   ```

3. Generate and install the completion script:
   ```bash
   # System-wide (requires sudo)
   togo completion bash > /tmp/togo
   sudo mv /tmp/togo /etc/bash_completion.d/togo
   
   # OR for current user only
   togo completion bash > ~/.bash_completion
   source ~/.bash_completion
   ```

#### Zsh

1. Create a custom completion directory and add to fpath:
   ```bash
   mkdir -p ~/.zsh/completion
   echo "fpath=(~/.zsh/completion \$fpath)" >> ~/.zshrc
   ```

2. Enable completions if not already enabled:
   ```bash
   echo "autoload -Uz compinit && compinit" >> ~/.zshrc
   ```

3. Generate and install the completion script:
   ```bash
   togo completion zsh > ~/.zsh/completion/_togo
   source ~/.zshrc
   ```

#### Fish

```bash
# Create the completions directory if it doesn't exist
mkdir -p ~/.config/fish/completions

# Generate and install the completion script
togo completion fish > ~/.config/fish/completions/togo.fish
```

#### PowerShell

```powershell
# Generate the completion script
togo completion powershell > ~/togo.ps1

# Add to your PowerShell profile
echo ". ~/togo.ps1" >> $PROFILE

# Reload your profile
. $PROFILE
```

#### How Shell Completion Works

Shell completion in Togo is powered by Cobra CLI's built-in completion functionality. When you install the completion script for your shell, it enables intelligent tab completion for:

- All Togo commands and subcommands
- Command flags and options
- Context-aware suggestions for todo items based on the current command

No additional code is needed to enable completion support, as Cobra analyzes Togo's command structure at runtime and generates the appropriate shell-specific scripts. This makes it easy to use Togo without having to remember all available commands and options - just press TAB to see what's available or to quickly complete your command.

After installation, you can verify that completion is working by typing `togo` followed by TAB to see all available commands, or `togo <command>` followed by TAB to see options for a specific command.

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 