package main

type Player struct {
	name   string
	health int
	token  int
	dices  [6]Dice
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

func InitPlayers() [2]Player {
	players := [2]Player{
		Player{
			name:   "Player 1",
			health: 15,
			token:  0,
			dices:  InitDices(),
		},
		Player{
			name:   "Player 2",
			health: 15,
			token:  0,
			dices:  InitDices(),
		},
	}

	// TODO: Choose name
	// TODO: Choose gods

	return players
}
