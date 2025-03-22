package ui

import (
	"fmt"
	"strings"

	"togo/model"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styling
var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Bold(false)

	statusCompleteStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("28")) // Darker green

	statusPendingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("136")) // Darker yellow/gold

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
			Foreground(lipgloss.Color("246")) // Light gray

	archivedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")) // Gray for archived items

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

// Mode represents the current UI mode
type Mode int

const (
	ModeNormal Mode = iota
	ModeViewDetail
	ModeDeleteConfirm
	ModeArchiveConfirm
	ModeAddTask
)

// TodoTableModel represents the UI state
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
}

func NewTodoTable(todoList *model.TodoList) TodoTableModel {
	// Get display width (use default first, will be updated from terminal)
	displayWidth := 80

	// Calculate dynamic column widths
	checkboxColWidth := 5
	statusColWidth := 15
	createdAtColWidth := 15
	// Let title take the rest of the space, leaving some margin
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

	// Create and configure text input for adding tasks
	ti := textinput.New()
	ti.Placeholder = "Enter new task title"
	ti.Focus()
	ti.CharLimit = 120
	ti.Width = titleColWidth

	// Check if we're showing archived tasks
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
	}
	m.updateRows()

	return m
}

func (m *TodoTableModel) updateRows() {
	// Account for table borders, padding, and margins
	// The bubbles table adds padding and borders that we need to consider
	tableBorderWidth := 4 // 2 for left border + 2 for right border

	// Set dynamic column widths based on current terminal width
	checkboxColWidth := 5
	statusColWidth := 15
	createdAtColWidth := 15

	// Calculate available space for title, ensuring the table fits within terminal
	availableWidth := m.width - tableBorderWidth
	titleColWidth := availableWidth - checkboxColWidth - statusColWidth - createdAtColWidth - 6 // additional padding and separator space

	if titleColWidth < 20 {
		titleColWidth = 20 // Minimum width for title
		// If we had to force a minimum title width, we might need to recalculate other widths
		if availableWidth-titleColWidth-statusColWidth-createdAtColWidth-6 < checkboxColWidth {
			// If we're really tight on space, reduce status column width slightly
			statusColWidth = 12
			createdAtColWidth = 12
		}
	}

	// Final check to make sure we're not exceeding terminal width
	totalTableWidth := checkboxColWidth + titleColWidth + statusColWidth + createdAtColWidth + tableBorderWidth + 6
	if totalTableWidth > m.width {
		// If still too wide, adjust title width as last resort
		titleColWidth = m.width - checkboxColWidth - statusColWidth - createdAtColWidth - tableBorderWidth - 6
		if titleColWidth < 10 {
			titleColWidth = 10 // Absolute minimum
		}
	}

	// Update text input width
	m.textInput.Width = titleColWidth

	// Update column widths
	m.table.SetColumns([]table.Column{
		{Title: "✓", Width: checkboxColWidth},
		{Title: "Title", Width: titleColWidth},
		{Title: "Status", Width: statusColWidth},
		{Title: "Created", Width: createdAtColWidth},
	})

	var rows []table.Row
	for _, todo := range m.todoList.Todos {
		// Create checkbox
		checkbox := checkboxEmpty
		if m.selectedTodoIDs[todo.ID] {
			checkbox = checkboxFilled
		}

		// For title, just use the title without padding - we'll let the table handle it
		title := todo.Title
		// If archived, style differently
		if todo.Archived {
			title = archivedStyle.Render(title)
		}

		// Create status text
		var status string
		if todo.Completed {
			status = statusCompleteStyle.Render("Completed")
		} else {
			status = statusPendingStyle.Render("Pending")
		}

		// Format created at as relative time
		createdAt := model.FormatTimeAgo(todo.CreatedAt)

		rows = append(rows, table.Row{checkbox, title, status, createdAt})
	}
	m.table.SetRows(rows)
}

func (m TodoTableModel) findTodoByID(id int) *model.Todo {
	for _, todo := range m.todoList.Todos {
		if todo.ID == id {
			return &todo
		}
	}
	return nil
}

