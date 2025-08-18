package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prime-run/togo/model"
)

func (m TodoTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = msg.Width
		m.height = msg.Height
		m.updateRows()
		return m, nil
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
						count := len(m.selectedTodoIDs)
						for id := range m.selectedTodoIDs {
							m.todoList.Delete(id)
						}
						m.selectedTodoIDs = make(map[int]bool)
						m.bulkActionActive = false
						m.SetStatusMessage(fmt.Sprintf("%d tasks deleted", count))
					} else {
						found := false
						for i, todo := range m.todoList.Todos {
							if todo.Title == m.actionTitle || strings.Contains(m.actionTitle, todo.Title) {
								m.todoList.Todos = append(m.todoList.Todos[:i], m.todoList.Todos[i+1:]...)
								found = true
								m.SetStatusMessage("Task deleted")
								break
							}
						}
						if !found {
							id, err := strconv.Atoi(m.actionTitle)
							if err == nil {
								for _, todo := range m.todoList.Todos {
									if todo.ID == id {
										m.todoList.Delete(id)
										m.SetStatusMessage("Task deleted")
										break
									}
								}
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
				return m, m.forceRelayoutCmd()
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
					m.SetStatusMessage("New task added")
				}
				m.mode = ModeNormal
				return m, m.forceRelayoutCmd()
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
			case ".":
				m.showHelp = !m.showHelp
				m.updateRows()
				return m, m.forceRelayoutCmd()
			case "esc", "q":
				return m, tea.Quit
			case "enter":
				if len(m.table.Rows()) > 0 {
					selectedTitle := m.table.SelectedRow()[1]
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
						count := 0
						for id := range m.selectedTodoIDs {
							todo := m.findTodoByID(id)
							if todo != nil {
								if todo.Archived {
									m.todoList.Unarchive(id)
									count++
								} else {
									m.todoList.Toggle(id)
									count++
								}
							}
						}
						if count > 0 {
							m.SetStatusMessage(fmt.Sprintf("%d tasks updated", count))
						}
					} else {
						selectedTitle := m.table.SelectedRow()[1]
						cleanTitle := strings.Replace(selectedTitle, archivedStyle.Render(""), "", -1)

						for _, todo := range m.todoList.Todos {
							if strings.Contains(selectedTitle, todo.Title) || todo.Title == cleanTitle {
								if todo.Archived {
									m.todoList.Unarchive(todo.ID)
									m.SetStatusMessage("Task unarchived")
								} else {
									m.todoList.Toggle(todo.ID)
									m.SetStatusMessage("Task updated")
								}
								break
							}
						}
					}
					m.updateRows()
					return m, m.forceRelayoutCmd()
				}
			case "n":
				if len(m.table.Rows()) > 0 {
					if len(m.selectedTodoIDs) > 0 && m.bulkActionActive {
						count := 0
						for id := range m.selectedTodoIDs {
							todo := m.findTodoByID(id)
							if todo != nil {
								if todo.Archived {
									m.todoList.Unarchive(id)
									count++
								} else {
									m.todoList.Archive(id)
									count++
								}
							}
						}
						if count > 0 {
							m.SetStatusMessage(fmt.Sprintf("%d tasks updated", count))
						}
						m.updateRows()
						return m, m.forceRelayoutCmd()
					} else {
						selectedTitle := m.table.SelectedRow()[1]
						cleanTitle := strings.Replace(selectedTitle, archivedStyle.Render(""), "", -1)

						for _, todo := range m.todoList.Todos {
							if strings.Contains(selectedTitle, todo.Title) || todo.Title == cleanTitle {
								if todo.Archived {
									m.todoList.Unarchive(todo.ID)
									m.SetStatusMessage("Task unarchived")
								} else {
									m.todoList.Archive(todo.ID)
									m.SetStatusMessage("Task archived")
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
						selectedRow := m.table.SelectedRow()
						if len(selectedRow) > 1 {
							selectedTitle := selectedRow[1]
							cleanTitle := strings.Replace(selectedTitle, archivedStyle.Render(""), "", -1)

							m.mode = ModeDeleteConfirm
							m.confirmAction = "delete"
							m.actionTitle = cleanTitle
						}
					}
				}
			case " ":
				if len(m.table.Rows()) > 0 {
					selectedIndex := m.table.Cursor()
					if selectedIndex >= 0 && selectedIndex < len(m.todoList.Todos) {
						var filteredTodos []model.Todo

						if m.showAll {
							filteredTodos = m.todoList.Todos
						} else if m.showArchivedOnly {
							filteredTodos = m.todoList.GetArchivedTodos()
						} else {
							filteredTodos = m.todoList.GetActiveTodos()
						}

						if selectedIndex < len(filteredTodos) {
							todo := filteredTodos[selectedIndex]
							if m.selectedTodoIDs[todo.ID] {
								delete(m.selectedTodoIDs, todo.ID)
							} else {
								m.selectedTodoIDs[todo.ID] = true
							}
							m.bulkActionActive = len(m.selectedTodoIDs) > 0
							m.updateRows()
							return m, m.forceRelayoutCmd()
						}
					}
					return m, nil
				}
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m TodoTableModel) forceRelayoutCmd() tea.Cmd {
	width, height := m.width, m.height
	return func() tea.Msg {
		return tea.WindowSizeMsg{Width: width, Height: height}
	}
}
