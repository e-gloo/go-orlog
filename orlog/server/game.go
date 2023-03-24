package server

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/e-gloo/orlog/orlog/commons"
	"github.com/google/uuid"
)

type Game struct {
	uuid    string
	player1 *commons.Player
	player2 *commons.Player
}

func (g *Game) PlayTurn(turn int) {
	players := [2]*commons.Player{g.player1, g.player2}
	for player_idx := range players {
		fmt.Println("Turn", turn, players[player_idx].Name)
		players[player_idx].RollDices()
		commons.PrintDices(players[player_idx].Dices)

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
			players[player_idx].Dices[i-1].Kept = true
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
	if g.player2.Health <= 0 {
		return
	}
	g.player2.AttackPlayer(g.player1)

	// thief phase
	g.player1.StealTokens(g.player2)
	g.player2.StealTokens(g.player1)

	fmt.Printf("%s: %dHP, %dTK\n", g.player1.Name, g.player1.Health, g.player1.Token)
	fmt.Printf("%s: %dHP, %dTK\n", g.player2.Name, g.player2.Health, g.player2.Token)

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

func InitGame(player *commons.Player) (*Game, error) {
	newuuid, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("Error at generating uuid", err)
		return nil, err
	}
	game := &Game{
		uuid:    newuuid.String(),
		player1: player,
	}
	//game.selectFirstPlayer()

	return game, nil
}

func (g *Game) Play() {
}
