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
	Long:  "Create a .togo.json file in the current directory for project-local storage.",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error determining current directory:", err)
			os.Exit(1)
		}

		todosPath := filepath.Join(cwd, TodoFileName)

		if _, err := os.Stat(todosPath); !os.IsNotExist(err) {
			fmt.Println("Project storage already initialized in:", cwd)
			return
		}

		tlist := model.NewTodoList()
		data, err := json.Marshal(tlist)
		if err != nil {
			fmt.Println("Error creating initial data:", err)
			os.Exit(1)
		}
		if err := os.WriteFile(todosPath, data, 0644); err != nil {
			fmt.Println("Error writing .togo.json:", err)
			os.Exit(1)
		}

		fmt.Println("Initialized project storage in:", cwd)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
