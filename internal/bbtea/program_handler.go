package bbtea

import tea "github.com/charmbracelet/bubbletea"

type programHandler struct {
	p *tea.Program
}

func (ph *programHandler) Send(msg interface{}) {
	ph.p.Send(msg)
}
