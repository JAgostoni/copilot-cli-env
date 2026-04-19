package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type warningModel struct {
	app *AppModel
}

func initialWarningModel(app *AppModel) warningModel {
	return warningModel{
		app: app,
	}
}

func (m warningModel) Init() tea.Cmd {
	return nil
}

func (m warningModel) Update(msg tea.Msg) (warningModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			return m, m.app.nextStep()
		case "n", "N", "enter":
			m.app.cancelled = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m warningModel) View() string {
	content := TitleStyle.Render("Security Warning") + "\n\n"
	
	warningBox := ErrorStyle.Copy().BorderForeground(ErrorStyle.GetForeground()).BorderStyle(lipgloss.RoundedBorder()).Padding(1, 2)
	
	warningMsg := "WARNING: You are about to write a plaintext API key to disk."
	
	content += warningBox.Render(warningMsg) + "\n\n"
	
	content += "Are you absolutely sure you want to proceed? [y/N]"

	return content
}
