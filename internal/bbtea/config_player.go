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

func initialConfigPlayerModel(client c.Client) tea.Model {
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

func (su configPlayerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if su.textInput.Value() != "" {
				su.validated = true
				su.client.AddPlayerName(su.textInput.Value())
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
		s += fmt.Sprintf(
			"Enter your player name:\n\n %s\n",
			su.textInput.View(),
		)
		return s
	}

	if su.validated {
		s = fmt.Sprintf("Get ready %s, the game is starting soon!\n", su.textInput.Value())
	} else {
		s = fmt.Sprintf(
			"Enter your player name:\n\n %s\n",
			su.textInput.View(),
		)
	}
	return s
}
