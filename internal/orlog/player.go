package orlog

import (
	"fmt"
)

type Player struct {
	Name   string `json:"name"`
	Health int    `json:"health"`
	Tokens int    `json:"tokens"`
	Dice   [6]Die `json:"dice"`
	// gods []God
}

func NewPlayer(name string) *Player {
	return &Player{
		Name:   name,
		Health: 15,
		Tokens: 0,
		Dice:   InitDice(),
	}
}

func (p *Player) RollDice() {
	for idx := range p.Dice {
		if !p.Dice[idx].Kept {
			p.Dice[idx].Roll()
		}
	}
}

func (p *Player) UnkeepDice() {
	for idx := range p.Dice {
		p.Dice[idx].Kept = false
	}
}

func (p *Player) AttackPlayer(player *Player) {
	nbArrows := 0
	nbAxes := 0
	for _, die := range p.Dice {
		if die.Face().Kind == Arrow {
			nbArrows++
		}
		if die.Face().Kind == Axe {
			nbAxes++
		}
	}

	nbHelmets := 0
	nbShields := 0
	for _, die := range player.Dice {
		if die.Face().Kind == Helmet {
			nbHelmets++
		}
		if die.Face().Kind == Shield {
			nbShields++
		}
	}
	arrowDamage := max(nbArrows-nbShields, 0)
	axeDamage := max(nbAxes-nbHelmets, 0)

	player.Health = player.Health - arrowDamage - axeDamage
}

func (p *Player) GainTokens() {
	nbMagics := 0

	for _, die := range p.Dice {
		if die.Face().Magic {
			nbMagics++
		}
	}

	p.Tokens += nbMagics
}

func (p *Player) StealTokens(player *Player) {
	nbThieves := 0

	for _, die := range p.Dice {
		if die.Face().Kind == Thief {
			nbThieves++
		}
	}

	nbTokens := min(nbThieves, player.Tokens)
	p.Tokens += nbTokens
	player.Tokens -= nbTokens
}

func (p *Player) FormatDice() string {
	var res string

	for dice_nb, die := range p.Dice {
		res = fmt.Sprintf(
			"%s%d %s",
			res,
			1+dice_nb,
			die.Face().String(),
		)
	}

	return res
}
