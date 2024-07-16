package commands

import (
	cmn "github.com/e-gloo/orlog/internal/commons"
)

// Server to Client
const (
	CreateOrJoin    Command = "create_or_join"
	CreatedOrJoined Command = "created_or_joined"
	ConfigurePlayer Command = "configure_player"
	SelectDice      Command = "select_dice"
	WantToPlayGods  Command = "want_to_play_gods"
	GameStarting    Command = "starting"
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

type SelectDiceMessage struct {
	Turn    int
	Players cmn.PlayerMap[cmn.DiceState]
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

type CommandErrorMessage struct {
	Reason string
}
