package server

import (
	"fmt"

	og "github.com/e-gloo/orlog/internal/orlog"
	"github.com/google/uuid"
)

type ServerGame struct {
	Uuid string
	Data *og.Game
}

func NewServerGame() (*ServerGame, error) {
	newuuid, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("error generating uuid: %w", err)
	}

	return &ServerGame{
		Uuid: newuuid.String(),
		Data: og.NewGame(),
	}, nil
}