func (m TodoTableModel) findTodoByTitle(title string) *model.Todo {
	for _, todo := range m.todoList.Todos {
		if todo.Title == title {
			return &todo
		}
	}
	return nil
}

func (m TodoTableModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TodoTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Window size updates
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = msg.Width
		m.height = msg.Height
		m.updateRows()
	}

	// Handle different modes
	switch m.mode {
	case ModeViewDetail:
		// Task detail view mode
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
		// Confirmation mode (for delete or archive)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "Y":
				if m.mode == ModeDeleteConfirm {
					// Handle delete confirmation
					if len(m.selectedTodoIDs) > 0 && m.bulkActionActive {
						// Delete all selected todos
						for id := range m.selectedTodoIDs {
							m.todoList.Delete(id)
						}
						m.selectedTodoIDs = make(map[int]bool) // Clear selections
						m.bulkActionActive = false
					} else {
						// Delete by title
						for i, todo := range m.todoList.Todos {
							if todo.Title == m.actionTitle {
								m.todoList.Todos = append(m.todoList.Todos[:i], m.todoList.Todos[i+1:]...)
								break
							}
						}
					}
				} else if m.mode == ModeArchiveConfirm {
					// Handle archive confirmation
					if len(m.selectedTodoIDs) > 0 && m.bulkActionActive {
						// Archive all selected todos
						for id := range m.selectedTodoIDs {
							m.todoList.Archive(id)
						}
						m.selectedTodoIDs = make(map[int]bool) // Clear selections
						m.bulkActionActive = false
					} else {
						// Archive by title
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
		// Add task mode
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				// Add the new task
				title := strings.TrimSpace(m.textInput.Value())
				if title != "" {
					m.todoList.Add(title)
					m.textInput.Reset()
					m.updateRows()
				}
				m.mode = ModeNormal
				return m, nil
			case "esc":
				// Cancel adding task
				m.textInput.Reset()
				m.mode = ModeNormal
				return m, nil
			}
		}
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd

	case ModeNormal:
		// Normal mode - handle keyboard shortcuts
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "q":
				return m, tea.Quit
			case "enter":
				if len(m.table.Rows()) > 0 {
					// Get the todo title from the selected row
					selectedTitle := m.table.SelectedRow()[1]

					// Find the todo by title
					for _, todo := range m.todoList.Todos {
						if todo.Title == selectedTitle {
							m.mode = ModeViewDetail
							m.viewTaskID = todo.ID
							break
						}
					}
				}
			case "t":
				if len(m.table.Rows()) > 0 {
					if len(m.selectedTodoIDs) > 0 && m.bulkActionActive {
						// Toggle all selected todos
						for id := range m.selectedTodoIDs {
							m.todoList.Toggle(id)
						}
					} else {
						// Toggle by title
						selectedTitle := m.table.SelectedRow()[1]
						for _, todo := range m.todoList.Todos {
							if todo.Title == selectedTitle {
								m.todoList.Toggle(todo.ID)
								break
							}
						}
					}
					m.updateRows()
				}
			case "a":
				// Archive action (only for unarchived tasks)
				if len(m.table.Rows()) > 0 {
					if len(m.selectedTodoIDs) > 0 && m.bulkActionActive {
						// Count how many selected tasks are not archived
						activeCount := 0
						for id := range m.selectedTodoIDs {
							todo := m.findTodoByID(id)
							if todo != nil && !todo.Archived {
								activeCount++
							}
						}

						if activeCount > 0 {
							m.mode = ModeArchiveConfirm
							m.confirmAction = "archive"
						}
					} else {
						// Archive single task
						selectedTitle := m.table.SelectedRow()[1]
						for _, todo := range m.todoList.Todos {
							if todo.Title == selectedTitle && !todo.Archived {
								m.mode = ModeArchiveConfirm
								m.confirmAction = "archive"
								m.actionTitle = selectedTitle
								break
							}
						}
					}
				}
			case "u":
				// Unarchive action (only for archived tasks)
				if len(m.table.Rows()) > 0 {
					if len(m.selectedTodoIDs) > 0 && m.bulkActionActive {
						// Count how many selected tasks are archived
						archivedCount := 0
						for id := range m.selectedTodoIDs {
							todo := m.findTodoByID(id)
							if todo != nil && todo.Archived {
								archivedCount++
							}
						}

						if archivedCount > 0 {
							// Unarchive all selected archived todos
							for id := range m.selectedTodoIDs {
								todo := m.findTodoByID(id)
								if todo != nil && todo.Archived {
									m.todoList.Unarchive(id)
								}
							}
							m.updateRows()
						}
					} else {
						// Unarchive single task
						selectedTitle := m.table.SelectedRow()[1]
						for _, todo := range m.todoList.Todos {
							if todo.Title == selectedTitle && todo.Archived {
								m.todoList.Unarchive(todo.ID)
								m.updateRows()
								break
							}
						}
					}
				}
			case "d":
				if len(m.table.Rows()) > 0 {
					if len(m.selectedTodoIDs) > 0 && m.bulkActionActive {
						// Bulk delete - show confirmation
						m.mode = ModeDeleteConfirm
						m.confirmAction = "delete"
					} else {
						// Delete by title
						selectedTitle := m.table.SelectedRow()[1]
						m.mode = ModeDeleteConfirm
						m.confirmAction = "delete"
						m.actionTitle = selectedTitle
					}
				}
			case " ":
				if len(m.table.Rows()) > 0 {
					// Get todo ID by title from the selected row
					selectedTitle := m.table.SelectedRow()[1]

					// Find the todo by title to get its ID
					for _, todo := range m.todoList.Todos {
						if todo.Title == selectedTitle {
							// Toggle selection for this todo
							if m.selectedTodoIDs[todo.ID] {
								delete(m.selectedTodoIDs, todo.ID)
							} else {
								m.selectedTodoIDs[todo.ID] = true
							}
							break
						}
					}

					// Update bulk action flag
					m.bulkActionActive = len(m.selectedTodoIDs) > 0

					m.updateRows()
					return m, nil // Don't pass space to table to avoid moving cursor
				}
			case "n":
				// Add new task
				m.mode = ModeAddTask
				m.textInput.Focus()
				return m, textinput.Blink
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m TodoTableModel) View() string {
	// Task detail view
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

	// Confirmation dialogs
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

	// Add task view
	if m.mode == ModeAddTask {
		inputView := inputStyle.Render(
			inputPromptStyle.Render("Add New Task") + "\n\n" +
				m.textInput.View() + "\n\n" +
				helpStyle.Render("Press Enter to save, Esc to cancel"))

		return fullScreenStyle.Width(m.width).Height(m.height).Render(inputView)
	}

	// Empty todo list
	if len(m.todoList.Todos) == 0 {
		return baseStyle.Render("No tasks found. Press 'n' to add a new task!")
	}

	// Regular view - dynamic help text based on selection state
	var helpText string

	// Add title to indicate archived/active view
	listTitle := "Active Tasks"
	if m.showArchived {
		listTitle = "Archived Tasks"
	}

	if m.bulkActionActive {
		helpText = "\n" + listTitle + " - Bulk Mode:" +
			"\n→ t: toggle completion for all selected" +
			"\n→ a: archive selected (active tasks)" +
			"\n→ u: unarchive selected (archived tasks)" +
			"\n→ d: delete selected" +
			"\n→ space: toggle selection" +
			"\n→ enter: view details" +
			"\n→ n: add new task" +
			"\n→ q: quit"
	} else {
		helpText = "\n" + listTitle + ":" +
			"\n→ t: toggle completion" +
			"\n→ a: archive (active tasks)" +
			"\n→ u: unarchive (archived tasks)" +
			"\n→ d: delete" +
			"\n→ space: select" +
			"\n→ enter: view details" +
			"\n→ n: add new task" +
			"\n→ q: quit"
	}

	help := helpStyle.Render(helpText)
	return baseStyle.Render(m.table.View()) + help
}
