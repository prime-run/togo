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
	showHelp         bool
	sourceLabel      string
	todoFileName     string
	projectName      string
}

func (m TodoTableModel) GetSourceLabel() string {
	return m.sourceLabel
}

func (m TodoTableModel) GetTodoList() *model.TodoList {
	return m.todoList
}

func (m *TodoTableModel) SetSource(label, filename string) {
	m.sourceLabel = label
	m.todoFileName = filename
	
	if label == "project" {
		if projectName, hasProject := model.GetProjectRootName(); hasProject {
			m.projectName = projectName
		} else {
			m.projectName = ""
		}
	} else {
		m.projectName = ""
	}
}
