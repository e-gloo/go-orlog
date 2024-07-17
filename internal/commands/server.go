package commands

import (
	cmn "github.com/e-gloo/orlog/internal/commons"
)

// Server to Client
const (
	CreateOrJoin    Command = "create_or_join"
	CreatedOrJoined Command = "created_or_joined"
	ConfigurePlayer Command = "configure_player"
	DiceRoll        Command = "dice_roll"
	SelectDice      Command = "select_dice"
	WantToPlayGods  Command = "want_to_play_gods"
	GameStarting    Command = "starting"
	GameFinished    Command = "finished"
	CommandError    Command = "error"
)

type CreateOrJoinMessage struct {
	Welcome string
	// Lobbies []string
}

type CreatedOrJoinedMessage struct {
	Uuid string
}

type ConfigurePlayerMessage struct {
	Gods []cmn.GodDefinition
}

type DiceRollMessage struct {
	Players cmn.PlayerMap[cmn.DiceState]
}

type SelectDiceMessage struct {
	Turn int
}

type WantToPlaysGodsMessage struct {
	// maybe the list of selected gods,
	// or nothing if we trust client
}

type GameStartingMessage struct {
	YourUsername string
	Dice         []cmn.InitGameDie
	Players      cmn.PlayerMap[cmn.InitGamePlayer]
}

type GameFinishedMessage struct {
	Winner string
}

type CommandErrorMessage struct {
	Reason string
}
