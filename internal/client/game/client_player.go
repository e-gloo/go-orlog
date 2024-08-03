package client_game

import (
	cmn "github.com/e-gloo/orlog/internal/commons"
)

type ClientPlayer struct {
	username string
	health   int
	tokens   int
	dice     PlayerDice
	gods     [3]int
}

func NewClientPlayer(init cmn.InitGamePlayer) *ClientPlayer {
	dice := PlayerDice{}
	for idx := range dice {
		dice[idx] = NewPlayerDie()
	}

	return &ClientPlayer{
		username: init.Username,
		health:   init.Health,
		tokens:   0,
		dice:     dice,
		gods:     init.GodIndexes,
	}
}

func (p *ClientPlayer) GetUsername() string {
	return p.username
}

func (p *ClientPlayer) GetHealth() int {
	return p.health
}

func (p *ClientPlayer) GetTokens() int {
	return p.tokens
}

func (p *ClientPlayer) GetDice() PlayerDice {
	return p.dice
}

func (p *ClientPlayer) GetGods() [3]int {
	return p.gods
}

func (p *ClientPlayer) update(update cmn.UpdateGamePlayer) {
	p.health = update.Health
	p.tokens = update.Tokens
}

func (p *ClientPlayer) updateDice(update cmn.DiceState) {
	for i, state := range update {
		p.dice[i].SetFaceId(state.Index)
		p.dice[i].SetKept(state.Kept)
	}
}
