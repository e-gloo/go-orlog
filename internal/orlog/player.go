package orlog

type Player struct {
	Name   string `json:"name"`
	Health int    `json:"health"`
	Tokens int    `json:"tokens"`
	Dices  [6]Die `json:"dices"`
	// gods []God
}

func NewPlayer(name string) *Player {
	return &Player{
		Name:   name,
		Health: 15,
		Tokens: 0,
		Dices:  InitDices(),
	}
}

func (p *Player) RollDices() {
	for idx := range p.Dices {
		if !p.Dices[idx].Kept {
			p.Dices[idx].Roll()
		}
	}
}

func (p *Player) UnkeepDices() {
	for idx := range p.Dices {
		p.Dices[idx].Kept = false
	}
}

func (p *Player) AttackPlayer(player *Player) {
	nbArrows := 0
	nbAxes := 0
	for _, die := range p.Dices {
		if die.Face().Kind == Arrow {
			nbArrows++
		}
		if die.Face().Kind == Axe {
			nbAxes++
		}
	}

	nbHelmets := 0
	nbShields := 0
	for _, die := range player.Dices {
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

	for _, die := range p.Dices {
		if die.Face().Magic {
			nbMagics++
		}
	}

	p.Tokens += nbMagics
}

func (p *Player) StealTokens(player *Player) {
	nbThieves := 0

	for _, die := range p.Dices {
		if die.Face().Kind == Thief {
			nbThieves++
		}
	}

	nbTokens := min(nbThieves, player.Tokens)
	p.Tokens += nbTokens
	player.Tokens -= nbTokens
}
