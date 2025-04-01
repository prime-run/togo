package cmd

import (
	"fmt"
	"os"

	"github.com/prime-run/togo/model"
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
