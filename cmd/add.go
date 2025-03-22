package cmd
import (
	"fmt"
	"os"
	"strings"
	"togo/model"
	"github.com/spf13/cobra"
)
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new todo",
	Long:  `Add a new todo to your list. The todo will be marked as pending by default.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: Todo title is required")
			fmt.Println("Usage: todooo add <title>")
			os.Exit(1)
		}
		title := strings.Join(args, " ")
		todoList, err := model.LoadTodoList(TodoFileName)
		if err != nil {
			fmt.Println("Error loading todos:", err)
			os.Exit(1)
		}
		todo := todoList.Add(title)
		if err := todoList.Save(TodoFileName); err != nil {
			fmt.Println("Error saving todos:", err)
			os.Exit(1)
		}
		fmt.Printf("Todo added successfully with ID: %d\n", todo.ID)
		fmt.Printf("Title: %s\n", todo.Title)
	},
}
func init() {
	rootCmd.AddCommand(addCmd)
}
