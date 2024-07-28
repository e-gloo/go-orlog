package client_game

import (
	"fmt"

	cmn "github.com/e-gloo/orlog/internal/commons"
)

type ClientGame struct {
	MyUsername string
	Players    cmn.PlayerMap[*ClientPlayer]
	Dice       [6]ClientDie
	Gods       []ClientGod
}

func NewClientGame(
	playerUsername string,
	initGameDice []cmn.InitGameDie,
	initGameGods []cmn.InitGod,
	initPlayers cmn.PlayerMap[cmn.InitGamePlayer],
) *ClientGame {
	players := make(cmn.PlayerMap[*ClientPlayer], 2)

	for u, p := range initPlayers {
		players[u] = NewClientPlayer(p)
	}

	return &ClientGame{
		MyUsername: playerUsername,
		Players:    players,
		Gods:       mapGameGods(initGameGods),
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

func mapGameGods(initGods []cmn.InitGod) []ClientGod {
	var res []ClientGod

	for i := range initGods {
		res = append(res, ClientGod{
			Emoji:       initGods[i].Emoji,
			Name:        initGods[i].Name,
			Description: initGods[i].Description,
			Priority:    initGods[i].Priority,
			Levels:      mapGameGodPowers(initGods[i].Levels),
		})
	}

	return res
}

func mapGameGodPowers(initPowers [3]cmn.InitGodPower) [3]ClientGodPower {
	var res [3]ClientGodPower

	for i := range initPowers {
		res[i] = ClientGodPower{
			Description: initPowers[i].Description,
			TokenCost:   initPowers[i].TokenCost,
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
	res := "\n"

	opponent := cg.Players[cg.GetOpponentName()]
	player := cg.Players[cg.MyUsername]

	res += fmt.Sprintf(" * %s HP: %d t: %d\t[%s]\n", opponent.username, opponent.health, opponent.tokens, formatGodEmojis(cg.Gods, opponent.GetGods()))
	res += formatDice(cg.Dice, opponent.GetDice())
	res += "\n\n"
	res += formatDice(cg.Dice, player.GetDice())
	res += fmt.Sprintf(" * %s HP: %d t: %d\t[%s]\n", player.username, player.health, player.tokens, formatGodEmojis(cg.Gods, player.GetGods()))

	return res
}

func formatGodEmojis(gods []ClientGod, playerGods [3]int) string {
	return fmt.Sprintf("%s, %s, %s", gods[playerGods[0]].Emoji, gods[playerGods[1]].Emoji, gods[playerGods[2]].Emoji)
}

func formatDice(dice [6]ClientDie, state PlayerDice) string {
	res := ""
	for dieIdx, die := range dice {
		res += fmt.Sprintf(
			"%d %s",
			1+dieIdx,
			die.FormatDie(state[dieIdx]),
		)
	}
	res += "\n"
	return res
}
