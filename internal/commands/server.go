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
	AskToPlayGod    Command = "ask_to_play_god"
	TurnFinished    Command = "turn_finished"
	GameStarting    Command = "game_starting"
	GameFinished    Command = "game_finished"
	CommandError    Command = "error"
)

type CreateOrJoinMessage struct{}

type CreatedOrJoinedMessage struct {
	Uuid string
}

type ConfigurePlayerMessage struct {
	Gods []cmn.InitGod
}

type DiceRollMessage struct {
	Players cmn.PlayerMap[cmn.DiceState]
}

type SelectDiceMessage struct {
	Turn int
}

type AskToPlayGodMessage struct{}

type TurnFinishedMessage struct {
	Turn    int
	Players cmn.PlayerMap[cmn.UpdateGamePlayer]
}

type GameStartingMessage struct {
	YourUsername string
	Dice         []cmn.InitGameDie
	Gods         []cmn.InitGod
	Players      cmn.PlayerMap[cmn.InitGamePlayer]
}

type GameFinishedMessage struct {
	Winner string
}

type CommandErrorMessage struct {
	Reason string
}
