package client_game

import (
	"fmt"

	cmn "github.com/e-gloo/orlog/internal/commons"
)

type ClientGame struct {
	MyUsername string
	Players    cmn.PlayerMap[*ClientPlayer]
	Dice       [6]ClientDie
}

func NewClientGame(playerUsername string, initGameDice []cmn.InitGameDie, initPlayers cmn.PlayerMap[cmn.InitGamePlayer]) *ClientGame {
	players := make(cmn.PlayerMap[*ClientPlayer], 2)

	for u, p := range initPlayers {
		players[u] = NewClientPlayer(p)
	}

	return &ClientGame{
		MyUsername: playerUsername,
		Players:    players,
		Dice:       mapGameDice(initGameDice),
	}
}

func mapGameDice(initDice []cmn.InitGameDie) [6]ClientDie {
	var res [6]ClientDie

	for i := range initDice {
		res[i] = ClientDie{
			faces: mapGameDieFaces(initDice[i].Faces),
		}
	}

	return res
}

func mapGameDieFaces(initFaces []cmn.InitGameDieFace) [6]ClientFace {
	var res [6]ClientFace

	for i := range initFaces {
		res[i] = ClientFace{
			kind:  initFaces[i].Kind,
			magic: initFaces[i].Magic,
		}
	}

	return res
}

func (cg *ClientGame) UpdatePlayers(update cmn.PlayerMap[cmn.UpdateGamePlayer]) {
	for username, player := range update {
		cg.Players[username].update(player)
	}
}

func (cg *ClientGame) UpdatePlayersDice(update cmn.PlayerMap[cmn.DiceState]) {
	for username, dice := range update {
		cg.Players[username].updateDice(dice)
	}
}

func (cg *ClientGame) GetOpponentName() string {
	for name := range cg.Players {
		if name != cg.MyUsername {
			return name
		}
	}
	return ""

}

func (cg *ClientGame) FormatGame() string {
	res := ""

	opponent := cg.Players[cg.GetOpponentName()]
	player := cg.Players[cg.MyUsername]

	for _, p := range []*ClientPlayer{opponent, player} {
		res += fmt.Sprintf("%s HP: %d t: %d\n", p.username, p.health, p.tokens)
		for dieIdx, die := range cg.Dice {
			dieState := p.GetDice()[dieIdx]
			res += fmt.Sprintf(
				"%d %s",
				1+dieIdx,
				die.FormatDie(dieState),
			)
		}
		res += "\n"
	}

	return res
}
