package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

		// We dont pick the dices
		if turn > 2 {
			continue
		}
		input := ""
		fmt.Scanln(&input)

		to_keep := strings.Split(input, ",")

		for _, dice_nb := range to_keep {
			i, err := strconv.ParseInt(dice_nb, 10, 32)
			if err != nil {
				continue
			}
			players[player_idx].dices[i-1].kept = true
		}
	}
}

func (g *Game) PlayRound() {
	for i := 0; i < 3; i++ {
		g.PlayTurn(i + 1)
	}
	// ask if should use god

	// gain tokens
	g.player1.GainTokens()
	g.player2.GainTokens()

	// damage phase
	g.player1.AttackPlayer(g.player2)
	if (g.player2.health <= 0) {
		return
	}
	g.player2.AttackPlayer(g.player1)

	// thief phase
	g.player1.StealTokens(g.player2)
	g.player2.StealTokens(g.player1)

	fmt.Printf("%s: %dHP, %dTK\n", g.player1.name, g.player1.health, g.player1.token)
	fmt.Printf("%s: %dHP, %dTK\n", g.player2.name, g.player2.health, g.player2.token)

	g.player1.UnkeepDices()
	g.player2.UnkeepDices()
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
	game := &Game{
		player1: InitPlayer(),
		player2: InitPlayer(),
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
