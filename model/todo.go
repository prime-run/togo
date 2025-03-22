package model
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	Archived  bool      `json:"archived"`
	CreatedAt time.Time `json:"created_at"`
}
type TodoList struct {
	Todos  []Todo `json:"todos"`
	NextID int    `json:"next_id"`
}
func NewTodoList() *TodoList {
	return &TodoList{
		Todos:  []Todo{},
		NextID: 1,
	}
}
func (tl *TodoList) Add(title string) *Todo {
	todo := Todo{
		ID:        tl.NextID,
		Title:     title,
		Completed: false,
		Archived:  false,
		CreatedAt: time.Now(),
	}
	tl.Todos = append(tl.Todos, todo)
	tl.NextID++
	return &todo
}
func (tl *TodoList) Toggle(id int) bool {
	for i, todo := range tl.Todos {
		if todo.ID == id {
			tl.Todos[i].Completed = !tl.Todos[i].Completed
			return true
		}
	}
	return false
}
func (tl *TodoList) Archive(id int) bool {
	for i, todo := range tl.Todos {
		if todo.ID == id {
			tl.Todos[i].Archived = true
			return true
		}
	}
	return false
}
func (tl *TodoList) Unarchive(id int) bool {
	for i, todo := range tl.Todos {
		if todo.ID == id {
			tl.Todos[i].Archived = false
			return true
		}
	}
	return false
}
func (tl *TodoList) GetActiveTodos() []Todo {
	var activeTodos []Todo
	for _, todo := range tl.Todos {
		if !todo.Archived {
			activeTodos = append(activeTodos, todo)
		}
	}
	return activeTodos
}
func (tl *TodoList) GetArchivedTodos() []Todo {
	var archivedTodos []Todo
	for _, todo := range tl.Todos {
		if todo.Archived {
			archivedTodos = append(archivedTodos, todo)
		}
	}
	return archivedTodos
}
func (tl *TodoList) Delete(id int) bool {
	for i, todo := range tl.Todos {
		if todo.ID == id {
			tl.Todos = append(tl.Todos[:i], tl.Todos[i+1:]...)
			return true
		}
	}
	return false
}
func (tl *TodoList) Save(filename string) error {
	dataDir := getDataDir()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}
	data, err := json.Marshal(tl)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dataDir, filename), data, 0644)
}
func LoadTodoList(filename string) (*TodoList, error) {
	dataDir := getDataDir()
	filePath := filepath.Join(dataDir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return NewTodoList(), nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var tl TodoList
	if err := json.Unmarshal(data, &tl); err != nil {
		return nil, err
	}
	for i, todo := range tl.Todos {
		if todo.CreatedAt.IsZero() {
			tl.Todos[i].CreatedAt = time.Now()
		}
	}
	return &tl, nil
}
func getDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".togo"
	}
	return filepath.Join(homeDir, ".togo")
}
func (tl *TodoList) FindByTitle(title string, caseSensitive bool) (*Todo, bool) {
	for _, todo := range tl.Todos {
		if caseSensitive {
			if todo.Title == title {
				return &todo, true
			}
		} else {
			if strings.EqualFold(todo.Title, title) {
				return &todo, true
			}
		}
	}
	return nil, false
}
func (tl *TodoList) DeleteByTitle(title string, caseSensitive bool) bool {
	for i, todo := range tl.Todos {
		var matches bool
		if caseSensitive {
			matches = todo.Title == title
		} else {
			matches = strings.EqualFold(todo.Title, title)
		}
		if matches {
			tl.Todos = append(tl.Todos[:i], tl.Todos[i+1:]...)
			return true
		}
	}
	return false
}
func (tl *TodoList) GetTodoTitles() []string {
	titles := make([]string, len(tl.Todos))
	for i, todo := range tl.Todos {
		titles[i] = todo.Title
	}
	return titles
}
func (tl *TodoList) GetActiveAndArchivedTodoTitles() ([]string, []string) {
	var activeTitles, archivedTitles []string
	for _, todo := range tl.Todos {
		if todo.Archived {
			archivedTitles = append(archivedTitles, todo.Title)
		} else {
			activeTitles = append(activeTitles, todo.Title)
		}
	}
	return activeTitles, archivedTitles
}
func FormatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	hours := int(diff.Hours())
	minutes := int(diff.Minutes()) % 60
	if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	} else {
		return "now"
	}
}
