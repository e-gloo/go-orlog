package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

type Game struct {
	players [2]Player
}

func (g *Game) PlayTurn(turn int) {
	for player_idx, _ := range g.players {
		fmt.Println("Turn", turn, g.players[player_idx].name)
		g.players[player_idx].RollDices()
		printDices(g.players[player_idx].dices)

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
			g.players[player_idx].dices[i-1].kept = true
		}
	}
}

func (g *Game) PlayRound() {
	for i := 0; i < 1; i++ {
		g.PlayTurn(i + 1)
	}
    g.players[0].AttackPlayer(&g.players[1])
    g.players[1].AttackPlayer(&g.players[0])
}

func (g *Game) changePlayersPosition() {
	for idx := range g.players {
		if g.players[idx].position == 1 {
			g.players[idx].position = 2
		} else {
			g.players[idx].position = 1
		}
	}
	sort.Slice(g.players[:], func(i, j int) bool {
		return g.players[i].position < g.players[j].position
	})
}

func (g *Game) selectFirstPlayer() {
	firstPlayer := rand.Intn(2)
	for idx := range g.players {
		if idx == firstPlayer {
			g.players[idx].position = 1
		} else {
			g.players[idx].position = 2
		}
	}
	sort.Slice(g.players[:], func(i, j int) bool {
		return g.players[i].position < g.players[j].position
	})

}

func InitGame() *Game {
	game := &Game{
		players: InitPlayers(),
	}
	game.selectFirstPlayer()
	return game
}

func (g *Game) Play() {
gameLoop:
	for {
		g.PlayRound()
		for idx := range g.players {
			if g.players[idx].health == 0 {
				break gameLoop
			}
		}
		g.changePlayersPosition()
	}
}
