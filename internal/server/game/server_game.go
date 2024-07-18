package server_game

import (
	"fmt"
	"log/slog"
	"math/rand"

	cmn "github.com/e-gloo/orlog/internal/commons"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ServerGame struct {
	Uuid         string
	count        int
	turn         int
	Rolls        int
	Dice         [6]ServerDie
	Players      cmn.PlayerMap[*ServerPlayer]
	PlayersOrder []string
}

func NewServerGame() (*ServerGame, error) {
	newuuid, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("error generating uuid: %w", err)
	}

	return &ServerGame{
		Uuid:         newuuid.String(),
		count:        0,
		turn:         0,
		Rolls:        0,
		Dice:         InitDice(),
		Players:      make(cmn.PlayerMap[*ServerPlayer]),
		PlayersOrder: make([]string, 0, 2),
	}, nil
}

func (g *ServerGame) AddPlayer(conn *websocket.Conn, name string) error {
	if len(g.PlayersOrder) != 0 && g.Players[g.PlayersOrder[0]] != nil && g.Players[g.PlayersOrder[0]].GetUsername() == name {
		return fmt.Errorf("name already exists")
	} else if len(g.PlayersOrder) >= 2 {
		return fmt.Errorf("game is full")
	} else {
		g.PlayersOrder = append(g.PlayersOrder, name)
		g.Players[name] = NewServerPlayer(conn, name)

		return nil
	}
}

func (g *ServerGame) GetTurn() int {
	return g.turn
}

func (g *ServerGame) IsGameReady() bool {
	return len(g.PlayersOrder) == 2 &&
		g.Players[g.PlayersOrder[0]] != nil &&
		g.Players[g.PlayersOrder[1]] != nil
}

func (g *ServerGame) ChangePlayersPosition() {
	tmp := g.PlayersOrder[0]
	g.PlayersOrder[0] = g.PlayersOrder[1]
	g.PlayersOrder[1] = tmp
}

func (g *ServerGame) SelectFirstPlayer() {
	firstPlayer := rand.Intn(2)
	if firstPlayer == 1 {
		g.ChangePlayersPosition()
	}
}

func (g *ServerGame) Restart() {
	for _, p := range g.Players {
		p.Reset()
	}
	g.count++
	g.turn = 0
	g.Rolls = 0
}

func (g *ServerGame) GetOpponentName(you string) string {
	for name := range g.Players {
		if name != you {
			return name
		}
	}
	return ""
}

func (g *ServerGame) ComputeRound() cmn.PlayerMap[cmn.UpdateGamePlayer] {
	p1res := cmn.UpdateGamePlayer{}
	p2res := cmn.UpdateGamePlayer{}

	// gain tokens
	p1res.TokensGained = g.Players[g.PlayersOrder[0]].GainTokens(g.Dice)
	p2res.TokensGained = g.Players[g.PlayersOrder[1]].GainTokens(g.Dice)

	// damage phase
	p2res.ArrowDamageReceived, p2res.AxeDamageReceived = g.Players[g.PlayersOrder[0]].AttackPlayer(g.Dice, g.Players[g.PlayersOrder[1]])
	if g.Players[g.PlayersOrder[1]].GetHealth() > 0 {
		p1res.ArrowDamageReceived, p1res.AxeDamageReceived = g.Players[g.PlayersOrder[1]].AttackPlayer(g.Dice, g.Players[g.PlayersOrder[0]])

		// thief phase
		p1res.TokensStolen = g.Players[g.PlayersOrder[0]].StealTokens(g.Dice, g.Players[g.PlayersOrder[1]])
		p2res.TokensStolen = g.Players[g.PlayersOrder[1]].StealTokens(g.Dice, g.Players[g.PlayersOrder[0]])
	}

	g.Players[g.PlayersOrder[0]].ResetDice()
	g.Players[g.PlayersOrder[1]].ResetDice()

	g.Rolls = 0
	g.turn++

	slog.Debug(
		fmt.Sprintf(
			"Turn %d",
			g.turn,
		),
		g.Players[g.PlayersOrder[0]].username,
		fmt.Sprintf(
			"%dHP, %dTK",
			g.Players[g.PlayersOrder[0]].health,
			g.Players[g.PlayersOrder[0]].tokens,
		),
		g.Players[g.PlayersOrder[1]].username,
		fmt.Sprintf(
			"%dHP, %dTK",
			g.Players[g.PlayersOrder[1]].health,
			g.Players[g.PlayersOrder[1]].tokens,
		),
	)

	p1res.Health = g.Players[g.PlayersOrder[0]].GetHealth()
	p1res.Tokens = g.Players[g.PlayersOrder[0]].GetTokens()
	p2res.Health = g.Players[g.PlayersOrder[1]].GetHealth()
	p2res.Tokens = g.Players[g.PlayersOrder[1]].GetTokens()

	return cmn.PlayerMap[cmn.UpdateGamePlayer]{
		g.PlayersOrder[0]: p1res,
		g.PlayersOrder[1]: p2res,
	}
}
