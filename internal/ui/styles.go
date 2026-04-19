package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Base styles
	BaseStyle = lipgloss.NewStyle().Margin(1, 2)

	// Typography
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63")). // Purple
			Bold(true).
			MarginBottom(1)

	SubTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")). // Lighter Purple
			MarginBottom(1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")). // Gray
			MarginTop(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")). // Red
			MarginTop(1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")). // Green
			MarginTop(1)

	// Inputs & Lists
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	// List item styles
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("205"))
)
