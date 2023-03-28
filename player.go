package main

import (
	"fmt"
)

type Player struct {
	name     string
	health   int
	tokens   int
	dices    [6]Die
	position int
	gods     [3]*God
}

func (p *Player) RollDices() {
	for idx, _ := range p.dices {
		if p.dices[idx].kept == false {
			p.dices[idx].Roll()
		}
	}
}

func (p *Player) ResetDices() {
	for idx, _ := range p.dices {
		p.dices[idx].ResetDie()
	}
}

func (p *Player) AttackPlayer(player *Player) {
	nbArrows := AssertFaces(p.dices, func(face *Face) bool { return face.kind == Arrow })
	nbAxes := AssertFaces(p.dices, func(face *Face) bool { return face.kind == Axe })

	nbHelmets := AssertFaces(player.dices, func(face *Face) bool { return face.kind == Helmet })
	nbShields := AssertFaces(player.dices, func(face *Face) bool { return face.kind == Shield })

	arrowDamage := Max(nbArrows-nbShields, 0)
	axeDamage := Max(nbAxes-nbHelmets, 0)

	player.health = player.health - arrowDamage - axeDamage
}

func (p *Player) GainTokens() {
	nbMagics := AssertFaces(p.dices, func(face *Face) bool { return face.magic == true })

	p.tokens += nbMagics
}

func (p *Player) StealTokens(player *Player) {
	nbThieves := AssertFaces(player.dices, func(face *Face) bool { return face.kind == Thief })

	nbToken := Min(nbThieves, player.tokens)
	p.tokens += nbToken
	player.tokens -= nbToken
}

func (p *Player) AskForGod() (*God, int) {
	if p.tokens <= 0 {
		return nil, -1
	}

	PrintGods(p.gods[:])
	fmt.Println("Activate a god : ")
	input := ""
	fmt.Scanln(&input)

	if input == "" {
		return nil, -1
	}

	choosen, err := StringToIntArray(input)
	if err != nil {
		return nil, -1
	}

	PrintGodLevels(p.gods[choosen[0]-1])
	fmt.Scanln(&input)
	levels, err := StringToIntArray(input)
	if err != nil {
		return nil, -1
	}

	return p.gods[choosen[0] - 1], levels[0] - 1
}

func (p *Player) ActivateGod(god *God, level int, opponent *Player, currentPriority int) bool {
	if god == nil || level == -1 || god.Priority != currentPriority || god.Levels[level].TokenCost > p.tokens {
		return false
	}

	p.tokens -= god.Levels[level].TokenCost
	fmt.Printf("%s activates %s : %s\n", p.name, god.Name, god.Levels[level].Description)
	god.Activate(p, opponent, god, level)
	return true
}

func InitPlayer(gods []*God) *Player {
	player := &Player{
		name:     "Player",
		health:   15,
		tokens:   50,
		dices:    InitDices(),
		position: 1,
		gods:     [3]*God{nil, nil, nil},
	}

	fmt.Println("Enter your name : ")
	fmt.Scanln(&player.name)

	PrintGods(gods)

	fmt.Println("Choose 3 gods : ")
	input := ""
	fmt.Scanln(&input)

	choosen, err := StringToIntArray(input)
	if err != nil {
		return player
	}

	for idx, godIdx := range choosen[:3] {
		player.gods[idx] = gods[godIdx-1]
	}

	return player
}
