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
var toggleCmd = &cobra.Command{
	Use:   "toggle <title>",
	Short: "Toggle todo completion status",
	Long:  `Toggle the completion status of a todo. It marks a pending todo as completed and vice versa.`,
	Run: func(cmd *cobra.Command, args []string) {
		todoList, err := model.LoadTodoList(TodoFileName)
		if err != nil {
			fmt.Println("Error loading todos:", err)
			os.Exit(1)
		}
		if len(todoList.Todos) == 0 {
			fmt.Println("No todos found. Add some todos with the 'add' command.")
			os.Exit(1)
		}
		if len(args) > 0 {
			todoTitle := args[0]
			todo, found := todoList.FindByTitle(todoTitle, false)
			if found {
				todoList.Toggle(todo.ID)
				status := "Pending"
				for _, t := range todoList.Todos {
					if t.ID == todo.ID {
						if t.Completed {
							status = "Completed"
						}
						break
					}
				}
				if err := todoList.Save(TodoFileName); err != nil {
					fmt.Println("Error saving todos:", err)
					os.Exit(1)
				}
				fmt.Printf("Todo \"%s\" toggled successfully\n", todoTitle)
				fmt.Printf("Status: %s\n", status)
				return
			}
			id, err := strconv.Atoi(todoTitle)
			if err == nil {
				for _, todo := range todoList.Todos {
					if todo.ID == id {
						todoList.Toggle(id)
						var status string
						if !todo.Completed { 
							status = "Completed"
						} else {
							status = "Pending"
						}
						if err := todoList.Save(TodoFileName); err != nil {
							fmt.Println("Error saving todos:", err)
							os.Exit(1)
						}
						fmt.Printf("Todo \"%s\" toggled successfully\n", todo.Title)
						fmt.Printf("Status: %s\n", status)
						return
					}
				}
			}
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
				todoList.Toggle(matches[0].ID)
				status := "Completed"
				if matches[0].Completed {
					status = "Pending"
				}
				if err := todoList.Save(TodoFileName); err != nil {
					fmt.Println("Error saving todos:", err)
					os.Exit(1)
				}
				fmt.Printf("Todo \"%s\" toggled successfully\n", matches[0].Title)
				fmt.Printf("Status: %s\n", status)
				return
			} else {
				selectedTodo, err := selectTodoForToggle(matches)
				if err != nil {
					fmt.Println("Operation cancelled")
					os.Exit(0)
				}
				todoList.Toggle(selectedTodo.ID)
				var status string
				if !selectedTodo.Completed {
					status = "Completed"
				} else {
					status = "Pending"
				}
				if err := todoList.Save(TodoFileName); err != nil {
					fmt.Println("Error saving todos:", err)
					os.Exit(1)
				}
				fmt.Printf("Todo \"%s\" toggled successfully\n", selectedTodo.Title)
				fmt.Printf("Status: %s\n", status)
				return
			}
		} else {
			todos := todoList.Todos
			if len(todos) == 0 {
				fmt.Println("No todos found. Add some todos with the 'add' command.")
				os.Exit(0)
			}
			selectedTodo, err := selectTodoForToggle(todos)
			if err != nil {
				fmt.Println("Operation cancelled")
				os.Exit(0)
			}
			todoList.Toggle(selectedTodo.ID)
			var status string
			if !selectedTodo.Completed {
				status = "Completed"
			} else {
				status = "Pending"
			}
			if err := todoList.Save(TodoFileName); err != nil {
				fmt.Println("Error saving todos:", err)
				os.Exit(1)
			}
			fmt.Printf("Todo \"%s\" toggled successfully\n", selectedTodo.Title)
			fmt.Printf("Status: %s\n", status)
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
func selectTodoForToggle(todos []model.Todo) (model.Todo, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ .Title | cyan }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
		Inactive: "  {{ .Title }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
		Selected: "✓ {{ .Title | green }} {{ if .Completed }}(Completed){{ else }}(Pending){{ end }}",
	}
	prompt := promptui.Select{
		Label:     "Select a todo to toggle status",
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
	rootCmd.AddCommand(toggleCmd)
}
