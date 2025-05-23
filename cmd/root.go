package cmd

import (
	"github.com/prime-run/togo/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var TodoFileName = "todos.json"

var rootCmd = &cobra.Command{
	Use:   "togo",
	Short: "A simple todo application",
	Long:  `A simple todo application that lets you manage your tasks from the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		todoList := loadTodoListOrExit()

		if checkEmptyTodoList(todoList, "No todos found. Add some todos with 'add' command!") {
			return
		}

		tableModel := ui.NewTodoTable(todoList)
		_, err := tea.NewProgram(tableModel, tea.WithAltScreen()).Run()
		handleErrorAndExit(err, "Error running program:")

		saveTodoListOrExit(todoList)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todooo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
