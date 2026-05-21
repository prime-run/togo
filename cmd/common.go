package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/prime-run/togo/model"
	"github.com/prime-run/togo/ui"
)

func loadTodoListOrExit() *model.TodoList {
	todoList, err := model.LoadTodoListWithSource(TodoFileName, sourceFlag)
	if err != nil {
		fmt.Println("Error loading todos:", err)
		os.Exit(1)
	}
	return todoList
}

func saveTodoListOrExit(todoList *model.TodoList) {
	if err := todoList.SaveWithSource(TodoFileName, sourceFlag); err != nil {
		fmt.Println("Error saving todos:", err)
		os.Exit(1)
	}
}

func saveTodoTableAfterTUIOrExit(finalModel tea.Model) {
	m, ok := finalModel.(ui.TodoTableModel)
	if !ok {
		fmt.Println("Error saving todos: unexpected TUI model type")
		os.Exit(1)
	}
	source := m.GetSourceLabel()
	if source == "" {
		source = sourceFlag
	}
	if err := m.GetTodoList().SaveWithSource(TodoFileName, source); err != nil {
		fmt.Println("Error saving todos:", err)
		os.Exit(1)
	}
}

func checkEmptyTodoList(todoList *model.TodoList, emptyMessage string) bool {
	if len(todoList.Todos) == 0 {
		fmt.Println(emptyMessage)
		return true
	}
	return false
}

func handleErrorAndExit(err error, message string) {
	if err != nil {
		fmt.Println(message, err)
		os.Exit(1)
	}
}
