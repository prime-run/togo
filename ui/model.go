package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/prime-run/togo/model"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeViewDetail
	ModeDeleteConfirm
	ModeArchiveConfirm
	ModeAddTask
)

type TodoTableModel struct {
	todoList         *model.TodoList
	table            table.Model
	err              error
	mode             Mode
	confirmAction    string
	actionTitle      string
	viewTaskID       int
	width            int
	height           int
	selectedTodoIDs  map[int]bool
	bulkActionActive bool
	textInput        textinput.Model
	showArchived     bool
	showAll          bool
	showArchivedOnly bool
	statusMessage    string
}
