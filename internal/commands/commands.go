package commands

import (
	"github.com/e-gloo/orlog/internal/game"
)

const (
	Create     = "create"
	Join       = "join"
	ChooseGods = "choose_gods"
	PlayGods   = "play_gods"
	KeepDices  = "keep_dices"
)

const (
    SelectDices = "select_dices"
    WantToPlayGods = "want_to_play_gods"
)

type CreateData struct {
	Uuid   string  `json:"uuid"`
	Player *game.Player `json:"player"`
}

type Packet struct {
	Command string `json:"command"`
	Data    []byte `json:"data"`
}

