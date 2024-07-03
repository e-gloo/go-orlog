package orlog

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Game struct {
	Players      map[string]*Player
	PlayersOrder []string
}

func NewGame() *Game {
	return &Game{
		Players:      make(map[string]*Player),
		PlayersOrder: make([]string, 0),
	}
}

func (g *Game) PlayTurn(turn int) {
	for _, username := range g.PlayersOrder {
		fmt.Println("Turn", turn, g.Players[username].Name)
		g.Players[username].RollDices()
		PrintDices(g.Players[username].Dices)

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
			g.Players[username].Dices[i-1].Kept = true
		}
	}
}

func (g *Game) PlayRound() {
	for i := 0; i < 3; i++ {
		g.PlayTurn(i + 1)
	}
	// ask if should use god

	// gain tokens
	g.Players[g.PlayersOrder[0]].GainTokens()
	g.Players[g.PlayersOrder[1]].GainTokens()

	// damage phase
	g.Players[g.PlayersOrder[0]].AttackPlayer(g.Players[g.PlayersOrder[1]])
	if g.Players[g.PlayersOrder[1]].Health <= 0 {
		return
	}
	g.Players[g.PlayersOrder[1]].AttackPlayer(g.Players[g.PlayersOrder[0]])

	// thief phase
	g.Players[g.PlayersOrder[0]].StealTokens(g.Players[g.PlayersOrder[1]])
	g.Players[g.PlayersOrder[1]].StealTokens(g.Players[g.PlayersOrder[0]])

	fmt.Printf("%s: %dHP, %dTK\n", g.Players[g.PlayersOrder[0]].Name, g.Players[g.PlayersOrder[0]].Health, g.Players[g.PlayersOrder[0]].Tokens)
	fmt.Printf("%s: %dHP, %dTK\n", g.Players[g.PlayersOrder[1]].Name, g.Players[g.PlayersOrder[1]].Health, g.Players[g.PlayersOrder[1]].Tokens)

	g.Players[g.PlayersOrder[0]].UnkeepDices()
	g.Players[g.PlayersOrder[1]].UnkeepDices()
}

func (g *Game) changePlayersPosition() {
	tmp := g.PlayersOrder[0]
	g.PlayersOrder[0] = g.PlayersOrder[1]
	g.PlayersOrder[1] = tmp
}

func (g *Game) SelectFirstPlayer() {
	firstPlayer := rand.Intn(2)
	if firstPlayer == 1 {
		g.changePlayersPosition()
	}
}

func (g *Game) AddPlayer(name string) error {
	if len(g.PlayersOrder) != 0 && g.Players[g.PlayersOrder[0]] != nil && g.Players[g.PlayersOrder[0]].Name == name {
		return fmt.Errorf("name already exists")
	} else if len(g.PlayersOrder) >= 2 {
		return fmt.Errorf("game is full")
	} else {
		g.PlayersOrder = append(g.PlayersOrder, name)
		g.Players[name] = NewPlayer(name)

		return nil
	}
}

func (g *Game) IsGameReady() bool {
	return len(g.PlayersOrder) == 2 &&
		g.Players[g.PlayersOrder[0]] != nil &&
		g.Players[g.PlayersOrder[1]] != nil
}
