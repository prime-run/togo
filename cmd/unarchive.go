package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/prime-run/togo/model"
	"github.com/spf13/cobra"
)

var unarchiveCmd = &cobra.Command{
	Use:   "unarchive <title>",
	Short: "Unarchive a todo",
	Long:  `Unarchive a todo from your archive using its title. This returns it to the active list.`,
	Run: func(cmd *cobra.Command, args []string) {
		todoList, err := model.LoadTodoList(TodoFileName)
		if err != nil {
			fmt.Println("Error loading todos:", err)
			os.Exit(1)
		}

		if len(todoList.GetArchivedTodos()) == 0 {
			fmt.Println("No archived todos found.")
			os.Exit(1)
		}

		var todo *model.Todo
		if len(args) > 0 {
			todo, err = findTodoByTitleOrID(todoList, args[0], true)
		} else {
			todo, err = selectTodoFromList(todoList.GetArchivedTodos())
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		todoList.Unarchive(todo.ID)
		if err := todoList.Save(TodoFileName); err != nil {
			fmt.Println("Error saving todos:", err)
			os.Exit(1)
		}
		fmt.Printf("Todo \"%s\" unarchived successfully\n", todo.Title)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		todoList, err := model.LoadTodoList(TodoFileName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		_, archivedTitles := todoList.GetActiveAndArchivedTodoTitles()
		if toComplete != "" {
			var filtered []string
			for _, title := range archivedTitles {
				if strings.Contains(strings.ToLower(title), strings.ToLower(toComplete)) {
					filtered = append(filtered, title)
				}
			}
			return filtered, cobra.ShellCompDirectiveNoFileComp
		}
		return archivedTitles, cobra.ShellCompDirectiveNoFileComp
	},
}

func selectTodoForUnarchive(todos []model.Todo) (model.Todo, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ .Title | cyan }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
		Inactive: "  {{ .Title }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
		Selected: "✓ {{ .Title | green }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
	}
	prompt := promptui.Select{
		Label:     "Select a todo to unarchive",
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

func init() {
	rootCmd.AddCommand(unarchiveCmd)
}
