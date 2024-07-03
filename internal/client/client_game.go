package client

import (
	og "github.com/e-gloo/orlog/internal/orlog"
)

type ClientGame struct {
	Uuid string
	Data *og.Game
}

func NewClientGame(uuid string) *ClientGame {
	game := og.NewGame()

	// TODO: hydrate game with data from server

	return &ClientGame{
		Uuid: uuid,
		Data: game,
	}
}
