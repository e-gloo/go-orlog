package main

import "fmt"

type Player struct {
	name   string
	health int
	token  int
	dices  [6]Dice
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
    nbMagics := 0
    nbThieves := 0
    nbArrows := 0
    nbAxes := 0
    for _, dice := range p.dices {
        if dice.Face().magic == true {
            nbMagics++
        }
        if dice.Face().kind == Thief {
            nbThieves++
        }
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
    var arrowDamage int
    var axeDamage int
    if nbArrows - nbShields > 0 {
        arrowDamage = nbArrows - nbShields
    } else {
        arrowDamage = 0
    }
    if nbAxes - nbHelmets > 0 {
        axeDamage = nbAxes - nbHelmets
    } else {
        axeDamage = 0
    }
    player.health = player.health - arrowDamage - axeDamage

    if nbThieves > player.token {
        p.token += player.token
        player.token = 0
    } else {
        p.token += nbThieves
        player.token -= nbThieves
    }

    p.token += nbMagics

    fmt.Printf("New stats %s, %d, %d\n", p.name, p.health, p.token)
    fmt.Printf("New stats %s, %d, %d\n", player.name, player.health, player.token)
}

func InitPlayers() [2]Player {
	players := [2]Player{
		{
			name:   "Player 1",
			health: 15,
			token:  0,
			dices:  InitDices(),
            position: 1,
		},
		{
			name:   "Player 2",
			health: 15,
			token:  0,
			dices:  InitDices(),
            position: 2,
		},
	}

	// TODO: Choose name
	// TODO: Choose gods

	return players
}
