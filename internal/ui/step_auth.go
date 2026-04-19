package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

type authModel struct {
	app *AppModel
	
	keyInput textinput.Model
	urlInput textinput.Model

	focusIndex int
	needsURL   bool
	needsKey   bool
}

func initialAuthModel(app *AppModel) authModel {
	ti := textinput.New()
	ti.Placeholder = "sk-..."
	ti.Focus()
	ti.CharLimit = 1024
	ti.Width = 40
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '*'

	ui := textinput.New()
	ui.Placeholder = "https://api.example.com"
	ui.CharLimit = 256
	ui.Width = 40

	return authModel{
		app:        app,
		keyInput:   ti,
		urlInput:   ui,
		focusIndex: 0,
	}
}

func (m authModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m authModel) Update(msg tea.Msg) (authModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// If both are needed, move focus
			if m.needsURL && m.needsKey && m.focusIndex == 0 {
				m.focusIndex = 1
				m.urlInput.Blur()
				m.keyInput.Focus()
				return m, nil
			}

			// Save values and proceed
			if m.needsKey {
				m.app.configData.APIKey = m.keyInput.Value()
			}
			if m.needsURL {
				m.app.configData.ProviderBaseURL = m.urlInput.Value()
			}
			return m, m.app.nextStep()

		case "tab", "shift+tab", "up", "down":
			if m.needsURL && m.needsKey {
				if m.focusIndex == 0 {
					m.focusIndex = 1
					m.urlInput.Blur()
					m.keyInput.Focus()
				} else {
					m.focusIndex = 0
					m.keyInput.Blur()
					m.urlInput.Focus()
				}
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	if m.focusIndex == 0 && m.needsURL {
		m.urlInput, cmd = m.urlInput.Update(msg)
	} else if m.needsKey {
		m.keyInput, cmd = m.keyInput.Update(msg)
	}

	return m, cmd
}

func (m authModel) View() string {
	var content string
	content += TitleStyle.Render("Authentication for " + m.app.selectedProvider.Name()) + "\n\n"

	if m.needsURL {
		content += "Base URL:\n"
		content += m.urlInput.View() + "\n\n"
	}

	if m.needsKey {
		content += "API Key:\n"
		content += m.keyInput.View() + "\n\n"
	}

	if !m.needsKey && !m.needsURL {
		content += "No authentication required for this provider.\n\n"
		content += HelpStyle.Render("Press [Enter] to continue.")
		return content
	}

	content += HelpStyle.Render("Press [Enter] to continue or [Tab] to switch fields.")
	return content
}
