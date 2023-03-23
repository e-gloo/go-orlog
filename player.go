package main

import "fmt"

type Player struct {
	name     string
	health   int
	token    int
	dices    [6]Dice
	position int
	// gods []God
}

func (p *Player) RollDices() {
	for idx, _ := range p.dices {
		if p.dices[idx].kept == false {
			p.dices[idx].Roll()
		}
		p.dices[idx].kept = false
	}
}

func (p *Player) AttackPlayer(player *Player) {
	nbArrows := 0
	nbAxes := 0
	for _, dice := range p.dices {
		if dice.Face().kind == Arrow {
			nbArrows++
		}
		if dice.Face().kind == Axe {
			nbAxes++
		}
	}

	nbHelmets := 0
	nbShields := 0
	for _, dice := range player.dices {
		if dice.Face().kind == Helmet {
			nbHelmets++
		}
		if dice.Face().kind == Shield {
			nbShields++
		}
	}
	arrowDamage := Max(nbArrows-nbShields, 0)
	axeDamage := Max(nbAxes-nbHelmets, 0)

	player.health = player.health - arrowDamage - axeDamage
}

func (p *Player) GainTokens() {
	nbMagics := 0

	for _, dice := range p.dices {
		if dice.Face().magic == true {
			nbMagics++
		}
	}

	p.token += nbMagics
}

func (p *Player) StealTokens(player *Player) {
	nbThieves := 0

	for _, dice := range p.dices {
		if dice.Face().kind == Thief {
			nbThieves++
		}
	}

	nbToken := Min(nbThieves, player.token)
	p.token += nbToken
	player.token -= nbToken
}

func InitPlayer() *Player {
	player := &Player{
		name:     "Player",
		health:   2,
		token:    0,
		dices:    InitDices(),
		position: 1,
	}

	fmt.Println("Enter your name : ")
	fmt.Scanln(&player.name)

	// TODO: Choose gods

	return player
}
