package server

import (
	"fmt"
	"log/slog"
	"math"
	"os"
	"sync"

	c "github.com/e-gloo/orlog/internal/commands"
	cmn "github.com/e-gloo/orlog/internal/commons"
	g "github.com/e-gloo/orlog/internal/server/game"
	"github.com/gorilla/websocket"
)

var joinableGames = sync.Map{}

type CommandHandler struct {
	Conn     *websocket.Conn
	Username string
	game     *g.ServerGame
	isHost   bool
}

func NewCommandHandler(conn *websocket.Conn) *CommandHandler {
	return &CommandHandler{
		Conn:     conn,
		isHost:   false,
		Username: "Player",
	}
}

func (ch *CommandHandler) Handle(packet *c.Packet) error {
	// if !slices.Contains(ch.game.Players[ch.Username].GetExpectedCommands(), packet.Command) {
	// 	return ch.handleUnexpectedCommand(packet.Command)
	// }

	switch packet.Command {
	case c.CreateGame:
		return ch.handleCreateGame()
	case c.JoinGame:
		return ch.handleJoinGame(packet)
	case c.AddPlayer:
		return ch.handleAddPlayer(packet)
	case c.KeepDice:
		return ch.handleKeepDice(packet)
	default:
		return ch.handleDefaultCase(packet.Command)
	}
}

