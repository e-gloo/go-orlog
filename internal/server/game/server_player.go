package server_game

import (
	"log/slog"

	c "github.com/e-gloo/orlog/internal/commands"
	"github.com/gorilla/websocket"
)

type Player interface {
	GetUsername() string
	GetHealth() int
	GetTokens() int
}

type GodChoice struct {
	index int
	level int
}

type ServerPlayer struct {
	username  string
	health    int
	tokens    int
	dice      PlayerDice
	gods      [3]int
	godChoice *GodChoice

	Conn     *websocket.Conn
	Expected []c.Command
}

func NewServerPlayer(conn *websocket.Conn, username string, godIndexes [3]int) *ServerPlayer {
	dice := PlayerDice{}
	for idx := range dice {
		dice[idx] = NewPlayerDie()
	}

	return &ServerPlayer{
		username:  username,
		health:    15,
		tokens:    50, // FIXME: should be 0
		dice:      dice,
		gods:      godIndexes,
		godChoice: nil,
		Conn:      conn,
		Expected:  []c.Command{},
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

func (sp *ServerPlayer) GetGods() [3]int {
	return sp.gods
}

func (p *ServerPlayer) GetGodChoice() *GodChoice {
	return p.godChoice
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
	sp.tokens = 50 // FIXME: should be 0
	sp.godChoice = nil
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

func (sp *ServerPlayer) SelectGod(godIndex int, godLevel int) {
	if godIndex != -1 && godLevel != -1 {
		sp.godChoice = &GodChoice{
			index: godIndex,
			level: godLevel,
		}
	} else {
		sp.godChoice = nil
	}
}

func (sp *ServerPlayer) GainTokens(dices [6]ServerDie) int {
	nbMagics := sp.assertFaces(
		dices,
		func(face *ServerFace) bool { return face.magic },
	)

	sp.tokens += nbMagics

	return nbMagics
}

func (sp *ServerPlayer) AttackPlayer(dices [6]ServerDie, player *ServerPlayer) (int, int) {
	nbArrows := sp.assertFaces(dices, func(face *ServerFace) bool { return face.kind == Arrow })
	nbAxes := sp.assertFaces(dices, func(face *ServerFace) bool { return face.kind == Axe })

	nbHelmets := player.assertFaces(dices, func(face *ServerFace) bool { return face.kind == Helmet })
	nbShields := player.assertFaces(dices, func(face *ServerFace) bool { return face.kind == Shield })

	arrowDamage := max(nbArrows-nbShields, 0)
	axeDamage := max(nbAxes-nbHelmets, 0)

	player.health = player.health - arrowDamage - axeDamage

	return arrowDamage, axeDamage
}

func (sp *ServerPlayer) StealTokens(dices [6]ServerDie, player *ServerPlayer) int {
	nbThieves := sp.assertFaces(dices, func(face *ServerFace) bool { return face.kind == Thief })

	nbToken := min(nbThieves, player.tokens)
	sp.tokens += nbToken
	player.tokens -= nbToken

	return nbToken
}

func (sp *ServerPlayer) activateGod(god *God, level int, game *ServerGame, opponent *ServerPlayer) bool {
	if god.Levels[level].TokenCost >= sp.tokens {
		return false
	}

	sp.tokens -= god.Levels[level].TokenCost
	slog.Debug("god activation", "username", sp.username, "god", god.Name, "level", god.Levels[level].Description)
	god.Activate(game, sp, opponent, god, level)
	return true
}
