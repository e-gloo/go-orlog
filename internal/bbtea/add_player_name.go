package bbtea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	c "github.com/e-gloo/orlog/internal/client"
	l "github.com/e-gloo/orlog/internal/client/lobby"
)

type addPlayerNameModel struct {
	client    c.Client
	textInput textinput.Model
	validated bool
}

func initialAddPlayerNameModel(client c.Client) addPlayerNameModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return addPlayerNameModel{
		client:    client,
		textInput: ti,
		validated: false,
	}
}

func (su addPlayerNameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (su addPlayerNameModel) Update(msg tea.Msg) (addPlayerNameModel, tea.Cmd) {
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

func (su addPlayerNameModel) View() string {
	var s string

	lobby := su.client.GetLobby()
	if lobby.Phase == l.AddPlayerName && lobby.Err != "" {
		s = lobby.Err + "\n"
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
