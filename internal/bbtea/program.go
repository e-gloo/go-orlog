package bbtea

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	c "github.com/e-gloo/orlog/internal/client"
	l "github.com/e-gloo/orlog/internal/client/lobby"
)

type model struct {
	client        c.Client
	serverUrl     serverUrlModel
	createOrJoin  createOrJoinModel
	addPlayerName addPlayerNameModel
}

type errMsg error

type Cmd interface{}

var ph *programHandler

func NewClient() *tea.Program {
	p := tea.NewProgram(initialModel())
	ph = &programHandler{p: p}
	return p
}

func initialModel() model {
	su := initialServerUrlModel()
	return model{
		serverUrl: su,
	}
}

func (m model) Init() tea.Cmd {
	return m.serverUrl.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// handle server url input before creating client with ws connection
	if m.client == nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			m.serverUrl, cmd = m.serverUrl.Update(msg)

		case ClientConnection:
			m.client = msg.client
			go m.client.Run(ph)
		}
		return m, cmd
	}

	switch m.client.GetState() {
	case c.LobbyState:
		var newModel tea.Model
		newModel, cmd = m.handleUpdateLobbyState(msg)
		m = newModel.(model)
	}

	return m, cmd
}

func (m model) handleUpdateLobbyState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case l.Phase:
		switch msg {
		case l.CreateOrJoin:
			m.createOrJoin = initialCreateOrJoinModel(m.client)
			return m, m.createOrJoin.Init()
		case l.AddPlayerName:
			m.addPlayerName = initialAddPlayerNameModel(m.client)
			return m, m.addPlayerName.Init()
		}
	default:
		switch m.client.GetLobby().Phase {
		case l.CreateOrJoin:
			m.createOrJoin, cmd = m.createOrJoin.Update(msg)
		case l.AddPlayerName:
			m.addPlayerName, cmd = m.addPlayerName.Update(msg)
		}
	}
	return m, cmd
}

func (m model) View() string {
	if m.client == nil {
		return m.serverUrl.View()
	}

	var s string
	switch m.client.GetState() {
	case c.LobbyState:
		s += m.handleViewLobbyState()
	case c.GameState:
		s += m.handleViewGameState()
	}
	return s
}

func (m model) handleViewLobbyState() string {
	s := m.serverUrl.View()
	phase := m.client.GetLobby().Phase

	if phase >= l.CreateOrJoin {
		s += m.createOrJoin.View()
	}
	if phase >= l.AddPlayerName {
		s += m.addPlayerName.View()
	}
	return s
}

func (m model) handleViewGameState() string {
	return fmt.Sprintf("Game is starting, get ready!\n\n%s\n", "(esc to quit)")
}

func setCmd(cmd Cmd) tea.Cmd {
	return func() tea.Msg {
		return cmd
	}
}
