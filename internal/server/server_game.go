package server

import (
	"encoding/json"
	"fmt"

	og "github.com/e-gloo/orlog/internal/orlog"
	"github.com/google/uuid"
)

type ServerGame struct {
	Uuid  string
	Data  *og.Game
	Rolls int
}

func NewServerGame() (*ServerGame, error) {
	newuuid, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("error generating uuid: %w", err)
	}

	return &ServerGame{
		Uuid:  newuuid.String(),
		Data:  og.NewGame(),
		Rolls: 0,
	}, nil
}

func (sg *ServerGame) String() (string, error) {
	jsonStr, err := json.Marshal(sg.Data)
	if err != nil {
		return "", fmt.Errorf("error marshalling game: %w", err)
	}

	return string(jsonStr), nil
}
