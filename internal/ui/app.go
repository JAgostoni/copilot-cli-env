package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jagostoni/copilot-cli-env/internal/config"
	"github.com/jagostoni/copilot-cli-env/internal/provider"
)

type Step int

const (
	StepEnv Step = iota
	StepProvider
	StepAuth
	StepModel
	StepOutput
	StepSummary
	StepWarning
	StepDone
)

type AppModel struct {
	step       Step
	configData config.CopilotEnv
	
	// Temporary state
	selectedProvider provider.Provider
	detectedShell    string
	detectedOS       string
	outputMode       OutputMode
	cancelled        bool

	// Sub-models
	envModel      envModel
	providerModel providerModel
	authModel     authModel
	modelPicker   modelPickerModel
	outputModel   outputModel
	summaryModel  summaryModel
	warningModel  warningModel

	err error
}

func InitialModel() *AppModel {
	m := &AppModel{
		step: StepEnv,
	}
	m.envModel = initialEnvModel(m)
	m.providerModel = initialProviderModel(m)
	m.authModel = initialAuthModel(m)
	m.modelPicker = initialModelPickerModel(m)
	m.outputModel = initialOutputModel(m)
	m.summaryModel = initialSummaryModel(m)
	m.warningModel = initialWarningModel(m)
	return m
}

func (m *AppModel) Init() tea.Cmd {
	return m.envModel.Init()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.cancelled = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.step {
	case StepEnv:
		m.envModel, cmd = m.envModel.Update(msg)
	case StepProvider:
		m.providerModel, cmd = m.providerModel.Update(msg)
	case StepAuth:
		m.authModel, cmd = m.authModel.Update(msg)
	case StepModel:
		m.modelPicker, cmd = m.modelPicker.Update(msg)
	case StepOutput:
		m.outputModel, cmd = m.outputModel.Update(msg)
	case StepSummary:
		m.summaryModel, cmd = m.summaryModel.Update(msg)
	case StepWarning:
		m.warningModel, cmd = m.warningModel.Update(msg)
	case StepDone:
		return m, tea.Quit
	}

	return m, cmd
}

func (m *AppModel) View() string {
	if m.err != nil {
		return BaseStyle.Render(TitleStyle.Render("Error") + "\n" + ErrorStyle.Render(m.err.Error()) + "\nPress Ctrl+C to exit.")
	}
	
	if m.cancelled {
		return BaseStyle.Render(ErrorStyle.Render("Configuration Cancelled."))
	}

	var content string
	switch m.step {
	case StepEnv:
		content = m.envModel.View()
	case StepProvider:
		content = m.providerModel.View()
	case StepAuth:
		content = m.authModel.View()
	case StepModel:
		content = m.modelPicker.View()
	case StepOutput:
		content = m.outputModel.View()
	case StepSummary:
		content = m.summaryModel.View()
	case StepWarning:
		content = m.warningModel.View()
	case StepDone:
		content = SuccessStyle.Render("Configuration Complete!") + "\n" + HelpStyle.Render("Proceeding to apply phase...")
	}

	return BaseStyle.Render(content)
}

func (m *AppModel) nextStep() tea.Cmd {
	m.step++
	switch m.step {
	case StepProvider:
		return m.providerModel.Init()
	case StepAuth:
		if m.selectedProvider.ID() == "reset" {
			// Skip auth and model pickers directly to output selection
			env, _ := m.selectedProvider.MapEnv("", "", "")
			m.configData = env
			m.step = StepOutput
			return m.outputModel.Init()
		}

		m.authModel.needsKey = m.selectedProvider.RequiresAPIKey()
		m.authModel.needsURL = m.selectedProvider.DefaultBaseURL() == ""
		if m.authModel.needsURL {
			m.authModel.focusIndex = 0
			m.authModel.urlInput.Focus()
			m.authModel.keyInput.Blur()
		} else if m.authModel.needsKey {
			m.authModel.focusIndex = 1
			m.authModel.keyInput.Focus()
			m.authModel.urlInput.Blur()
		}
		return m.authModel.Init()
	case StepModel:
		return m.modelPicker.Init()
	case StepOutput:
		// Map the raw inputs into the final Copilot Env using the selected provider
		env, err := m.selectedProvider.MapEnv(m.configData.ProviderBaseURL, m.configData.Model, m.configData.APIKey)
		if err == nil {
			env.MaxPromptTokens = m.configData.MaxPromptTokens
			env.MaxOutputTokens = m.configData.MaxOutputTokens
			m.configData = env
		}
		return m.outputModel.Init()
	case StepSummary:
		return m.summaryModel.Init()
	case StepWarning:
		if (m.outputMode == ModeProfile || m.outputMode == ModeEnvFile) && m.configData.APIKey != "" && m.configData.APIKey != "dummy" {
			return m.warningModel.Init()
		}
		m.step++ // Skip warning if not writing secret to disk
		fallthrough
	case StepDone:
		return tea.Quit
	}
	return nil
}

func (m *AppModel) ConfigData() config.CopilotEnv {
	return m.configData
}

func (m *AppModel) OutputMode() string {
	return string(m.outputMode)
}

func (m *AppModel) DetectedShell() string {
	return m.detectedShell
}

func (m *AppModel) DetectedOS() string {
	return m.detectedOS
}

func (m *AppModel) IsCancelled() bool {
	return m.cancelled
}
