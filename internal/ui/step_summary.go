package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jagostoni/copilot-cli-env/internal/render"
)

type summaryModel struct {
	app *AppModel
}

func initialSummaryModel(app *AppModel) summaryModel {
	return summaryModel{
		app: app,
	}
}

func (m summaryModel) Init() tea.Cmd {
	return nil
}

func (m summaryModel) Update(msg tea.Msg) (summaryModel, tea.Cmd) {
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

func (m summaryModel) View() string {
	content := TitleStyle.Render("Configuration Summary") + "\n\n"

	content += fmt.Sprintf("Target OS:    %s\n", m.app.detectedOS)
	content += fmt.Sprintf("Target Shell: %s\n", m.app.detectedShell)
	content += fmt.Sprintf("Provider:     %s (%s)\n", m.app.selectedProvider.Name(), m.app.configData.ProviderType)
	
	baseURL := m.app.configData.ProviderBaseURL
	if baseURL == "" {
		baseURL = m.app.selectedProvider.DefaultBaseURL()
	}
	content += fmt.Sprintf("Base URL:     %s\n", baseURL)
	content += fmt.Sprintf("Model:        %s\n", m.app.configData.Model)
	
	if m.app.configData.APIKey != "" {
		content += fmt.Sprintf("API Key:      %s\n", render.MaskSecret(m.app.configData.APIKey))
	} else if !m.app.selectedProvider.RequiresAPIKey() {
		content += "API Key:      Not required\n"
	}

	content += "\n"

	var action string
	switch m.app.outputMode {
	case ModeProfile:
		action = fmt.Sprintf("Will update your %s profile securely.", m.app.detectedShell)
	case ModeEnvFile:
		action = "Will generate a .env file."
	case ModeConsole:
		action = "Will output commands to the console."
	}

	content += SubTitleStyle.Render("Action: ") + action + "\n\n"
	content += "Proceed? [y/N]"

	return content
}
