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
		var status string
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

	sourceText := ""
	if m.sourceLabel != "" {
		sourceText = "  |  source: " + m.sourceLabel
		if m.sourceLabel == "project" && m.projectName != "" {
			sourceText += " (" + m.projectName + ")"
		}
	}
	leftSide := titleBarStyle.Render(listTitle + sourceText)
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
			"\n→ " + confirmBtnStyle.Render("t") + ": toggle completion for all selected" +
			"\n→ " + confirmBtnStyle.Render("n") + ": toggle archive/unarchive for selected" +
			"\n→ " + confirmBtnStyle.Render("d") + ": delete selected" +
			"\n→ " + confirmBtnStyle.Render("space") + ": toggle selection" +
			"\n→ " + confirmBtnStyle.Render("enter") + ": view details" +
			"\n→ " + confirmBtnStyle.Render("a") + ": add new task" +
			"\n→ " + confirmBtnStyle.Render("s") + ": switch source (project/global)" +
			"\n→ " + confirmBtnStyle.Render("q") + ": quit" +
			"\n→ " + confirmBtnStyle.Render(".") + ": toggle help"
	} else {
		helpText = "\n" + statusBar + "\n" +
			"→ " + confirmBtnStyle.Render("t") + ": toggle completion" +
			"\n→ " + confirmBtnStyle.Render("n") + ": toggle archive/unarchive" +
			"\n→ " + confirmBtnStyle.Render("d") + ": delete" +
			"\n→ " + confirmBtnStyle.Render("space") + ": select" +
			"\n→ " + confirmBtnStyle.Render("enter") + ": view details" +
			"\n→ " + confirmBtnStyle.Render("a") + ": add new task" +
			"\n→ " + confirmBtnStyle.Render("s") + ": switch source (project/global)" +
			"\n→ " + confirmBtnStyle.Render("q") + ": quit" +
			"\n→ " + confirmBtnStyle.Render(".") + ": toggle help"
	}

	tableView := tableContainerStyle.Render(m.table.View())
	if m.mode == ModeNormal {
		if m.showHelp {
			help := helpStyle.Render(helpText)
			return tableView + help
		}
		hint := helpStyle.Render("\n→ " + confirmBtnStyle.Render(".") + ": toggle help")
		return tableView + hint
	}
	return tableView
}
