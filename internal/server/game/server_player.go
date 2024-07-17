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
		health:   15,
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

func (sp *ServerPlayer) Reset() {
	sp.health = 15
	sp.tokens = 0
	sp.ResetDice()
}

func (sp *ServerPlayer) assertFaces(dices [6]ServerDie, assert func(f *ServerFace) bool) int {
	count := 0
	for dieIdx, die := range dices {
		if assert(&die.faces[sp.dice[dieIdx].faceIndex]) {
			count += sp.dice[dieIdx].quantity
		}
	}
	return count
}

func (sp *ServerPlayer) GainTokens(dices [6]ServerDie) {
	nbMagics := sp.assertFaces(
		dices,
		func(face *ServerFace) bool { return face.magic },
	)

	sp.tokens += nbMagics
}

func (sp *ServerPlayer) AttackPlayer(dices [6]ServerDie, player *ServerPlayer) {
	nbArrows := sp.assertFaces(dices, func(face *ServerFace) bool { return face.kind == Arrow })
	nbAxes := sp.assertFaces(dices, func(face *ServerFace) bool { return face.kind == Axe })

	nbHelmets := player.assertFaces(dices, func(face *ServerFace) bool { return face.kind == Helmet })
	nbShields := player.assertFaces(dices, func(face *ServerFace) bool { return face.kind == Shield })

	arrowDamage := max(nbArrows-nbShields, 0)
	axeDamage := max(nbAxes-nbHelmets, 0)

	player.health = player.health - arrowDamage - axeDamage
}

func (sp *ServerPlayer) StealTokens(dices [6]ServerDie, player *ServerPlayer) {
	nbThieves := sp.assertFaces(dices, func(face *ServerFace) bool { return face.kind == Thief })

	nbToken := min(nbThieves, player.tokens)
	sp.tokens += nbToken
	player.tokens -= nbToken
}
