package server_game

import (
	c "github.com/e-gloo/orlog/internal/commands"
	"github.com/gorilla/websocket"
)

type Player interface {
	GetUsername() string
	GetHealth() int
	GetTokens() int
}

type ServerPlayer struct {
	username string
	health   int
	tokens   int
	dice     PlayerDice

	Conn     *websocket.Conn
	Expected []c.Command
}

func NewServerPlayer(conn *websocket.Conn, username string) *ServerPlayer {
	dice := PlayerDice{}
	for idx := range dice {
		dice[idx] = NewPlayerDie()
	}

	return &ServerPlayer{
		Conn:     conn,
		Expected: []c.Command{},
		username: username,
		dice:     dice,
		health:   0,
		tokens:   0,
	}
}

func (sp *ServerPlayer) GetUsername() string {
	return sp.username
}

func (sp *ServerPlayer) GetHealth() int {
	return sp.health
}

func (sp *ServerPlayer) GetTokens() int {
	return sp.tokens
}

func (sp *ServerPlayer) GetDice() PlayerDice {
	return sp.dice
}

func (sp *ServerPlayer) RollDice() {
	for idx := range sp.dice {
		if !sp.dice[idx].IsKept() {
			sp.dice[idx].Roll()
		}
	}
}

func (sp *ServerPlayer) ResetDice() {
	for idx := range sp.dice {
		sp.dice[idx].Reset()
	}
}
