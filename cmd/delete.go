package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/prime-run/togo/model"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <title>",
	Short: "Delete a todo",
	Long:  `Delete a todo from your list using its title.`,
	Run: func(cmd *cobra.Command, args []string) {
		todoList := loadTodoListOrExit()

		if checkEmptyTodoList(todoList, "No todos found. Add some todos with the 'add' command.") {
			return
		}

		var todo *model.Todo
		var err error
		if len(args) > 0 {
			todo, err = findTodoByTitleOrID(todoList, args[0], false)
		} else {
			todo, err = selectTodoFromList(todoList.Todos)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if confirmDelete(todo.Title) {
			todoList.Delete(todo.ID)
			saveTodoListOrExit(todoList)
			fmt.Printf("Todo \"%s\" deleted successfully\n", todo.Title)
		}
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		todoList := loadTodoListOrExit()
		titles := todoList.GetTodoTitles()

		if toComplete != "" {
			var filtered []string
			for _, title := range titles {
				if strings.Contains(strings.ToLower(title), strings.ToLower(toComplete)) {
					filtered = append(filtered, title)
				}
			}
			return filtered, cobra.ShellCompDirectiveNoFileComp
		}
		return titles, cobra.ShellCompDirectiveNoFileComp
	},
}

func selectTodo(todos []model.Todo) (model.Todo, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ .Title | cyan }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
		Inactive: "  {{ .Title }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
		Selected: "✓ {{ .Title | green }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
	}
	prompt := promptui.Select{
		Label:     "Select a todo to delete",
		Items:     todos,
		Templates: templates,
		Size:      10,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return model.Todo{}, err
	}
	return todos[index], nil
}

func confirmDelete(title string) bool {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Are you sure you want to delete \"%s\"", title),
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil {
		return false
	}
	return strings.ToLower(result) == "y"
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
