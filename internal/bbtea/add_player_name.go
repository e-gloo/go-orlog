package bbtea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type addPlayerNameModel struct {
	textInput textinput.Model
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
		err:       nil,
	}
}

func (su addPlayerNameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (su addPlayerNameModel) Update(msg tea.KeyMsg) (addPlayerNameModel, tea.Cmd) {
	var cmd tea.Cmd

	if msg.Type == tea.KeyEnter {
		cmd = setPhaseCmd(GameStarting)
	} else {
		su.textInput, cmd = su.textInput.Update(msg)
	}
	return su, cmd
}

func (su addPlayerNameModel) View() string {
	return fmt.Sprintf(
		"Enter your player name:\n\n %s\n\n%s",
		su.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

