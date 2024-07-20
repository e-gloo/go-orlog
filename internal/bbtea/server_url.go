package bbtea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type serverUrlModel struct {
	textInput textinput.Model
	err       error
}

func initialServerUrlModel() serverUrlModel {
	ti := textinput.New()
	ti.Placeholder = "localhost"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return serverUrlModel{
		textInput: ti,
		err:       nil,
	}
}

func (su serverUrlModel) Init() tea.Cmd {
	return textinput.Blink
}

func (su serverUrlModel) Update(msg tea.KeyMsg) (serverUrlModel, tea.Cmd) {
	var cmd tea.Cmd

	if msg.Type == tea.KeyEnter {
		cmd = setPhaseCmd(CreateOrJoinGame)
	} else {
		su.textInput, cmd = su.textInput.Update(msg)
	}
	return su, cmd
}

func (su serverUrlModel) View() string {
	return fmt.Sprintf(
		"What is the server url (blank for localhost)?\n\n%s\n\n%s",
		su.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
