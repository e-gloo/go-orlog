package bbtea

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Phase int

const (
	ServerConnection Phase = iota
	CreateOrJoinGame
	AddPlayerName
	GameStarting
)

func NewClient() *tea.Program {
	p := tea.NewProgram(initialModel())
	return p
}

type (
	errMsg error
)

type model struct {
	serverUrl     serverUrlModel
	createOrJoin  createOrJoinModel
	addPlayerName addPlayerNameModel
	phase         Phase
	err           error
}

func initialModel() model {
	su := initialServerUrlModel()
	return model{
		serverUrl: su,
		phase:     ServerConnection,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return m.serverUrl.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
		switch m.phase {
		case ServerConnection:
			m.serverUrl, cmd = m.serverUrl.Update(msg)
		case CreateOrJoinGame:
			m.createOrJoin, cmd = m.createOrJoin.Update(msg)
		case AddPlayerName:
			m.addPlayerName, cmd = m.addPlayerName.Update(msg)
		}

	case Phase:
		switch msg {
		case CreateOrJoinGame:
			m.createOrJoin = initialCreateOrJoinModel()
			m.phase = CreateOrJoinGame
		case AddPlayerName:
			m.phase = AddPlayerName
			m.addPlayerName = initialAddPlayerNameModel()
			cmd = m.addPlayerName.Init()
		case GameStarting:
			m.phase = GameStarting
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

func (m model) View() string {
	var s string
	if m.phase >= ServerConnection {
		s += m.serverUrl.View()
	}
	if m.phase >= CreateOrJoinGame {
		s += m.createOrJoin.View()
	}
	if m.phase >= AddPlayerName {
		s += m.addPlayerName.View()
	}
	if m.phase >= GameStarting {
		s += fmt.Sprintf("Game is about to start...\n\n%s\n",
			"(esc to quit)")
	}
	return s
}

func setPhaseCmd(phase Phase) tea.Cmd {
	return func() tea.Msg {
		return Phase(phase)
	}
}
