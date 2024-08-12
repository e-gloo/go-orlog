package bbtea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	c "github.com/e-gloo/orlog/internal/client"
)

type gameModel struct {
	opponentHUD  tea.Model
	playerHUD    tea.Model
	opponentDice tea.Model
	playerDice   tea.Model
	opponentGods tea.Model
	playerGods   tea.Model
	client       c.Client
}

func initialGameModel(client c.Client) tea.Model {
	return gameModel{
		opponentHUD:  initalHudModel(client, false),
		opponentDice: initialDiceModel(client, 6, false),
		opponentGods: initialGodModel(client, false),
		playerHUD:    initalHudModel(client, true),
		playerDice:   initialDiceModel(client, 6, true),
		playerGods:   initialGodModel(client, true),
		client:       client,
	}
}

func (gm gameModel) Init() tea.Cmd {
	return nil
}

func (gm gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case c.Phase:
		ph.SetPhase(msg)
		switch msg {
		case c.RollDice:
			cmd = gm.playerDice.Init()
		case c.DiceRoll:
			gm.playerDice, cmd = gm.playerDice.Update(Cmd(false))
		}
	default:
		switch ph.Phase() {
		case c.RollDice, c.PickDice:
			gm.playerDice, cmd = gm.playerDice.Update(msg)
		case c.SelectGod:
			gm.playerGods, cmd = gm.playerGods.Update(msg)
		}
	}

	return gm, cmd
}

func (gm gameModel) View() string {
	return fmt.Sprintf(
		"%s\n%s\n%s\n\n%s\n%s\n%s\n\n",
		gm.opponentHUD.View(),
		gm.opponentDice.View(),
		gm.opponentGods.View(),
		gm.playerHUD.View(),
		gm.playerDice.View(),
		gm.playerGods.View(),
	)
}
