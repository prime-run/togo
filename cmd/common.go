package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/prime-run/togo/model"
	"github.com/manifoldco/promptui"
)

func loadTodoListOrExit() *model.TodoList {
	todoList, err := model.LoadTodoList(TodoFileName)
	if err != nil {
		fmt.Println("Error loading todos:", err)
		os.Exit(1)
	}
	return todoList
}

func saveTodoListOrExit(todoList *model.TodoList) {
	if err := todoList.Save(TodoFileName); err != nil {
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

func findTodoByTitleOrID(todoList *model.TodoList, input string, archived bool) (*model.Todo, error) {
	// Try to find by title
	todo, found := todoList.FindByTitle(input, false)
	if found && todo.Archived == archived {
		return todo, nil
	}

	// Try to find by ID
	id, err := strconv.Atoi(input)
	if err == nil {
		for _, todo := range todoList.Todos {
			if todo.ID == id && todo.Archived == archived {
				return &todo, nil
			}
		}
	}

	// Search for partial matches
	var matches []model.Todo
	for _, todo := range todoList.Todos {
		if strings.Contains(strings.ToLower(todo.Title), strings.ToLower(input)) && todo.Archived == archived {
			matches = append(matches, todo)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no todos found matching \"%s\"", input)
	} else if len(matches) == 1 {
		return &matches[0], nil
	} else {
		selectedTodo, err := selectTodoFromList(matches)
		if err != nil {
			return nil, fmt.Errorf("operation cancelled")
		}
		return selectedTodo, nil
	}
}

func selectTodoFromList(todos []model.Todo) (*model.Todo, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ .Title | cyan }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
		Inactive: "  {{ .Title }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
		Selected: "✓ {{ .Title | green }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
	}
	prompt := promptui.Select{
		Label:     "Select a todo",
		Items:     todos,
		Templates: templates,
		Size:      10,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &todos[index], nil
}
