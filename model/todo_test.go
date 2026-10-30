package model

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeNextID(t *testing.T) {
	tests := []struct {
		name     string
		ids      []int
		nextID   int
		expected int
	}{
		{"stale counter", []int{1, 50}, 2, 51},
		{"counter already free", []int{1, 50}, 51, 51},
		{"counter past every id", []int{1, 50}, 90, 90},
		{"empty list", nil, 1, 1},
		{"zero counter", []int{3}, 0, 4},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tl := NewTodoList()
			for _, id := range tc.ids {
				tl.Todos = append(tl.Todos, Todo{ID: id})
			}
			tl.NextID = tc.nextID

			tl.normalizeNextID()

			if tl.NextID != tc.expected {
				t.Fatalf("NextID = %d, want %d", tl.NextID, tc.expected)
			}
		})
	}
}

func TestAddSkipsExistingIDs(t *testing.T) {
	tl := NewTodoList()
	tl.Todos = []Todo{{ID: 1, Title: "one"}, {ID: 50, Title: "fifty"}}
	tl.rebuildIndex()
	tl.NextID = 2

	todo := tl.Add("new")

	if todo.ID != 51 {
		t.Fatalf("Add assigned id %d, want 51", todo.ID)
	}
	if tl.NextID != 52 {
		t.Fatalf("NextID = %d, want 52", tl.NextID)
	}
	if got := len(tl.Todos); got != 3 {
		t.Fatalf("got %d todos, want 3", got)
	}
	if idx, ok := tl.TodoByID[50]; !ok || tl.Todos[idx].Title != "fifty" {
		t.Fatal("TodoByID[50] no longer points at the original task")
	}
}

func TestLoadNormalizesStaleNextID(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := os.WriteFile(filepath.Join(dir, ".togo"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	data := `{"todos":[{"id":1,"title":"one"},{"id":103,"title":"done"}]` +
		`,"next_id":103}`
	if err := os.WriteFile(filepath.Join(dir, "todos.json"), []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	tl, err := LoadTodoListWithSource("todos.json", "project")
	if err != nil {
		t.Fatal(err)
	}
	if tl.NextID != 104 {
		t.Fatalf("NextID = %d, want 104", tl.NextID)
	}

	added := tl.Add("new")
	if added.ID != 104 {
		t.Fatalf("Add assigned id %d, want 104", added.ID)
	}
}
