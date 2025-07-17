package ui

import "github.com/charmbracelet/lipgloss"

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Width(80)
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
	successMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("125")).
				Bold(true)
	errorMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("160")).
				Bold(true)
	titleBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Bold(true)
	statusBarContainerStyle = lipgloss.NewStyle().
				Width(100).
				Padding(0, 0)
	tableContainerStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")).
				Padding(1, 2)
)
