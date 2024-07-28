package bbtea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	c "github.com/e-gloo/orlog/internal/client"
)

type configPlayerModel struct {
	client    c.Client
	textInput textinput.Model
	validated bool
}

func initialConfigPlayerModel(client c.Client) configPlayerModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return configPlayerModel{
		client:    client,
		textInput: ti,
		validated: false,
	}
}

func (su configPlayerModel) Init() tea.Cmd {
	return textinput.Blink
}

func (su configPlayerModel) Update(msg tea.Msg) (configPlayerModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return su, tea.Quit
		case tea.KeyEnter:
			if su.textInput.Value() != "" {
				su.validated = true
			}
		default:
			su.textInput, cmd = su.textInput.Update(msg)
		}
	}

	return su, cmd
}

func (su configPlayerModel) View() string {
	var s string

	if ph.Phase() == c.ConfigPlayer && su.client.Error() != "" {
		s = su.client.Error() + "\n"
	}

	if su.validated {
		s = fmt.Sprintf("Get ready %s!\n", su.textInput.Value())
	} else {
		s = fmt.Sprintf(
			"Enter your player name:\n\n %s\n\n%s",
			su.textInput.View(),
			"(esc to quit)",
		) + "\n"
	}
	return s
}
