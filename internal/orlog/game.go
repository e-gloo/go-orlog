package orlog

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Game struct {
	Uuid    string
	Player1 *Player
	Player2 *Player
}

func (g *Game) PlayTurn(turn int) {
	players := [2]*Player{g.Player1, g.Player2}
	for player_idx := range players {
		fmt.Println("Turn", turn, players[player_idx].Name)
		players[player_idx].RollDices()
		PrintDices(players[player_idx].Dices)

		// We dont pick the dices
		if turn > 2 {
			continue
		}
		input := ""
		_, err := fmt.Scanln(&input)
		if err != nil && err.Error() != "unexpected newline" {
			return
		}

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
	g.Player1.GainTokens()
	g.Player2.GainTokens()

	// damage phase
	g.Player1.AttackPlayer(g.Player2)
	if g.Player2.Health <= 0 {
		return
	}
	g.Player2.AttackPlayer(g.Player1)

	// thief phase
	g.Player1.StealTokens(g.Player2)
	g.Player2.StealTokens(g.Player1)

	fmt.Printf("%s: %dHP, %dTK\n", g.Player1.Name, g.Player1.Health, g.Player1.Tokens)
	fmt.Printf("%s: %dHP, %dTK\n", g.Player2.Name, g.Player2.Health, g.Player2.Tokens)

	g.Player1.UnkeepDices()
	g.Player2.UnkeepDices()
}

func (g *Game) changePlayersPosition() {
	tmp := g.Player1
	g.Player1 = g.Player2
	g.Player2 = tmp
}

func (g *Game) selectFirstPlayer() {
	firstPlayer := rand.Intn(2)
	if firstPlayer == 1 {
		g.changePlayersPosition()
	}
}

func (g *Game) SetPlayer1(name string) {
	g.Player1 = NewPlayer(name)
}

func (g *Game) SetPlayer2(name string) {
	g.Player2 = NewPlayer(name)
}

func (g *Game) IsGameReady() bool {
	return g.Player1 != nil && g.Player2 != nil
}

func InitGame() (*Game, error) {
	newuuid, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("Error at generating uuid", err)
		return nil, err
	}
	game := &Game{
		Uuid: newuuid.String(),
	}

	return game, nil
}
