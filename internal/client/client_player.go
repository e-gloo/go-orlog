package client

import (
	og "github.com/e-gloo/orlog/internal/orlog"
)

type ClientPlayer struct {
	Data *og.Player
}

func NewClientPlayer(username string) *ClientPlayer {
	return &ClientPlayer{
		Data: og.NewPlayer(username),
	}
}
