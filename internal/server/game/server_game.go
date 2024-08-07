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
	Gods         []*God
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
		Gods:         InitGods(),
		Players:      make(cmn.PlayerMap[*ServerPlayer]),
		PlayersOrder: make([]string, 0, 2),
	}, nil
}

func (g *ServerGame) AddPlayer(conn *websocket.Conn, name string, godIndexes [3]int) error {
	if len(g.PlayersOrder) != 0 && g.Players[g.PlayersOrder[0]] != nil && g.Players[g.PlayersOrder[0]].GetUsername() == name {
		return fmt.Errorf("name already exists")
	} else if len(g.PlayersOrder) >= 2 {
		return fmt.Errorf("game is full")
	} else {
		for _, i := range godIndexes {
			if i < 0 || i >= len(g.Gods) || g.Gods[i] == nil {
				return fmt.Errorf("god not found: %d", i)
			}
		}

		g.PlayersOrder = append(g.PlayersOrder, name)
		g.Players[name] = NewServerPlayer(conn, name, godIndexes)

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

func (g *ServerGame) IsGameFinished() bool {
	for _, p := range g.Players {
		if p.GetHealth() <= 0 {
			return true
		}
	}
	return false
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

func (g *ServerGame) GetStartingDefinition() cmn.PlayerMap[cmn.InitGamePlayer] {
	players := make(cmn.PlayerMap[cmn.InitGamePlayer], 2)
	for u := range g.Players {
		players[u] = cmn.InitGamePlayer{
			Username:   u,
			Health:     g.Players[u].GetHealth(),
			GodIndexes: g.Players[u].GetGods(),
		}
	}
	return players
}

func (g *ServerGame) GetRollDiceState() cmn.PlayerMap[cmn.DiceState] {
	dice := make(cmn.PlayerMap[cmn.DiceState], len(g.Players))
	for u := range g.Players {
		dice[u] = make(cmn.DiceState, len(g.Players[u].GetDice()))
		for die := range g.Players[u].GetDice() {
			dice[u][die].Index = g.Players[u].GetDice()[die].GetFaceIndex()
			dice[u][die].Kept = g.Players[u].GetDice()[die].IsKept()
		}
	}
	return dice
}

func (g *ServerGame) GetPlayerRollDiceState(username string) cmn.DiceState {
	dice := make(cmn.DiceState, len(g.Players[username].GetDice()))
	for die := range g.Players[username].GetDice() {
		dice[die].Index = g.Players[username].GetDice()[die].GetFaceIndex()
		dice[die].Kept = g.Players[username].GetDice()[die].IsKept()
	}
	return dice
}

func (g *ServerGame) GetDiceDefinition() []cmn.InitGameDie {
	dice := make([]cmn.InitGameDie, 6)
	for i := 0; i < 6; i++ {
		dice[i].Faces = make([]cmn.InitGameDieFace, 6)
		for j := 0; j < 6; j++ {
			dice[i].Faces[j] = cmn.InitGameDieFace{
				Kind:  g.Dice[i].GetFaces()[j].GetKind(),
				Magic: g.Dice[i].GetFaces()[j].IsMagic(),
			}
		}
	}
	return dice
}

func (g *ServerGame) GetGodsDefinition() []cmn.InitGod {
	res := make([]cmn.InitGod, len(g.Gods))
	for i, god := range g.Gods {
		if god == nil {
			continue
		}

		res[i] = cmn.InitGod{
			Emoji:       god.Emoji,
			Name:        god.Name,
			Description: god.Description,
			Priority:    god.Priority,
			Levels: [3]cmn.InitGodPower{
				{
					Description: god.Levels[0].Description,
					TokenCost:   god.Levels[0].TokenCost,
				},
				{
					Description: god.Levels[1].Description,
					TokenCost:   god.Levels[1].TokenCost,
				},
				{
					Description: god.Levels[2].Description,
					TokenCost:   god.Levels[2].TokenCost,
				},
			},
		}
	}

	return res
}

func (g *ServerGame) ComputeRound() cmn.PlayerMap[cmn.UpdateGamePlayer] {
	p1res := cmn.UpdateGamePlayer{}
	p2res := cmn.UpdateGamePlayer{}

	g.activateGods(1, 2)

	// gain tokens
	p1res.TokensGained = g.Players[g.PlayersOrder[0]].GainTokens(g.Dice)
	p2res.TokensGained = g.Players[g.PlayersOrder[1]].GainTokens(g.Dice)

	g.activateGods(3)

	// damage phase
	p2res.ArrowDamageReceived, p2res.AxeDamageReceived = g.Players[g.PlayersOrder[0]].AttackPlayer(g.Dice, g.Players[g.PlayersOrder[1]])
	if g.Players[g.PlayersOrder[1]].GetHealth() > 0 {
		p1res.ArrowDamageReceived, p1res.AxeDamageReceived = g.Players[g.PlayersOrder[1]].AttackPlayer(g.Dice, g.Players[g.PlayersOrder[0]])

		g.activateGods(4)

		// thief phase
		p1res.TokensStolen = g.Players[g.PlayersOrder[0]].StealTokens(g.Dice, g.Players[g.PlayersOrder[1]])
		p2res.TokensStolen = g.Players[g.PlayersOrder[1]].StealTokens(g.Dice, g.Players[g.PlayersOrder[0]])

		g.activateGods(5, 6, 7)
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

func (g *ServerGame) activateGods(phases ...int) {
	for _, phase := range phases {
		for _, u := range g.PlayersOrder {
			godChoice := g.Players[u].GetGodChoice()
			if godChoice != nil && godChoice.index != -1 {
				god := g.Gods[godChoice.index]
				if god.Priority == phase {
					opponent := g.Players[g.GetOpponentName(u)]
					g.Players[u].activateGod(god, godChoice.level, g, opponent)
				}
			}
		}
	}
}
