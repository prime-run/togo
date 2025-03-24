package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ashkansamadiyan/togo/model"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Bold(false)
	statusCompleteStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("28"))
	statusPendingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("136"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
	confirmStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2).
			Margin(1, 0).
			Width(60).
			Align(lipgloss.Center)
	confirmTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Bold(true).
				Margin(1, 0)
	confirmBtnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("236")).
			Padding(0, 1).
			MarginRight(1)
	cancelBtnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)
	fullScreenStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(2)
	fullTaskViewStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")).
				Padding(1, 2).
				Width(60)
	taskTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("252")).
			MarginBottom(1)
	taskContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				MarginBottom(1)
	checkboxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
	selectedCheckboxStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("34"))
	createdAtStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246"))
	archivedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
	inputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2).
			Width(60)
	inputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Bold(true).
				MarginBottom(1)
)

const (
	checkboxEmpty  = "[ ]"
	checkboxFilled = "[×]"
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
}

func NewTodoTable(todoList *model.TodoList) TodoTableModel {
	displayWidth := 80
	checkboxColWidth := 5
	statusColWidth := 15
	createdAtColWidth := 15
	titleColWidth := displayWidth - checkboxColWidth - statusColWidth - createdAtColWidth - 8
	columns := []table.Column{
		{Title: "✓", Width: checkboxColWidth},
		{Title: "Title", Width: titleColWidth},
		{Title: "Status", Width: statusColWidth},
		{Title: "Created", Width: createdAtColWidth},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("252"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("236")).
		Bold(true)
	t.SetStyles(s)
	ti := textinput.New()
	ti.Placeholder = "Enter new task title"
	ti.Focus()
	ti.CharLimit = 120
	ti.Width = titleColWidth
	showArchived := false
	for _, todo := range todoList.Todos {
		if todo.Archived {
			showArchived = true
			break
		}
	}
	m := TodoTableModel{
		todoList:         todoList,
		table:            t,
		mode:             ModeNormal,
		confirmAction:    "",
		actionTitle:      "",
		viewTaskID:       0,
		width:            displayWidth,
		height:           24,
		selectedTodoIDs:  make(map[int]bool),
		bulkActionActive: false,
		textInput:        ti,
		showArchived:     showArchived,
		showAll:          true,
		showArchivedOnly: false,
	}
	m.updateRows()
	return m
}

// SetShowArchivedOnly sets the filter to show only archived todos
func (m *TodoTableModel) SetShowArchivedOnly(show bool) {
	m.showArchivedOnly = show
	m.showAll = false
	m.updateRows()
}

// SetShowAll sets the filter to show all todos
func (m *TodoTableModel) SetShowAll(show bool) {
	m.showAll = show
	m.showArchivedOnly = false
	m.updateRows()
}

// SetShowActiveOnly sets the filter to show only active (non-archived) todos
func (m *TodoTableModel) SetShowActiveOnly(show bool) {
	m.showAll = false
	m.showArchivedOnly = false
	m.updateRows()
}

func (m *TodoTableModel) updateRows() {
	tableBorderWidth := 4
	checkboxColWidth := 5
	statusColWidth := 15
	createdAtColWidth := 15
	availableWidth := m.width - tableBorderWidth
	titleColWidth := availableWidth - checkboxColWidth - statusColWidth - createdAtColWidth - 6
	if titleColWidth < 20 {
		titleColWidth = 20
		if availableWidth-titleColWidth-statusColWidth-createdAtColWidth-6 < checkboxColWidth {
			statusColWidth = 12
			createdAtColWidth = 12
		}
	}
	totalTableWidth := checkboxColWidth + titleColWidth + statusColWidth + createdAtColWidth + tableBorderWidth + 6
	if totalTableWidth > m.width {
		titleColWidth = m.width - checkboxColWidth - statusColWidth - createdAtColWidth - tableBorderWidth - 6
		if titleColWidth < 10 {
			titleColWidth = 10
		}
	}
	m.textInput.Width = titleColWidth
	m.table.SetColumns([]table.Column{
		{Title: "✓", Width: checkboxColWidth},
		{Title: "Title", Width: titleColWidth},
		{Title: "Status", Width: statusColWidth},
		{Title: "Created", Width: createdAtColWidth},
	})
	var rows []table.Row
	var filteredTodos []model.Todo

	// Apply filters
	if m.showAll {
		filteredTodos = m.todoList.Todos
	} else if m.showArchivedOnly {
		filteredTodos = m.todoList.GetArchivedTodos()
	} else {
		filteredTodos = m.todoList.GetActiveTodos()
	}

	for _, todo := range filteredTodos {
		checkbox := checkboxEmpty
		if m.selectedTodoIDs[todo.ID] {
			checkbox = checkboxFilled
		}
		title := todo.Title
		if todo.Archived {
			title = archivedStyle.Render(title)
		}
		var status string
		if todo.Completed {
			status = statusCompleteStyle.Render("Completed")
		} else {
			status = statusPendingStyle.Render("Pending")
		}
		createdAt := model.FormatTimeAgo(todo.CreatedAt)
		rows = append(rows, table.Row{checkbox, title, status, createdAt})
	}
	m.table.SetRows(rows)
}

func (m TodoTableModel) findTodoByID(id int) *model.Todo {
	for i, todo := range m.todoList.Todos {
		if todo.ID == id {
			return &m.todoList.Todos[i]
		}
	}
	return nil
}

func (m TodoTableModel) findTodoByTitle(title string) *model.Todo {
	for i, todo := range m.todoList.Todos {
		if todo.Title == title {
			return &m.todoList.Todos[i]
		}
	}
	return nil
}

func (m TodoTableModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TodoTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = msg.Width
		m.height = msg.Height
		m.updateRows()
	}
	switch m.mode {
	case ModeViewDetail:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "q", "enter":
				m.mode = ModeNormal
				return m, nil
			}
		}
		return m, nil
	case ModeDeleteConfirm, ModeArchiveConfirm:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "Y":
				if m.mode == ModeDeleteConfirm {
					if len(m.selectedTodoIDs) > 0 && m.bulkActionActive {
						for id := range m.selectedTodoIDs {
							m.todoList.Delete(id)
						}
						m.selectedTodoIDs = make(map[int]bool)
						m.bulkActionActive = false
					} else {
						found := false
						for i, todo := range m.todoList.Todos {
							if todo.Title == m.actionTitle || strings.Contains(m.actionTitle, todo.Title) {
								m.todoList.Todos = append(m.todoList.Todos[:i], m.todoList.Todos[i+1:]...)
								found = true
								break
							}
						}
						if !found {
							id, err := strconv.Atoi(m.actionTitle)
							if err == nil {
								m.todoList.Delete(id)
							}
						}
					}
				} else if m.mode == ModeArchiveConfirm {
					if len(m.selectedTodoIDs) > 0 && m.bulkActionActive {
						for id := range m.selectedTodoIDs {
							m.todoList.Archive(id)
						}
						m.selectedTodoIDs = make(map[int]bool)
						m.bulkActionActive = false
					} else {
						for _, todo := range m.todoList.Todos {
							if todo.Title == m.actionTitle {
								m.todoList.Archive(todo.ID)
								break
							}
						}
					}
				}
				m.updateRows()
				m.mode = ModeNormal
				return m, nil
			case "n", "N", "esc", "q":
				m.mode = ModeNormal
				return m, nil
			}
		}
		return m, nil
	case ModeAddTask:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				title := strings.TrimSpace(m.textInput.Value())
				if title != "" {
					m.todoList.Add(title)
					m.textInput.Reset()
					m.updateRows()
				}
				m.mode = ModeNormal
				return m, nil
			case "esc":
				m.textInput.Reset()
				m.mode = ModeNormal
				return m, nil
			}
		}
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	case ModeNormal:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "q":
				return m, tea.Quit
			case "enter":
				if len(m.table.Rows()) > 0 {
					selectedTitle := m.table.SelectedRow()[1]
					// strip archived formatting from the title if necessary
					//NOTE: if this project ctually gets more ppular i might drop lipgloss and raw dog styles! I lowkey feel like im writing js!!!!
					cleanTitle := strings.Replace(selectedTitle, archivedStyle.Render(""), "", -1)

					for _, todo := range m.todoList.Todos {
						if strings.Contains(selectedTitle, todo.Title) || todo.Title == cleanTitle {
							m.mode = ModeViewDetail
							m.viewTaskID = todo.ID
							break
						}
					}
				}
			case "t":
				if len(m.table.Rows()) > 0 {
					if len(m.selectedTodoIDs) > 0 && m.bulkActionActive {
						for id := range m.selectedTodoIDs {
							todo := m.findTodoByID(id)
							if todo != nil {
								if todo.Archived {
									m.todoList.Unarchive(id)
								} else {
									m.todoList.Toggle(id)
								}
							}
						}
					} else {
						selectedTitle := m.table.SelectedRow()[1]
						cleanTitle := strings.Replace(selectedTitle, archivedStyle.Render(""), "", -1)

						for _, todo := range m.todoList.Todos {
							if strings.Contains(selectedTitle, todo.Title) || todo.Title == cleanTitle {
								if todo.Archived {
									m.todoList.Unarchive(todo.ID)
								} else {
									m.todoList.Toggle(todo.ID)
								}
								break
							}
						}
					}
					m.updateRows()
				}
			case "n":
				if len(m.table.Rows()) > 0 {
					if len(m.selectedTodoIDs) > 0 && m.bulkActionActive {
						for id := range m.selectedTodoIDs {
							todo := m.findTodoByID(id)
							if todo != nil {
								if todo.Archived {
									m.todoList.Unarchive(id)
								} else {
									m.todoList.Archive(id)
								}
							}
						}
						m.updateRows()
					} else {
						selectedTitle := m.table.SelectedRow()[1]
						cleanTitle := strings.Replace(selectedTitle, archivedStyle.Render(""), "", -1)

						for _, todo := range m.todoList.Todos {
							if strings.Contains(selectedTitle, todo.Title) || todo.Title == cleanTitle {
								if todo.Archived {
									m.todoList.Unarchive(todo.ID)
								} else {
									m.todoList.Archive(todo.ID)
								}
								m.updateRows()
								break
							}
						}
					}
				}
			case "a":
				m.mode = ModeAddTask
				m.textInput.Focus()
				return m, textinput.Blink
			case "d":
				if len(m.table.Rows()) > 0 {
					if len(m.selectedTodoIDs) > 0 && m.bulkActionActive {
						m.mode = ModeDeleteConfirm
						m.confirmAction = "delete"
					} else {
						selectedTitle := m.table.SelectedRow()[1]
						cleanTitle := strings.Replace(selectedTitle, archivedStyle.Render(""), "", -1)

						m.mode = ModeDeleteConfirm
						m.confirmAction = "delete"
						m.actionTitle = cleanTitle
					}
				}
			case " ":
				if len(m.table.Rows()) > 0 {
					selectedTitle := m.table.SelectedRow()[1]
					cleanTitle := strings.Replace(selectedTitle, archivedStyle.Render(""), "", -1)

					for _, todo := range m.todoList.Todos {
						if strings.Contains(selectedTitle, todo.Title) || todo.Title == cleanTitle {
							if m.selectedTodoIDs[todo.ID] {
								delete(m.selectedTodoIDs, todo.ID)
							} else {
								m.selectedTodoIDs[todo.ID] = true
							}
							break
						}
					}
					m.bulkActionActive = len(m.selectedTodoIDs) > 0
					m.updateRows()
					return m, nil
				}
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m TodoTableModel) View() string {
	if m.mode == ModeViewDetail {
		todo := m.findTodoByID(m.viewTaskID)
		if todo == nil {
			return fullScreenStyle.Width(m.width).Height(m.height).Render(
				fullTaskViewStyle.Render("Task not found."))
		}
		status := "Pending"
		if todo.Completed {
			status = statusCompleteStyle.Render("Completed")
		} else {
			status = statusPendingStyle.Render("Pending")
		}
		archivedStatus := ""
		if todo.Archived {
			archivedStatus = "\nArchived: " + archivedStyle.Render("Yes")
		}
		createdAt := model.FormatTimeAgo(todo.CreatedAt)
		taskView := fullTaskViewStyle.Render(
			taskTitleStyle.Render(todo.Title) + "\n\n" +
				"Status: " + status + archivedStatus + "\n" +
				"Created: " + createdAtStyle.Render(createdAt) + "\n\n" +
				helpStyle.Render("Press Enter to go back"))
		return fullScreenStyle.Width(m.width).Height(m.height).Render(taskView)
	}
	if m.mode == ModeDeleteConfirm || m.mode == ModeArchiveConfirm {
		var confirmMessage string
		action := "delete"
		if m.mode == ModeArchiveConfirm {
			action = "archive"
		}
		if len(m.selectedTodoIDs) > 0 && m.bulkActionActive {
			confirmMessage = fmt.Sprintf("Are you sure you want to %s %d selected tasks?", action, len(m.selectedTodoIDs))
		} else {
			confirmMessage = fmt.Sprintf("Are you sure you want to %s task: \"%s\"?", action, m.actionTitle)
		}
		confirmBox := confirmStyle.Render(
			confirmTextStyle.Render(confirmMessage) + "\n\n" +
				confirmBtnStyle.Render("Y - Yes") + " " + cancelBtnStyle.Render("N - No"))
		return fullScreenStyle.Width(m.width).Height(m.height).Render(confirmBox)
	}
	if m.mode == ModeAddTask {
		inputView := inputStyle.Render(
			inputPromptStyle.Render("Add New Task") + "\n\n" +
				m.textInput.View() + "\n\n" +
				helpStyle.Render("Press Enter to save, Esc to cancel"))
		return fullScreenStyle.Width(m.width).Height(m.height).Render(inputView)
	}
	if len(m.todoList.Todos) == 0 {
		return baseStyle.Render("No tasks found. Press 'a' to add a new task!")
	}
	var helpText string
	var listTitle string

	if m.showArchivedOnly {
		listTitle = "Archived Tasks"
	} else if m.showAll {
		listTitle = "All Tasks"
	} else {
		listTitle = "Active Tasks"
	}

	if m.bulkActionActive {
		helpText = "\n" + listTitle + " - Bulk Mode:" +
			"\n→ t: toggle completion for all selected" +
			"\n→ n: toggle archive/unarchive for selected" +
			"\n→ d: delete selected" +
			"\n→ space: toggle selection" +
			"\n→ enter: view details" +
			"\n→ a: add new task" +
			"\n→ q: quit"
	} else {
		helpText = "\n" + listTitle + ":" +
			"\n→ t: toggle completion" +
			"\n→ n: toggle archive/unarchive" +
			"\n→ d: delete" +
			"\n→ space: select" +
			"\n→ enter: view details" +
			"\n→ a: add new task" +
			"\n→ q: quit"
	}
	help := helpStyle.Render(helpText)
	return baseStyle.Render(m.table.View()) + help
}
