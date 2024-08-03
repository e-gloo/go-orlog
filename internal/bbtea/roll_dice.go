package bbtea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	c "github.com/e-gloo/orlog/internal/client"
	g "github.com/e-gloo/orlog/internal/client/game"
)

type diceModel struct {
	client    c.Client
	cursor    int
	nbDie     int
	myDice    bool
	validated bool
}

var baseDieStyle = lipgloss.NewStyle().
	Width(5).
	Align(lipgloss.Center).
	BorderStyle(lipgloss.RoundedBorder())

var baseDieBoxStyle = lipgloss.NewStyle().
	MarginLeft(1).
	Align(lipgloss.Center)

func initialDiceModel(client c.Client, nbDice int, myDice bool) diceModel {
	return diceModel{
		client:    client,
		cursor:    0,
		nbDie:     nbDice,
		myDice:    myDice,
		validated: false,
	}
}

func (dm diceModel) Init() tea.Cmd {
	return nil
}

func (dm diceModel) Update(msg tea.Msg) (diceModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch ph.Phase() {
		case c.RollDice:
			switch msg.Type {
			case tea.KeyEnter:
				dm.validated = true
				dm.client.RollDice()
			}
		case c.PickDice:
			switch msg.Type {
			case tea.KeyLeft:
				if dm.cursor > 0 {
					dm.cursor--
				}

			// The "down" and "j" keys move the cursor down
			case tea.KeyRight:
				if dm.cursor < dm.nbDie-1 {
					dm.cursor++
				}
			case tea.KeySpace:
				dm.client.ToggleDieState(dm.cursor)
			}
		}
	}
	return dm, cmd
}

func (dm diceModel) View() string {
	var s string

	var dice g.PlayerDice

	if dm.myDice {
		dice = dm.client.GetMyDice()
	} else {
		dice = dm.client.GetOpponentDice()
	}

	var styledDice []string
	selector := " "
	gameDice := dm.client.GetGameDice()
	for idx, die := range gameDice {
		face := die.GetFaces()[dice[idx].GetFaceId()]

		if dm.myDice && idx == dm.cursor && ph.Phase() == c.PickDice {
			selector = "^"
		} else {
			selector = " "
		}

		dieStyle := baseDieStyle
		if face.IsMagic() {
			dieStyle = dieStyle.BorderForeground(lipgloss.Color("#e8d102"))
		} else {
			dieStyle = dieStyle.BorderForeground(lipgloss.Color("#6d6d6d"))
		}

		dieBoxStyle := baseDieBoxStyle
		if dice[idx].IsKept() {
			dieBoxStyle = dieBoxStyle.BorderStyle(lipgloss.ThickBorder()).BorderBottom(true).BorderForeground(lipgloss.Color("#05c11e"))
		}

		styledDice = append(styledDice, dieBoxStyle.Render(fmt.Sprintf("%s\n%s", dieStyle.Render(face.GetKind()), selector)))
	}

	s = fmt.Sprintf("%s\n", lipgloss.JoinHorizontal(lipgloss.Top, styledDice...))

	if ph.Phase() == c.RollDice && dm.myDice {
		if dm.validated {
			s += "Rolling dice...\n"
		} else {
			s += "\t> Roll dice\n"
		}
	}

	if ph.Phase() == c.WaitingDiceRoll && !dm.myDice {
		s += "Waiting for other player to roll dice...\n"
	}

	return s
}
