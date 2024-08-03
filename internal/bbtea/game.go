package bbtea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	c "github.com/e-gloo/orlog/internal/client"
)

type gameModel struct {
	opponentHUD  hudModel
	playerHUD    hudModel
	opponentDice diceModel
	playerDice   diceModel
	client       c.Client
}

func initialGameModel(client c.Client) gameModel {
	return gameModel{
		opponentHUD:  initalHudModel(client, false),
		opponentDice: initialDiceModel(client, 6, false),
		playerHUD:    initalHudModel(client, true),
		playerDice:   initialDiceModel(client, 6, true),
		client:       client,
	}
}

func (gm gameModel) Init() tea.Cmd {
	return nil
}

func (gm gameModel) Update(msg tea.Msg) (gameModel, tea.Cmd) {
	var cmd tea.Cmd

	switch ph.Phase() {
	case c.RollDice, c.PickDice:
		gm.playerDice, cmd = gm.playerDice.Update(msg)
	}
	return gm, cmd
}

func (gm gameModel) View() string {
	return fmt.Sprintf(
		"%s\n%s\n\n%s\n%s\n\n",
		gm.opponentHUD.View(),
		gm.opponentDice.View(),
		gm.playerHUD.View(),
		gm.playerDice.View(),
	)
}
