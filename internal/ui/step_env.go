package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jagostoni/copilot-cli-env/internal/envdetect"
)

type envModel struct {
	app      *AppModel
	os       string
	shell    string
	terminal string
	loading  bool
}

func initialEnvModel(app *AppModel) envModel {
	return envModel{
		app:     app,
		loading: true,
	}
}

type envDetectedMsg struct {
	env envdetect.Environment
	err error
}

func (m envModel) Init() tea.Cmd {
	return func() tea.Msg {
		detector := envdetect.NewDetector()
		env, err := detector.Detect()
		return envDetectedMsg{env: env, err: err}
	}
}

func (m envModel) Update(msg tea.Msg) (envModel, tea.Cmd) {
	switch msg := msg.(type) {
	case envDetectedMsg:
		m.loading = false
		if msg.err != nil {
			m.app.err = msg.err
			return m, nil
		}
		m.os = msg.env.OS
		m.shell = msg.env.Shell
		m.terminal = msg.env.Terminal
		
		// Save for later
		m.app.detectedShell = msg.env.Shell
		m.app.detectedOS = msg.env.OS
		
	case tea.KeyMsg:
		if !m.loading {
			switch msg.String() {
			case "enter", " ":
				return m, m.app.nextStep()
			}
		}
	}
	return m, nil
}

func (m envModel) View() string {
	if m.loading {
		return "Detecting environment..."
	}

	content := TitleStyle.Render("Environment Detected") + "\n"
	content += fmt.Sprintf("OS:       %s\n", m.os)
	content += fmt.Sprintf("Shell:    %s\n", m.shell)
	if m.terminal != "" {
		content += fmt.Sprintf("Terminal: %s\n", m.terminal)
	}

	content += HelpStyle.Render("\nPress [Enter] to continue or [Ctrl+C] to quit.")
	return content
}
