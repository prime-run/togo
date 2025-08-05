package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/prime-run/togo/model"
)

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

	leftSide := titleBarStyle.Render(listTitle)
	rightSide := successMessageStyle.Render(m.statusMessage)

	statusBar := lipgloss.JoinHorizontal(
		lipgloss.Center,
		leftSide,
		lipgloss.PlaceHorizontal(
			m.width-lipgloss.Width(leftSide)-4,
			lipgloss.Right,
			rightSide,
		),
	)

	if m.bulkActionActive {
		helpText = "\n" + statusBar + "\n" +
			"Bulk Mode:" +
			"\n→ t: toggle completion for all selected" +
			"\n→ n: toggle archive/unarchive for selected" +
			"\n→ d: delete selected" +
			"\n→ space: toggle selection" +
			"\n→ enter: view details" +
			"\n→ a: add new task" +
			"\n→ q: quit"
	} else {
		helpText = "\n" + statusBar + "\n" +
			"→ t: toggle completion" +
			"\n→ n: toggle archive/unarchive" +
			"\n→ d: delete" +
			"\n→ space: select" +
			"\n→ enter: view details" +
			"\n→ a: add new task" +
			"\n→ q: quit"
	}

	help := helpStyle.Render(helpText)

	tableView := tableContainerStyle.Render(m.table.View())
	return tableView + help
}
