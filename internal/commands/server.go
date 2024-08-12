package commands

import (
	cmn "github.com/e-gloo/orlog/internal/commons"
)

// Server to Client
const (
	CreateOrJoin    Command = "create_or_join"
	CreatedOrJoined Command = "created_or_joined"
	ConfigurePlayer Command = "configure_player"
	AskRollDice     Command = "ask_to_roll_dice"
	DiceRoll        Command = "dice_roll"
	SelectDice      Command = "select_dice"
	DiceState       Command = "dice_state"
	AskToPlayGod    Command = "ask_to_play_god"
	TurnFinished    Command = "turn_finished"
	GameStarting    Command = "game_starting"
	GameFinished    Command = "game_finished"
	CommandError    Command = "error"
)

type CreateOrJoinMessage struct{}

type CreatedOrJoinedMessage struct {
	Uuid string
	Dice []cmn.InitGameDie
	Gods []cmn.InitGod
}

type ConfigurePlayerMessage struct {
	Gods []cmn.InitGod
}

type AskRollDiceMessage struct {
	Player string
}

type DiceRollMessage struct {
	Player    string
	DiceState cmn.DiceState
}

type SelectDiceMessage struct {
	Player string
	Turn   int
}

type DiceStateMessage struct {
	DiceState cmn.PlayerMap[cmn.DiceState]
}

type AskToPlayGodMessage struct {
	Player string
}

type TurnFinishedMessage struct {
	Turn    int
	Players cmn.PlayerMap[cmn.UpdateGamePlayer]
}

type GameStartingMessage struct {
	YourUsername string
	Players      cmn.PlayerMap[cmn.InitGamePlayer]
}

type GameFinishedMessage struct {
	Winner string
}

type CommandErrorMessage struct {
	Reason string
}
