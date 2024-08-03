package bbtea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	c "github.com/e-gloo/orlog/internal/client"
	g "github.com/e-gloo/orlog/internal/client/game"
)

type hudModel struct {
	client c.Client
	isMe   bool
}

var godStyle = lipgloss.NewStyle().
	Bold(true).
	PaddingTop(1)

func initalHudModel(client c.Client, isMe bool) hudModel {
	return hudModel{client: client, isMe: isMe}
}

func (h hudModel) Init() tea.Cmd {
	return nil
}

func (h hudModel) Update(msg tea.Msg) (hudModel, tea.Cmd) {
	var cmd tea.Cmd
	return h, cmd
}

func (h hudModel) View() string {
	var s string
	var player *g.ClientPlayer

	var playerStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(true).
		Italic(true)

	var nameLen int
	if h.isMe {
		nameLen = len(h.client.GetMe().GetUsername())
	} else {
		nameLen = len(h.client.GetOpponent().GetUsername())
	}
	playerStyle.Width(nameLen)

	if h.isMe {
		player = h.client.GetMe()
	} else {
		player = h.client.GetOpponent()
	}

	s = fmt.Sprintf("%s\n", playerStyle.Render(player.GetUsername()))
	s += fmt.Sprintf("🤎: %d\t🪙: %d\n", player.GetHealth(), player.GetTokens())

	gameGods := h.client.GetGameGods()
	var gods []string
	for _, god := range player.GetGods() {
		currentGod := gameGods[god]
		gods = append(gods, godStyle.Render(fmt.Sprintf("\t%s %s", currentGod.Name, currentGod.Emoji)))
	}
	s += fmt.Sprintf("%s\n", lipgloss.JoinHorizontal(lipgloss.Top, gods...))

	return s
}
