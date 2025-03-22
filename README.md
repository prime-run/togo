# Togo CLI

A simple command-line todo application written in Go using [Cobra CLI](https://github.com/spf13/cobra) and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- Add todos from the command line
- View todos in an interactive table
- Toggle completion status of todos
- Delete todos
- Persistent storage (saves to `~/.togo/todos.json`)

## Installation

### Build from Source

1. Clone this repository
2. Run `go build` to build the executable
3. Move the executable to a directory in your PATH

## Usage

### Add a Todo

```bash
togo add Buy groceries
```

### List All Todos

```bash
togo list
```

### Toggle Todo Status

```bash
togo toggle 1
```

### Delete a Todo

```bash
togo delete 1
```

### Interactive Mode

Simply run `togo` without any subcommands to enter interactive mode:

```bash
togo
```

In interactive mode, you can:
- Navigate with arrow keys
- Press ENTER to toggle completion status
- Press d to delete a todo
- Press q to quit

## License

MIT 