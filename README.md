# Togo CLI

A modern, command-line todo application built with Go, featuring a terminal UI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea) and structured with [Cobra CLI](https://github.com/spf13/cobra).

![Togo CLI Screenshot](docs/screenshot-1.png)

[![Togo CLI Demo Video](docs/screenshot-1.png)](docs/demo.mp4)

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

## Installation

### Download Binary

Coming soon!

### Build from Source

1. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/togo.git
   cd togo
   ```

2. Build the executable:
   ```bash
   go build
   ```

3. Move the executable to a directory in your PATH:
   ```bash
   # Example
   sudo mv togo /usr/local/bin/
   ```

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

#### Archive/Unarchive Tasks

Archive a completed task by ID:
```bash
togo archive 1
```

Unarchive a task by ID:
```bash
togo unarchive 1
```

#### Delete Tasks

Delete a task by ID:
```bash
togo delete 1
```

## Data Storage

Togo stores your todos in a JSON file located at `~/.togo/todos.json`. This allows your todos to persist between sessions and be accessible from anywhere on your system.

## Keyboard Shortcuts (Interactive Mode)

| Key       | Action                    |
|-----------|---------------------------|
| ↑/↓       | Navigate between todos    |
| Enter     | Toggle completion status  |
| a         | Archive selected todo     |
| u         | Unarchive selected todo   |
| d         | Delete selected todo      |
| q         | Quit and save changes     |

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 