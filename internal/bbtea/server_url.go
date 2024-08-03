package bbtea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	c "github.com/e-gloo/orlog/internal/client"
)

type serverUrlModel struct {
	textInput textinput.Model
	validated bool
	err       error
}

type ClientConnection struct {
	client c.Client
}

func initialServerUrlModel() serverUrlModel {
	ti := textinput.New()
	ti.Placeholder = "localhost"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return serverUrlModel{
		textInput: ti,
		validated: false,
	}
}

func (su serverUrlModel) Init() tea.Cmd {
	return textinput.Blink
}

func (su serverUrlModel) Update(msg tea.Msg) (serverUrlModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if su.textInput.Value() == "" {
				su.textInput.SetValue("localhost:8080")
			}
			client, err := setClient(su.textInput.Value())
			if err != nil {
				su.err = err
				su.textInput.SetValue("")
				return su, cmd
			}
			su.validated = true
			cc := ClientConnection{
				client: client,
			}
			cmd = setCmd(cc)
		default:
			su.textInput, cmd = su.textInput.Update(msg)
		}
	}
	return su, cmd
}

func (su serverUrlModel) View() string {
	var s string
	if su.err != nil {
		s = su.err.Error() + "\n"
	}
	if su.validated {
		s += fmt.Sprintf("You're connected to %s\n", su.textInput.Value())
	} else {
		s += fmt.Sprintf(
			"What is the server url (blank for localhost)?\n\n%s\n",
			su.textInput.View(),
		)
	}
	return s
}

func setClient(serverAddr string) (c.Client, error) {
	client, err := c.NewClient(serverAddr)
	if err != nil {
		return nil, fmt.Errorf("Error creating client: %w", err)
	}

	return client, nil
}
