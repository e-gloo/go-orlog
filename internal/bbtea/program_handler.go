package bbtea

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/e-gloo/orlog/internal/client"
)

type programHandler struct {
	p     *tea.Program
	state client.State
	phase client.Phase
}

func (ph *programHandler) Send(msg interface{}) {
	ph.p.Send(msg)
}

func (ph *programHandler) State() client.State {
	return ph.state
}

func (ph *programHandler) SetState(state client.State) {
	ph.state = state
}

func (ph *programHandler) Phase() client.Phase {
	return ph.phase
}

func (ph *programHandler) SetPhase(phase client.Phase) {
	ph.phase = phase
}
