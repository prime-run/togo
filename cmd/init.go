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
	Short: "Initialize a project-local .togo file in the current directory",
	Long:  "Create an empty .togo todos file in the current directory to use project-local storage.",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error determining current directory:", err)
			os.Exit(1)
		}
		path := filepath.Join(cwd, ".togo")
		if _, err := os.Stat(path); err == nil {
			fmt.Println(".togo already exists in:", cwd)
			return
		}

		tlist := model.NewTodoList()
		data, err := json.Marshal(tlist)
		if err != nil {
			fmt.Println("Error creating initial data:", err)
			os.Exit(1)
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			fmt.Println("Error writing .togo:", err)
			os.Exit(1)
		}
		fmt.Println("Initialized .togo in:", cwd)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
