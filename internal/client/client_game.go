package client

import (
	og "github.com/e-gloo/orlog/internal/orlog"
)

type ClientGame struct {
	Data *og.Game
}

func NewClientGame(game *og.Game) *ClientGame {
	return &ClientGame{
		Data: game,
	}
}
