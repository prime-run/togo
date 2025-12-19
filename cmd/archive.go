package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/prime-run/togo/model"
	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive <title>",
	Short: "Archive a todo",
	Long:  `Archive a todo from your list using its title. Archived todos are hidden from the main list.`,
	Run: func(cmd *cobra.Command, args []string) {
		todoList, err := model.LoadTodoList(TodoFileName)
		if err != nil {
			fmt.Println("Error loading todos:", err)
			os.Exit(1)
		}
		if len(todoList.GetActiveTodos()) == 0 {
			fmt.Println("No active todos found. Add some todos with the 'add' command.")
			os.Exit(1)
		}
		if len(args) > 0 {
			todoTitle := args[0]
			todo, found := todoList.FindByTitle(todoTitle, false)
			if found && !todo.Archived {
				todoList.Archive(todo.ID)
				if err := todoList.Save(TodoFileName); err != nil {
					fmt.Println("Error saving todos:", err)
					os.Exit(1)
				}
				fmt.Printf("Todo \"%s\" archived successfully\n", todoTitle)
				return
			}
			id, err := strconv.Atoi(todoTitle)
			if err == nil {
				for _, todo := range todoList.Todos {
					if todo.ID == id && !todo.Archived {
						todoList.Archive(id)
						if err := todoList.Save(TodoFileName); err != nil {
							fmt.Println("Error saving todos:", err)
							os.Exit(1)
						}
						fmt.Printf("Todo \"%s\" archived successfully\n", todo.Title)
						return
					}
				}
			}
			var matches []model.Todo
			for _, todo := range todoList.GetActiveTodos() {
				if strings.Contains(strings.ToLower(todo.Title), strings.ToLower(todoTitle)) {
					matches = append(matches, todo)
				}
			}
			if len(matches) == 0 {
				fmt.Printf("Error: No active todos found matching \"%s\"\n", todoTitle)
				os.Exit(1)
			} else if len(matches) == 1 {
				todoList.Archive(matches[0].ID)
				if err := todoList.Save(TodoFileName); err != nil {
					fmt.Println("Error saving todos:", err)
					os.Exit(1)
				}
				fmt.Printf("Todo \"%s\" archived successfully\n", matches[0].Title)
				return
			} else {
				selectedTodo, err := selectTodoForArchive(matches)
				if err != nil {
					fmt.Println("Operation cancelled")
					os.Exit(0)
				}
				todoList.Archive(selectedTodo.ID)
				if err := todoList.Save(TodoFileName); err != nil {
					fmt.Println("Error saving todos:", err)
					os.Exit(1)
				}
				fmt.Printf("Todo \"%s\" archived successfully\n", selectedTodo.Title)
				return
			}
		} else {
			todos := todoList.GetActiveTodos()
			if len(todos) == 0 {
				fmt.Println("No active todos found. Add some todos with the 'add' command.")
				os.Exit(0)
			}
			selectedTodo, err := selectTodoForArchive(todos)
			if err != nil {
				fmt.Println("Operation cancelled")
				os.Exit(0)
			}
			todoList.Archive(selectedTodo.ID)
			if err := todoList.Save(TodoFileName); err != nil {
				fmt.Println("Error saving todos:", err)
				os.Exit(1)
			}
			fmt.Printf("Todo \"%s\" archived successfully\n", selectedTodo.Title)
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
		activeTitles, _ := todoList.GetActiveAndArchivedTodoTitles()
		if toComplete != "" {
			var filtered []string
			for _, title := range activeTitles {
				if strings.Contains(strings.ToLower(title), strings.ToLower(toComplete)) {
					filtered = append(filtered, title)
				}
			}
			return filtered, cobra.ShellCompDirectiveNoFileComp
		}
		return activeTitles, cobra.ShellCompDirectiveNoFileComp
	},
}

func selectTodoForArchive(todos []model.Todo) (model.Todo, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ .Title | cyan }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
		Inactive: "  {{ .Title }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
		Selected: "✓ {{ .Title | green }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
	}
	prompt := promptui.Select{
		Label:     "Select a todo to archive",
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
	rootCmd.AddCommand(archiveCmd)
}
