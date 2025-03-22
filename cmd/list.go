package cmd

import (
	"fmt"
	"os"

	"togo/model"
	"togo/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all todos",
	Long: `List all todos in a nice interactive UI.
	
You can use:
- list: to show active todos
- list --archived: to show archived todos
- list --all: to show both active and archived todos`,
	Run: func(cmd *cobra.Command, args []string) {
		todoList, err := model.LoadTodoList(TodoFileName)
		if err != nil {
			fmt.Println("Error loading todos:", err)
			os.Exit(1)
		}

		archivedFlag, _ := cmd.Flags().GetBool("archived")
		allFlag, _ := cmd.Flags().GetBool("all")

		var filteredList *model.TodoList

		if archivedFlag {
			// Show only archived todos
			filteredList = &model.TodoList{
				Todos:  todoList.GetArchivedTodos(),
				NextID: todoList.NextID,
			}
		} else if allFlag {
			// Show all todos (both active and archived)
			filteredList = todoList
		} else {
			// Default: show only active todos
			filteredList = &model.TodoList{
				Todos:  todoList.GetActiveTodos(),
				NextID: todoList.NextID,
			}
		}

		m := ui.NewTodoTable(filteredList)
		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		// Save the todos after any changes
		if err := todoList.Save(TodoFileName); err != nil {
			fmt.Println("Error saving todos:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags for archived and all todos
	listCmd.Flags().BoolP("archived", "a", false, "Show only archived todos")
	listCmd.Flags().Bool("all", false, "Show all todos (both active and archived)")
}
