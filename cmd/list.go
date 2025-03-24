package cmd

import (
	"github.com/ashkansamadiyan/togo/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all todos",
	Long: `List all todos in a nice interactive UI.
You can use:
- list: to show active todos
- list --archived: to show archived todos
- list --all: to show both active and archived todos`,
	Run: func(cmd *cobra.Command, args []string) {
		todoList := loadTodoListOrExit()

		if checkEmptyTodoList(todoList, "No todos found. Add some todos with 'add' command.") {
			return
		}

		archivedFlag, _ := cmd.Flags().GetBool("archived")
		allFlag, _ := cmd.Flags().GetBool("all")

		// Pass the actual todoList to the UI, along with filter information
		m := ui.NewTodoTable(todoList)

		// Set the filter mode
		if archivedFlag {
			m.SetShowArchivedOnly(true)
		} else if allFlag {
			m.SetShowAll(true)
		} else {
			m.SetShowActiveOnly(true)
		}

		_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
		handleErrorAndExit(err, "Error running program:")

		// Save the original todoList which contains all modifications
		saveTodoListOrExit(todoList)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("archived", "a", false, "Show only archived todos")
	listCmd.Flags().Bool("all", false, "Show all todos (both active and archived)")
}
