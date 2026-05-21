package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/prime-run/togo/model"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize project-local storage in the current directory",
	Long:  "Create a .togo marker and an empty todos.json in the current directory for project-local storage.",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error determining current directory:", err)
			os.Exit(1)
		}

		markerPath := filepath.Join(cwd, ".togo")
		todosPath := filepath.Join(cwd, TodoFileName)

		created := false

		if _, err := os.Stat(markerPath); os.IsNotExist(err) {
			if err := os.WriteFile(markerPath, []byte{}, 0644); err != nil {
				fmt.Println("Error writing .togo marker:", err)
				os.Exit(1)
			}
			created = true
		}

		if _, err := os.Stat(todosPath); os.IsNotExist(err) {
			tlist := model.NewTodoList()
			data, err := json.Marshal(tlist)
			if err != nil {
				fmt.Println("Error creating initial data:", err)
				os.Exit(1)
			}
			if err := os.WriteFile(todosPath, data, 0644); err != nil {
				fmt.Println("Error writing todos.json:", err)
				os.Exit(1)
			}
			created = true
		}

		if !created {
			fmt.Println("Project storage already initialized in:", cwd)
			return
		}

		fmt.Println("Initialized project storage in:", cwd)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
