package bbtea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type addPlayerNameModel struct {
	textInput textinput.Model
	validated bool
	err       error
}

func initialAddPlayerNameModel() addPlayerNameModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return addPlayerNameModel{
		textInput: ti,
		validated: false,
		err:       nil,
	}
}

func (su addPlayerNameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (su addPlayerNameModel) Update(msg tea.KeyMsg) (addPlayerNameModel, tea.Cmd) {
	var cmd tea.Cmd

	if msg.Type == tea.KeyEnter {
		if su.textInput.Value() != "" {
			su.validated = true
			cmd = setPhaseCmd(GameStarting)
		}
	} else {
		su.textInput, cmd = su.textInput.Update(msg)
	}
	return su, cmd
}

func (su addPlayerNameModel) View() string {
	var s string
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
