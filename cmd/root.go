package cmd

import (
	"fmt"
	"strings"

	"github.com/prime-run/togo/config"
	"github.com/prime-run/togo/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var TodoFileName = "todos.json"
var sourceFlag string = "project"
var skipConfirmations bool

var rootCmd = &cobra.Command{
	Use:   "togo",
	Short: "A simple todo application",
	Long:  `A simple todo application that lets you manage your tasks from the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			handleErrorAndExit(err, "Error loading config:")
		}

		// Flag overrides config
		if !skipConfirmations {
			skipConfirmations = cfg.SkipConfirmations
		}

		todoList := loadTodoListOrExit()

		tableModel := ui.NewTodoTable(todoList)
		tableModel.SetSource(sourceFlag, TodoFileName)
		tableModel.SetConfig(cfg)
		tableModel.SkipConfirmationsByDefault = skipConfirmations
		if skipConfirmations {
			tableModel.SetSkipConfirmationsStatus("on")
		} else {
			tableModel.SetSkipConfirmationsStatus("off")
		}
		_, err = tea.NewProgram(tableModel, tea.WithAltScreen()).Run()
		handleErrorAndExit(err, "Error running program:")

		finalSource := tableModel.GetSourceLabel()
		finalList := tableModel.GetTodoList()
		if err := finalList.SaveWithSource(TodoFileName, finalSource); err != nil {
			handleErrorAndExit(err, "Error saving todos:")
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&sourceFlag, "source", "s", "project", "todo source: project or global")
	rootCmd.PersistentFlags().BoolVarP(&skipConfirmations, "skip-confirmations", "y", false, "skip confirmations for delete/archive")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		s := strings.ToLower(strings.TrimSpace(sourceFlag))
		switch s {
		case "project", "global":
			sourceFlag = s
			return nil
		default:
			return fmt.Errorf("invalid value for --source: %q (must be 'project' or 'global')", sourceFlag)
		}
	}

	_ = rootCmd.RegisterFlagCompletionFunc("source", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"project", "global"}, cobra.ShellCompDirectiveNoFileComp
	})

}
