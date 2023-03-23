package main

import "fmt"

type Player struct {
	name     string
	health   int
	token    int
	dices    [6]Die
	position int
	// gods []God
}

func (p *Player) RollDices() {
	for idx, _ := range p.dices {
		if p.dices[idx].kept == false {
			p.dices[idx].Roll()
		}
	}
}

func (p *Player) UnkeepDices() {
	for idx, _ := range p.dices {
		p.dices[idx].kept = false
	}
}

func (p *Player) AttackPlayer(player *Player) {
	nbArrows := 0
	nbAxes := 0
	for _, die := range p.dices {
		if die.Face().kind == Arrow {
			nbArrows++
		}
		if die.Face().kind == Axe {
			nbAxes++
		}
	}

	nbHelmets := 0
	nbShields := 0
	for _, die := range player.dices {
		if die.Face().kind == Helmet {
			nbHelmets++
		}
		if die.Face().kind == Shield {
			nbShields++
		}
	}
	arrowDamage := Max(nbArrows-nbShields, 0)
	axeDamage := Max(nbAxes-nbHelmets, 0)

	player.health = player.health - arrowDamage - axeDamage
}

func (p *Player) GainTokens() {
	nbMagics := 0

	for _, die := range p.dices {
		if die.Face().magic == true {
			nbMagics++
		}
	}

	p.token += nbMagics
}

func (p *Player) StealTokens(player *Player) {
	nbThieves := 0

	for _, die := range p.dices {
		if die.Face().kind == Thief {
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
		health:   15,
		token:    0,
		dices:    InitDices(),
		position: 1,
	}

	fmt.Println("Enter your name : ")
	fmt.Scanln(&player.name)

	// TODO: Choose gods
    // https://www.thegamer.com/assassins-creed-valhalla-orlog-god-favors/

	return player
}
