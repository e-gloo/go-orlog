package bbtea

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	c "github.com/e-gloo/orlog/internal/client"
)

type model struct {
	client       c.Client
	serverUrl    serverUrlModel
	createOrJoin createOrJoinModel
	configPlayer configPlayerModel
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

	switch ph.State() {
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
	case c.Phase:
		ph.SetPhase(msg)
		switch msg {
		case c.CreateOrJoin:
			m.createOrJoin = initialCreateOrJoinModel(m.client)
			return m, m.createOrJoin.Init()
		case c.ConfigPlayer:
			m.configPlayer = initialConfigPlayerModel(m.client)
			return m, m.configPlayer.Init()
		}
	default:
		switch ph.Phase() {
		case c.CreateOrJoin:
			m.createOrJoin, cmd = m.createOrJoin.Update(msg)
		case c.ConfigPlayer:
			m.configPlayer, cmd = m.configPlayer.Update(msg)
		}
	}
	return m, cmd
}

func (m model) View() string {
	if m.client == nil {
		return m.serverUrl.View()
	}

	var s string
	switch ph.State() {
	case c.LobbyState:
		s += m.handleViewLobbyState()
	case c.GameState:
		s += m.handleViewGameState()
	}
	return s
}

func (m model) handleViewLobbyState() string {
	s := m.serverUrl.View()
	phase := ph.Phase()

	if phase >= c.CreateOrJoin {
		s += m.createOrJoin.View()
	}
	if phase >= c.ConfigPlayer {
		s += m.configPlayer.View()
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
