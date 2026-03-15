package tui

import "charm.land/lipgloss/v2"

var (
	// Title style
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("white"))

	// Cursor indicator
	Cursor = lipgloss.NewStyle().
		Foreground(lipgloss.Blue).
		Bold(true)

	// Selected item - blue foreground
	SelectedItem = lipgloss.NewStyle().
		Foreground(lipgloss.Blue).
		Bold(true)

	// Normal item - gray foreground
	NormalItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("242"))

	// Active field label
	FieldActive = lipgloss.NewStyle().
		Foreground(lipgloss.Blue).
		Bold(true)

	// Inactive field label
	FieldInactive = lipgloss.NewStyle().
		Foreground(lipgloss.Color("242"))

	// Field labels (static text)
	FieldLabel = lipgloss.NewStyle().
		Foreground(lipgloss.Color("white"))

	// Help text - muted gray
	HelpText = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	// Error messages - red
	ErrorMsg = lipgloss.NewStyle().
		Foreground(lipgloss.Red).
		Bold(true)
)
