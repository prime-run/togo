package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"togo/model"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// unarchiveCmd represents the unarchive command
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

		// If no archived todos, show error
		if len(todoList.GetArchivedTodos()) == 0 {
			fmt.Println("No archived todos found.")
			os.Exit(1)
		}

		// If an argument was provided
		if len(args) > 0 {
			todoTitle := args[0]

			// First try to find the todo by title (case insensitive)
			todo, found := todoList.FindByTitle(todoTitle, false)
			if found && todo.Archived {
				todoList.Unarchive(todo.ID)
				if err := todoList.Save(TodoFileName); err != nil {
					fmt.Println("Error saving todos:", err)
					os.Exit(1)
				}
				fmt.Printf("Todo \"%s\" unarchived successfully\n", todoTitle)
				return
			}

			// If not found by title, try ID for backward compatibility
			id, err := strconv.Atoi(todoTitle)
			if err == nil {
				for _, todo := range todoList.Todos {
					if todo.ID == id && todo.Archived {
						todoList.Unarchive(id)
						if err := todoList.Save(TodoFileName); err != nil {
							fmt.Println("Error saving todos:", err)
							os.Exit(1)
						}
						fmt.Printf("Todo \"%s\" unarchived successfully\n", todo.Title)
						return
					}
				}
			}

			// Find matching titles for selection (case insensitive)
			var matches []model.Todo
			for _, todo := range todoList.GetArchivedTodos() {
				if strings.Contains(strings.ToLower(todo.Title), strings.ToLower(todoTitle)) {
					matches = append(matches, todo)
				}
			}

			if len(matches) == 0 {
				fmt.Printf("Error: No archived todos found matching \"%s\"\n", todoTitle)
				os.Exit(1)
			} else if len(matches) == 1 {
				// Single match
				todoList.Unarchive(matches[0].ID)
				if err := todoList.Save(TodoFileName); err != nil {
					fmt.Println("Error saving todos:", err)
					os.Exit(1)
				}
				fmt.Printf("Todo \"%s\" unarchived successfully\n", matches[0].Title)
				return
			} else {
				// Multiple matches
				selectedTodo, err := selectTodoForUnarchive(matches)
				if err != nil {
					fmt.Println("Operation cancelled")
					os.Exit(0)
				}

				todoList.Unarchive(selectedTodo.ID)
				if err := todoList.Save(TodoFileName); err != nil {
					fmt.Println("Error saving todos:", err)
					os.Exit(1)
				}
				fmt.Printf("Todo \"%s\" unarchived successfully\n", selectedTodo.Title)
				return
			}
		} else {
			// No argument provided - show all archived todos for selection
			todos := todoList.GetArchivedTodos()
			if len(todos) == 0 {
				fmt.Println("No archived todos found.")
				os.Exit(0)
			}

			selectedTodo, err := selectTodoForUnarchive(todos)
			if err != nil {
				fmt.Println("Operation cancelled")
				os.Exit(0)
			}

			todoList.Unarchive(selectedTodo.ID)
			if err := todoList.Save(TodoFileName); err != nil {
				fmt.Println("Error saving todos:", err)
				os.Exit(1)
			}
			fmt.Printf("Todo \"%s\" unarchived successfully\n", selectedTodo.Title)
		}
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		todoList, err := model.LoadTodoList(TodoFileName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// Only suggest archived todos
		_, archivedTitles := todoList.GetActiveAndArchivedTodoTitles()

		// Filter titles based on toComplete (case insensitive)
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

// selectTodoForUnarchive presents an interactive selection prompt for todos to unarchive
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
