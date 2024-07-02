package server

import (
	og "github.com/e-gloo/orlog/internal/orlog"
)

type ServerPlayer struct {
	Data *og.Player
}

func NewServerPlayer(username string) *ServerPlayer {
	return &ServerPlayer{
		Data: og.NewPlayer(username),
	}
}
