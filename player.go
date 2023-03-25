package main

import (
	"fmt"
)

type Player struct {
	name     string
	health   int
	token    int
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
		p.dices[idx].kept = false
		for faceIdx, _ := range p.dices[idx].faces {
			p.dices[idx].faces[faceIdx].quantity = 1
		}
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

	p.token += nbMagics
}

func (p *Player) StealTokens(player *Player) {
	nbThieves := AssertFaces(player.dices, func(face *Face) bool { return face.kind == Thief })

	nbToken := Min(nbThieves, player.token)
	p.token += nbToken
	player.token -= nbToken
}

func (p *Player) AskForGod() (*God, int) {
	// FIXME: this SegFaults
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

	fmt.Scanln(&input)
	levels, err := StringToIntArray(input)
	if err != nil {
		return nil, -1
	}

	return p.gods[choosen[0]], levels[0]
}

func (p *Player) ActivateGod(god *God, level int, opponent *Player, currentPriority int) bool {
	if god == nil || level == -1 || god.Priority != currentPriority || god.Levels[level].TokenCost > p.token {
		return false
	}

	fmt.Printf("%s activates %s : %s\n", p.name, god.Name, god.Levels[level].Description)
	god.Activate(p, opponent, god, level)
	p.token -= god.Levels[level].TokenCost
	return true
}

func InitPlayer(gods []*God) *Player {
	player := &Player{
		name:     "Player",
		health:   15,
		token:    50,
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
	if err != nil || len(choosen) != 3 {
		return player
	}

	for idx, godIdx := range choosen[:3] {
		player.gods[idx] = gods[godIdx]
	}

	return player
}
