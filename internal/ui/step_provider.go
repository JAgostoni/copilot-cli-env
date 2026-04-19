package ui

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/jagostoni/copilot-cli-env/internal/provider"
)

type providerItem struct {
	provider provider.Provider
}

func (i providerItem) FilterValue() string { return i.provider.Name() }

type providerItemDelegate struct{}

func (d providerItemDelegate) Height() int                             { return 1 }
func (d providerItemDelegate) Spacing() int                            { return 0 }
func (d providerItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d providerItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(providerItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.provider.Name())

	if index == m.Index() {
		fmt.Fprint(w, SelectedItemStyle.Render("> "+str))
	} else {
		fmt.Fprint(w, ItemStyle.Render("  "+str))
	}
}

type providerModel struct {
	app  *AppModel
	list list.Model
}

func initialProviderModel(app *AppModel) providerModel {
	items := []list.Item{
		providerItem{provider.NewOpenAIProvider()},
		providerItem{provider.NewAnthropicProvider()},
		providerItem{provider.NewOpenRouterProvider()},
		providerItem{provider.NewAzureProvider()},
		providerItem{provider.NewOllamaProvider()},
		providerItem{provider.NewResetProvider()},
	}

	l := list.New(items, providerItemDelegate{}, 30, 14)
	l.Title = "Select AI Provider"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return providerModel{
		app:  app,
		list: l,
	}
}

func (m providerModel) Init() tea.Cmd {
	return nil
}

func (m providerModel) Update(msg tea.Msg) (providerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.list.SelectedItem().(providerItem)
			if ok {
				m.app.selectedProvider = i.provider
				m.app.configData.ProviderType = string(i.provider.ID())
				if i.provider.DefaultBaseURL() != "" {
					m.app.configData.ProviderBaseURL = i.provider.DefaultBaseURL()
				}
				return m, m.app.nextStep()
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m providerModel) View() string {
	return m.list.View()
}
