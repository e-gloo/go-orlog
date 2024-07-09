package server

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ServerGame struct {
	Uuid         string
	Rolls        int
	Dice         [6]ServerDie
	Players      map[string]*ServerPlayer
	PlayersOrder []string
}

func NewServerGame() (*ServerGame, error) {
	newuuid, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("error generating uuid: %w", err)
	}

	return &ServerGame{
		Uuid:         newuuid.String(),
		Rolls:        0,
		Dice:         InitDice(),
		Players:      make(map[string]*ServerPlayer),
		PlayersOrder: make([]string, 0, 2),
	}, nil
}

func (g *ServerGame) AddPlayer(conn *websocket.Conn, name string) error {
	if len(g.PlayersOrder) != 0 && g.Players[g.PlayersOrder[0]] != nil && g.Players[g.PlayersOrder[0]].GetUsername() == name {
		return fmt.Errorf("name already exists")
	} else if len(g.PlayersOrder) >= 2 {
		return fmt.Errorf("game is full")
	} else {
		g.PlayersOrder = append(g.PlayersOrder, name)
		g.Players[name] = NewServerPlayer(conn, name)

		return nil
	}
}

func (g *ServerGame) IsGameReady() bool {
	return len(g.PlayersOrder) == 2 &&
		g.Players[g.PlayersOrder[0]] != nil &&
		g.Players[g.PlayersOrder[1]] != nil
}

func (g *ServerGame) ChangePlayersPosition() {
	tmp := g.PlayersOrder[0]
	g.PlayersOrder[0] = g.PlayersOrder[1]
	g.PlayersOrder[1] = tmp
}

func (g *ServerGame) SelectFirstPlayer() {
	firstPlayer := rand.Intn(2)
	if firstPlayer == 1 {
		g.ChangePlayersPosition()
	}
}
