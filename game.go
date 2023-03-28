package main

import (
	"fmt"
	"math/rand"
)

type Game struct {
	player1 *Player
	player2 *Player
}

func (g *Game) PlayTurn(turn int) {
	players := [2]*Player{g.player1, g.player2}
	for player_idx, _ := range players {
		fmt.Println("Turn", turn, players[player_idx].name)
		players[player_idx].RollDices()
		printDices(players[player_idx].dices)

		keptDices := AssertDices(players[player_idx].dices, func (d *Die) bool { return d.kept == true })

		// We dont pick the dices
		if turn > 2 || keptDices >= 6 {
			continue
		}
		input := ""
		fmt.Scanln(&input)

		if input == "*" {
			input = "1,2,3,4,5,6"
		}
		to_keep, err := StringToIntArray(input)
		if err != nil {
			continue
		}

		for _, dice_nb := range to_keep {
			players[player_idx].dices[dice_nb-1].kept = true
		}
	}
}

func (g *Game) PlayRound() {
	for i := 0; i < 3; i++ {
		g.PlayTurn(i + 1)
	}
	// ask if should use god
	p1god, p1godLevel := g.player1.AskForGod()
	p2god, p2godLevel := g.player2.AskForGod()

	g.player1.ActivateGod(p1god, p1godLevel, g.player2, 1)
	g.player2.ActivateGod(p2god, p2godLevel, g.player1, 1)

	g.player1.ActivateGod(p1god, p1godLevel, g.player2, 2)
	g.player2.ActivateGod(p2god, p2godLevel, g.player1, 2)

	// gain tokens
	g.player1.GainTokens()
	g.player2.GainTokens()

	g.player1.ActivateGod(p1god, p1godLevel, g.player2, 3)
	g.player2.ActivateGod(p2god, p2godLevel, g.player1, 3)

	// damage phase
	g.player1.AttackPlayer(g.player2)
	if g.player2.health <= 0 {
		return
	}
	g.player2.AttackPlayer(g.player1)

	g.player1.ActivateGod(p1god, p1godLevel, g.player2, 4)
	g.player2.ActivateGod(p2god, p2godLevel, g.player1, 4)

	// thief phase
	g.player1.StealTokens(g.player2)
	g.player2.StealTokens(g.player1)

	g.player1.ActivateGod(p1god, p1godLevel, g.player2, 5)
	g.player2.ActivateGod(p2god, p2godLevel, g.player1, 5)

	g.player1.ActivateGod(p1god, p1godLevel, g.player2, 6)
	g.player2.ActivateGod(p2god, p2godLevel, g.player1, 6)

	g.player1.ActivateGod(p1god, p1godLevel, g.player2, 7)
	g.player2.ActivateGod(p2god, p2godLevel, g.player1, 7)

	fmt.Printf("%s: %dHP, %dTK\n", g.player1.name, g.player1.health, g.player1.tokens)
	fmt.Printf("%s: %dHP, %dTK\n", g.player2.name, g.player2.health, g.player2.tokens)

	g.player1.ResetDices()
	g.player2.ResetDices()
}

func (g *Game) changePlayersPosition() {
	tmp := g.player1
	g.player1 = g.player2
	g.player2 = tmp
}

func (g *Game) selectFirstPlayer() {
	firstPlayer := rand.Intn(2)
	if firstPlayer == 1 {
		g.changePlayersPosition()
	}
}

func InitGame() *Game {
	gods := InitGods()
	game := &Game{
		player1: InitPlayer(gods),
		player2: InitPlayer(gods),
	}
	game.selectFirstPlayer()

	return game
}

func (g *Game) Play() {
gameLoop:
	for {
		g.PlayRound()
		if g.player2.health <= 0 {
			// P1 won
			break gameLoop
		} else if g.player1.health <= 0 {
			// P2 won
			break gameLoop
		}
		g.changePlayersPosition()
	}
}
