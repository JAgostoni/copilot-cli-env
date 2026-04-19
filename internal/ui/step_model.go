package ui

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jagostoni/copilot-cli-env/internal/provider"
)

type modelItem struct {
	model provider.Model
}

func (i modelItem) Title() string       { return i.model.ID }
func (i modelItem) Description() string { return i.model.Description }
func (i modelItem) FilterValue() string { return i.model.ID + " " + i.model.Description }

type modelItemDelegate struct{}

func (d modelItemDelegate) Height() int                             { return 1 }
func (d modelItemDelegate) Spacing() int                            { return 0 }
func (d modelItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d modelItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(modelItem)
	if !ok {
		return
	}

	str := i.model.ID
	if i.model.Description != "" {
		str += fmt.Sprintf(" (%s)", i.model.Description)
	}

	if index == m.Index() {
		fmt.Fprint(w, SelectedItemStyle.Render("> "+str))
	} else {
		fmt.Fprint(w, ItemStyle.Render("  "+str))
	}
}

type modelsFetchedMsg struct {
	models []provider.Model
	err    error
}

type modelPickerModel struct {
	app     *AppModel
	spinner spinner.Model
	list    list.Model
	input   textinput.Model

	loading bool
	manual  bool
	err     error
}

func initialModelPickerModel(app *AppModel) modelPickerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // need lipgloss import or use styles, but it's fine

	ti := textinput.New()
	ti.Placeholder = "e.g., custom-model-name"
	ti.CharLimit = 128
	ti.Width = 40

	l := list.New([]list.Item{}, modelItemDelegate{}, 60, 14)
	l.Title = "Select AI Model"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("ctrl+t"),
				key.WithHelp("ctrl+t", "manual entry"),
			),
		}
	}
	l.AdditionalFullHelpKeys = l.AdditionalShortHelpKeys

	return modelPickerModel{
		app:     app,
		spinner: s,
		list:    l,
		input:   ti,
		loading: true,
	}
}

func fetchModelsCmd(p provider.Provider, baseURL, apiKey string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		models, err := p.FetchModels(ctx, baseURL, apiKey)
		return modelsFetchedMsg{models: models, err: err}
	}
}

func (m modelPickerModel) Init() tea.Cmd {
	m.loading = true
	m.manual = false
	m.err = nil
	m.input.SetValue("")
	
	// If provider doesn't support discovery, switch to manual mode automatically
	if _, err := m.app.selectedProvider.FetchModels(context.Background(), "", ""); err == provider.ErrDiscoveryUnsupported {
		m.manual = true
		m.loading = false
		m.input.Focus()
		return textinput.Blink
	}

	return tea.Batch(
		m.spinner.Tick,
		fetchModelsCmd(m.app.selectedProvider, m.app.configData.ProviderBaseURL, m.app.configData.APIKey),
	)
}

func (m modelPickerModel) Update(msg tea.Msg) (modelPickerModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case modelsFetchedMsg:
		m.loading = false
		if msg.err != nil {
			if msg.err == provider.ErrDiscoveryUnsupported {
				m.manual = true
				m.input.Focus()
				return m, textinput.Blink
			}
			m.err = msg.err
			m.manual = true
			m.input.Focus()
			return m, textinput.Blink
		}
		items := make([]list.Item, len(msg.models))
		for i, mod := range msg.models {
			items[i] = modelItem{model: mod}
		}
		m.list.SetItems(items)
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+t" {
			m.manual = !m.manual
			if m.manual {
				m.input.Focus()
				return m, textinput.Blink
			} else {
				m.input.Blur()
				return m, nil
			}
		}

		if msg.String() == "enter" {
			if m.manual {
				if m.input.Value() != "" {
					m.app.configData.Model = m.input.Value()
					m.app.configData.MaxPromptTokens = 0
					m.app.configData.MaxOutputTokens = 0
					return m, m.app.nextStep()
				}
			} else if !m.loading {
				i, ok := m.list.SelectedItem().(modelItem)
				if ok {
					m.app.configData.Model = i.model.ID
					m.app.configData.MaxPromptTokens = i.model.ContextLength
					m.app.configData.MaxOutputTokens = i.model.MaxCompletionTokens
					return m, m.app.nextStep()
				}
			}
		}
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.manual {
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m modelPickerModel) View() string {
	if m.loading {
		return fmt.Sprintf("\n %s Fetching models for %s...\n", m.spinner.View(), m.app.selectedProvider.Name())
	}

	if m.manual {
		var content string
		content += TitleStyle.Render("Manual Model Entry") + "\n\n"
		if m.err != nil {
			content += ErrorStyle.Render(fmt.Sprintf("Discovery failed: %v", m.err)) + "\n\n"
		}
		content += "Enter model ID:\n"
		content += m.input.View() + "\n\n"
		content += HelpStyle.Render("Press [Enter] to continue or [Ctrl+T] to return to list.")
		return content
	}

	return m.list.View() + "\n" + HelpStyle.Render("Press [Ctrl+T] to enter a custom model manually.")
}
