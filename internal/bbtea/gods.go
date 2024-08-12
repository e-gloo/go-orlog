package bbtea

import (
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	c "github.com/e-gloo/orlog/internal/client"
)

var unselectableGodStyle = lipgloss.NewStyle().Faint(true)

type godModel struct {
	client    c.Client
	isMe      bool
	choices   [3]int
	cursor    int
	validated bool
}

func initialGodModel(client c.Client, isMe bool) tea.Model {
	gm := godModel{
		client: client,
		isMe:   isMe,
		cursor: -1,
	}
	if isMe {
		gm.choices = client.GetMyGods()
	} else {
		gm.choices = client.GetOpponentGods()
	}
	return gm
}

func (gm godModel) Init() tea.Cmd {
	return nil
}

func (gm godModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch ph.Phase() {
		case c.SelectGod:
			switch msg.Type {
			case tea.KeyUp:
				if gm.cursor > -1 {
					gm.cursor--
				} else {
					gm.cursor = 8
				}

			// Less than 8 because 3 levels per god and we have 3 gods minus 1
			case tea.KeyDown:
				if gm.cursor < 8 {
					gm.cursor++
				} else {
					gm.cursor = -1
				}
			case tea.KeyEnter:
				if gm.cursor == -1 {
					gm.validated = true
					gm.client.PlayGod(-1, -1)
				} else {
					godIdx := int(math.Floor(float64(gm.cursor) / 3))
					levelIdx := gm.cursor % 3
					if gm.canPick(godIdx, levelIdx) {
						gm.validated = true
						gm.client.PlayGod(godIdx, levelIdx)
					}
				}
			}
		}
	}
	return gm, cmd
}

func (gm godModel) View() string {
	var s string

	if ph.Phase() == c.SelectGod && gm.isMe {
		if gm.validated {
			return fmt.Sprintf("Selection validated\n")
		}
		idx := -1

		cursor := " "
		if gm.cursor == idx {
			cursor = ">"
		}
		s += fmt.Sprintf("%s None\n", cursor)
		idx++

		allGods := gm.client.GetGameGods()

		for _, godIdx := range gm.choices {
			god := allGods[godIdx]
			s += fmt.Sprintf("  %s %s: %s\n", god.Name, god.Emoji, god.Description)
			for levelIdx, level := range god.Levels {
				cursor = " "
				if gm.cursor == idx {
					cursor = ">"
				}
				content := fmt.Sprintf("\t%s %d🪙 %s", cursor, level.TokenCost, level.Description)
				if gm.canPick(godIdx, levelIdx) {
					s += content
				} else {
					cursor = "x"
					s += unselectableGodStyle.Render(content)
				}
				s += "\n"
				idx++
			}
		}
	} else if ph.Phase() == c.WaitingGodSelection && !gm.isMe {
		s += "Waiting for other player to select their god\n"
	}

	return s
}

func (gm godModel) canPick(god int, level int) bool {
	cost := gm.client.GetGameGods()[god].Levels[level].TokenCost
	return cost <= gm.client.GetMe().GetTokens()
}