func (ch *CommandHandler) handleCreateGame() error {
	slog.Info("Creating new game")

	game, err := g.NewServerGame()
	if err != nil {
		return fmt.Errorf("error initializing game: %w", err)
	}

	joinableGames.Store(game.Uuid, game)
	ch.isHost = true
	ch.game = game
	// ch.ExpectedCommands = []c.Command{c.AddPlayer}
	slog.Info("Game created", "uuid", game.Uuid)

	if err := c.SendPacket(ch.Conn, c.CreatedOrJoined, &c.CreatedOrJoinedMessage{Uuid: game.Uuid}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	if err := c.SendPacket(ch.Conn, c.ConfigurePlayer, &c.ConfigurePlayerMessage{Gods: nil}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	return nil
}

func (ch *CommandHandler) handleJoinGame(packet *c.Packet) error {
	var joinGameMessage c.JoinGameMessage
	if err := c.ParsePacketData(packet, &joinGameMessage); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	slog.Info("Trying to join...", "uuid", joinGameMessage.Uuid)
	value, ok := joinableGames.Load(joinGameMessage.Uuid)
	if !ok {
		slog.Debug("Error joining game, uuid not found", "uuid", joinGameMessage.Uuid)
		if err := c.SendPacket(ch.Conn, c.CreateOrJoin, &c.CreateOrJoinMessage{Welcome: "Game not found, try again."}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
		return nil
	}

	game, ok := value.(*g.ServerGame)
	if !ok {
		return fmt.Errorf("could not retrieve a valid game")
	}
	ch.game = game

	// ch.ExpectedCommands = []c.Command{c.AddPlayer}
	slog.Info("Joined game", "uuid", joinGameMessage.Uuid)
	joinableGames.Delete(joinGameMessage.Uuid)

	if err := c.SendPacket(ch.Conn, c.CreatedOrJoined, &c.CreatedOrJoinedMessage{Uuid: ch.game.Uuid}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	if err := c.SendPacket(ch.Conn, c.ConfigurePlayer, &c.ConfigurePlayerMessage{Gods: nil}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}

func (ch *CommandHandler) handleAddPlayer(packet *c.Packet) error {
	var message c.AddPlayerMessage
	if err := c.ParsePacketData(packet, &message); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	slog.Debug("addPlayer", "game", ch.game, "username", message.Username)

	// TODO: add message.GodIndexes to AddPlayer
	if err := ch.game.AddPlayer(ch.Conn, message.Username); err != nil {
		if err := c.SendPacket(ch.Conn, c.CommandError, &c.CommandErrorMessage{Reason: err.Error()}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}

		if err := c.SendPacket(ch.Conn, c.ConfigurePlayer, &c.ConfigurePlayerMessage{Gods: nil}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
		return nil
	}

	ch.Username = message.Username

	if ch.isHost {
		slog.Info("Player 1 added", "name", message.Username)
	} else {
		slog.Info("Player 2 added", "name", message.Username)
	}

	// ch.ExpectedCommands = []c.Command{}

	if ch.game.IsGameReady() {
		if err := ch.handleStartingGame(); err != nil {
			return fmt.Errorf("error starting the game: %w", err)
		}
	}
	return nil
}

func (ch *CommandHandler) handleStartingGame() error {
	ch.game.SelectFirstPlayer()
	slog.Info("Game is starting...")

	var gameStartingMessage c.GameStartingMessage
	gameStartingMessage.Players = make(cmn.PlayerMap[cmn.InitGamePlayer], 2)
	for u := range ch.game.Players {
		gameStartingMessage.Players[u] = cmn.InitGamePlayer{
			Username:   u,
			Health:     ch.game.Players[u].GetHealth(),
			GodIndexes: [3]int{0, 0, 0},
		}
	}
	gameStartingMessage.Dice = make([]cmn.InitGameDie, 6)
	for i := 0; i < 6; i++ {
		gameStartingMessage.Dice[i].Faces = make([]cmn.InitGameDieFace, 6)
		for j := 0; j < 6; j++ {
			gameStartingMessage.Dice[i].Faces[j] = cmn.InitGameDieFace{
				Kind:  ch.game.Dice[i].GetFaces()[j].GetKind(),
				Magic: ch.game.Dice[i].GetFaces()[j].IsMagic(),
			}
		}
	}

	// Send every player the data to init the game
	for u := range ch.game.Players {
		gameStartingMessage.YourUsername = u
		if err := c.SendPacket(ch.game.Players[u].Conn, c.GameStarting, &gameStartingMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
	}

	// do the first roll
	firstUsername := ch.game.PlayersOrder[0]
	// secondUsername := ch.game.PlayersOrder[1]
	ch.game.Players[firstUsername].RollDice()

	var selectDiceMessage c.SelectDiceMessage
	// TODO: fill this selectDiceMessage struct with data from the game.
	selectDiceMessage.Turn = int(math.Ceil(float64(ch.game.Rolls) / 2))
	selectDiceMessage.Players = make(cmn.PlayerMap[cmn.DiceState], len(ch.game.Players))
	for u := range ch.game.Players {
		selectDiceMessage.Players[u] = make(cmn.DiceState, len(ch.game.Players[u].GetDice()))
		for die := range ch.game.Players[u].GetDice() {
			selectDiceMessage.Players[u][die].Index = ch.game.Players[u].GetDice()[die].GetFaceIndex()
			selectDiceMessage.Players[u][die].Kept = ch.game.Players[u].GetDice()[die].IsKept()
		}
	}

	if err := c.SendPacket(ch.game.Players[firstUsername].Conn, c.SelectDice, &selectDiceMessage); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	// if err := c.SendPacket(ch.game.Players[secondUsername].Conn, c.OpponentDiceRoll, &selectDiceMessage); err != nil {
	// 	return fmt.Errorf("error sending packet: %w", err)
	// }

	// ch.players[firstUsername].ExpectedCommands = []c.Command{c.KeepDice}
	// ch.players[secondUsername].ExpectedCommands = []c.Command{}

	ch.game.Rolls++

	return nil
}

func (ch *CommandHandler) handleKeepDice(packet *c.Packet) error {
	var message c.KeepDiceMessage
	if err := c.ParsePacketData(packet, &message); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	for dice_idx, dice_kept := range message.Kept {
		if dice_kept {
			ch.game.Players[ch.Username].GetDice()[dice_idx].Keep()
		} else {
			ch.game.Players[ch.Username].GetDice()[dice_idx].Unkeep()
		}
	}

	if ch.game.Rolls >= 4 {
		for u := range ch.game.Players {
			ch.game.Players[u].RollDice()
		}

		slog.Info("Round is over, computing ...")
		// ch.game.ComputeRound()
		ch.game.Rolls = 0

		// TODO: new round ... we need to find a way to hydrate both clients after computation

		// gameData, err := ch.game.String()
		// if err != nil {
		// 	return fmt.Errorf("error serializing game data: %w", err)
		// }

		// // Send every player the update of the game
		// for u := range ch.players {
		// 	if err := c.SendPacket(ch.players[u].Conn, &c.Packet{Command: c.GameInfo, Data: gameData}); err != nil {
		// 		return fmt.Errorf("error sending packet: %w", err)
		// 	}
		// }

		if ch.game.Players[ch.game.PlayersOrder[1]].GetHealth() <= 0 {
			// P1 won
			slog.Info("Congratulations P1, you won ! :)")
			os.Exit(1)
		} else if ch.game.Players[ch.game.PlayersOrder[0]].GetHealth() <= 0 {
			// P2 won
			slog.Info("Congratulations P2, you won ! :)")
			os.Exit(2)
		} else {
			ch.game.ChangePlayersPosition()

			firstUsername := ch.game.PlayersOrder[0]
			// secondUsername := ch.game.PlayersOrder[1]
			ch.game.Players[firstUsername].RollDice()

			var selectDiceMessage c.SelectDiceMessage
			// TODO: fill this selectDiceMessage struct with data from the game.
			selectDiceMessage.Turn = int(math.Ceil(float64(ch.game.Rolls) / 2))
			selectDiceMessage.Players = make(cmn.PlayerMap[cmn.DiceState], len(ch.game.Players))
			for u := range ch.game.Players {
				selectDiceMessage.Players[u] = make(cmn.DiceState, len(ch.game.Players[u].GetDice()))
				for die := range ch.game.Players[u].GetDice() {
					selectDiceMessage.Players[u][die].Index = ch.game.Players[u].GetDice()[die].GetFaceIndex()
					selectDiceMessage.Players[u][die].Kept = ch.game.Players[u].GetDice()[die].IsKept()
				}
			}

			if err := c.SendPacket(ch.game.Players[firstUsername].Conn, c.SelectDice, &selectDiceMessage); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}

			// if err := c.SendPacket(ch.game.Players[secondUsername].Conn, c.OpponentDiceRoll, &selectDiceMessage); err != nil {
			// 	return fmt.Errorf("error sending packet: %w", err)
			// }

			// ch.players[firstUsername].ExpectedCommands = []c.Command{c.KeepDice}
			// ch.players[secondUsername].ExpectedCommands = []c.Command{}

			ch.game.Rolls++
		}
	} else {
		otherUsername := ch.game.GetOpponentName(ch.Username)

		ch.game.Players[otherUsername].RollDice()

		var selectDiceMessage c.SelectDiceMessage
		// TODO: fill this selectDiceMessage struct with data from the game.
		selectDiceMessage.Turn = int(math.Ceil(float64(ch.game.Rolls) / 2))
		selectDiceMessage.Players = make(cmn.PlayerMap[cmn.DiceState], len(ch.game.Players))
		for u := range ch.game.Players {
			selectDiceMessage.Players[u] = make(cmn.DiceState, len(ch.game.Players[u].GetDice()))
			for die := range ch.game.Players[u].GetDice() {
				selectDiceMessage.Players[u][die].Index = ch.game.Players[u].GetDice()[die].GetFaceIndex()
				selectDiceMessage.Players[u][die].Kept = ch.game.Players[u].GetDice()[die].IsKept()
			}
		}

		if err := c.SendPacket(ch.game.Players[otherUsername].Conn, c.SelectDice, &selectDiceMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}

		// ch.players[ch.Username].ExpectedCommands = []c.Command{}
		// ch.players[otherUsername].ExpectedCommands = []c.Command{c.KeepDice}

		ch.game.Rolls++
	}

	return nil
}

func (ch *CommandHandler) handleUnexpectedCommand(command c.Command) error {
	slog.Warn("Unexpected command", "command", command)
	if err := c.SendPacket(ch.Conn, c.CommandError, &c.CommandErrorMessage{Reason: "unexpected command"}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}

func (ch *CommandHandler) handleDefaultCase(command c.Command) error {
	slog.Debug("Unknown command", "command", command)
	if err := c.SendPacket(ch.Conn, c.CommandError, &c.CommandErrorMessage{Reason: "unknown command"}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}
