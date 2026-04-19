package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type OutputMode string

const (
	ModeProfile OutputMode = "profile"
	ModeEnvFile OutputMode = "env"
	ModeConsole OutputMode = "console"
)

type outputItem struct {
	mode        OutputMode
	title       string
	description string
}

func (i outputItem) Title() string       { return i.title }
func (i outputItem) Description() string { return i.description }
func (i outputItem) FilterValue() string { return i.title }

type outputItemDelegate struct{}

func (d outputItemDelegate) Height() int                             { return 2 }
func (d outputItemDelegate) Spacing() int                            { return 1 }
func (d outputItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d outputItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(outputItem)
	if !ok {
		return
	}

	title := i.title
	desc := i.description

	if index == m.Index() {
		fmt.Fprint(w, SelectedItemStyle.Render("> "+title+"\n    "+desc))
	} else {
		fmt.Fprint(w, ItemStyle.Render("  "+title+"\n    "+desc))
	}
}

type outputModel struct {
	app  *AppModel
	list list.Model
}

func initialOutputModel(app *AppModel) outputModel {
	items := []list.Item{
		outputItem{mode: ModeProfile, title: "Profile Update", description: "Inject into your shell profile (e.g., ~/.bashrc) for persistence."},
		outputItem{mode: ModeEnvFile, title: ".env file", description: "Generate a .env file in the current directory."},
		outputItem{mode: ModeConsole, title: "Console Output", description: "Print the commands to the console to copy/paste manually."},
	}

	l := list.New(items, outputItemDelegate{}, 50, 15)
	l.Title = "Select Output Method"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return outputModel{
		app:  app,
		list: l,
	}
}

func (m outputModel) Init() tea.Cmd {
	return nil
}

func (m outputModel) Update(msg tea.Msg) (outputModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := m.list.SelectedItem().(outputItem)
			if ok {
				m.app.outputMode = i.mode
				return m, m.app.nextStep()
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m outputModel) View() string {
	return m.list.View()
}
