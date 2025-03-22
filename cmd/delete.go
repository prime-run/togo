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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <title>",
	Short: "Delete a todo",
	Long:  `Delete a todo from your list using its title.`,
	Run: func(cmd *cobra.Command, args []string) {
		todoList, err := model.LoadTodoList(TodoFileName)
		if err != nil {
			fmt.Println("Error loading todos:", err)
			os.Exit(1)
		}

		// If no todos, show error
		if len(todoList.Todos) == 0 {
			fmt.Println("No todos found. Add some todos with the 'add' command.")
			os.Exit(1)
		}

		var todoTitle string

		// If an argument was provided
		if len(args) > 0 {
			todoTitle = args[0]

			// Try direct match first (case insensitive)
			todo, found := todoList.FindByTitle(todoTitle, false)
			if found {
				// Found an exact match
				if confirmDelete(todo.Title) {
					todoList.Delete(todo.ID)
					if err := todoList.Save(TodoFileName); err != nil {
						fmt.Println("Error saving todos:", err)
						os.Exit(1)
					}
					fmt.Printf("Todo \"%s\" deleted successfully\n", todo.Title)
				}
				return
			}

			// If no direct match, try to see if it's an ID (for backward compatibility)
			id, err := strconv.Atoi(todoTitle)
			if err == nil {
				for _, todo := range todoList.Todos {
					if todo.ID == id {
						if confirmDelete(todo.Title) {
							todoList.Delete(id)
							if err := todoList.Save(TodoFileName); err != nil {
								fmt.Println("Error saving todos:", err)
								os.Exit(1)
							}
							fmt.Printf("Todo \"%s\" deleted successfully\n", todo.Title)
						}
						return
					}
				}
			}

			// Find matching titles for suggestion (case insensitive)
			var matches []model.Todo
			for _, todo := range todoList.Todos {
				if strings.Contains(strings.ToLower(todo.Title), strings.ToLower(todoTitle)) {
					matches = append(matches, todo)
				}
			}

			if len(matches) == 0 {
				fmt.Printf("Error: No todos found matching \"%s\"\n", todoTitle)
				os.Exit(1)
			} else if len(matches) == 1 {
				// Single match
				if confirmDelete(matches[0].Title) {
					todoList.Delete(matches[0].ID)
					if err := todoList.Save(TodoFileName); err != nil {
						fmt.Println("Error saving todos:", err)
						os.Exit(1)
					}
					fmt.Printf("Todo \"%s\" deleted successfully\n", matches[0].Title)
				}
				return
			} else {
				// Multiple matches - show interactive selection
				selectedTodo, err := selectTodo(matches)
				if err != nil {
					fmt.Println("Operation cancelled")
					os.Exit(0)
				}

				if confirmDelete(selectedTodo.Title) {
					todoList.Delete(selectedTodo.ID)
					if err := todoList.Save(TodoFileName); err != nil {
						fmt.Println("Error saving todos:", err)
						os.Exit(1)
					}
					fmt.Printf("Todo \"%s\" deleted successfully\n", selectedTodo.Title)
				}
				return
			}
		} else {
			// No argument provided - show all todos for selection
			todos := todoList.Todos
			if len(todos) == 0 {
				fmt.Println("No todos found. Add some todos with the 'add' command.")
				os.Exit(0)
			}

			selectedTodo, err := selectTodo(todos)
			if err != nil {
				fmt.Println("Operation cancelled")
				os.Exit(0)
			}

			if confirmDelete(selectedTodo.Title) {
				todoList.Delete(selectedTodo.ID)
				if err := todoList.Save(TodoFileName); err != nil {
					fmt.Println("Error saving todos:", err)
					os.Exit(1)
				}
				fmt.Printf("Todo \"%s\" deleted successfully\n", selectedTodo.Title)
			}
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

		titles := todoList.GetTodoTitles()

		// Filter titles based on toComplete (case insensitive)
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

// selectTodo presents an interactive selection prompt for todos
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

// confirmDelete asks the user to confirm deletion
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
