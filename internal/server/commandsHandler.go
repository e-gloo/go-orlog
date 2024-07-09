package server

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sync"

	c "github.com/e-gloo/orlog/internal/commands"
	"github.com/gorilla/websocket"
)

var joinableGames = sync.Map{}

type gameManager struct {
	game    *ServerGame
	players map[string]*CommandHandler
}

type CommandHandler struct {
	Conn             *websocket.Conn
	Username         string
	manager          *gameManager
	isHost           bool
	ExpectedCommands []c.Command
}

func NewCommandHandler(conn *websocket.Conn) *CommandHandler {
	return &CommandHandler{
		Conn:     conn,
		isHost:   false,
		Username: "Player",
		ExpectedCommands: []c.Command{
			c.CreateGame,
			c.JoinGame,
		},
	}
}

func (ch *CommandHandler) Handle(packet *c.Packet) error {
	if !slices.Contains(ch.ExpectedCommands, packet.Command) {
		return ch.handleUnexpectedCommand(packet.Command)
	}

	switch packet.Command {
	case c.CreateGame:
		return ch.handleCreateGame()
	case c.JoinGame:
		return ch.handleJoinGame(packet)
	case c.AddPlayer:
		return ch.handleAddPlayer(packet)
	case c.KeepDice:
		return ch.handleKeepDice(packet)
	case c.CommandError:
		slog.Debug("Oops désolé :D")
		return nil
	default:
		return ch.handleDefaultCase(packet.Command)
	}
}

func (ch *CommandHandler) handleCreateGame() error {
	slog.Info("Creating new game")

	game, err := NewServerGame()
	if err != nil {
		return fmt.Errorf("error initializing game: %w", err)
	}

	manager := &gameManager{
		game:    game,
		players: make(map[string]*CommandHandler),
	}
	joinableGames.Store(game.Uuid, manager)
	ch.isHost = true
	ch.manager = manager
	ch.ExpectedCommands = []c.Command{c.AddPlayer}
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
		return nil
	}

	manager := value.(*gameManager)
	ch.manager = manager
	ch.ExpectedCommands = []c.Command{c.AddPlayer}
	slog.Info("Joined game", "uuid", joinGameMessage.Uuid)
	joinableGames.Delete(joinGameMessage.Uuid)

	if err := c.SendPacket(ch.Conn, c.CreatedOrJoined, &c.CreatedOrJoinedMessage{Uuid: ch.manager.game.Uuid}); err != nil {
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

	// TODO: add message.GodIndexes to AddPlayer
	if err := ch.manager.game.Data.AddPlayer(message.Username); err != nil {
		if err := c.SendPacket(ch.Conn, c.CommandError, &c.CommandErrorMessage{Reason: err.Error()}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}

		if err := c.SendPacket(ch.Conn, c.ConfigurePlayer, &c.ConfigurePlayerMessage{Gods: nil}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
		return nil
	}

	ch.manager.players[message.Username] = ch
	ch.Username = message.Username

	if ch.isHost {
		slog.Info("Player 1 added", "name", message.Username)
	} else {
		slog.Info("Player 2 added", "name", message.Username)
	}

	ch.ExpectedCommands = []c.Command{}

	if ch.manager.game.Data.IsGameReady() {
		if err := ch.handleStartingGame(); err != nil {
			return fmt.Errorf("error starting the game: %w", err)
		}
	}
	return nil
}

func (ch *CommandHandler) handleStartingGame() error {
	ch.manager.game.Data.SelectFirstPlayer()
	slog.Info("Game is starting...")

	var gameStartingMessage c.GameStartingMessage
	// TODO: fill this gameStartingMessage struct with data from the game.
	// FIXME: how to identify p1 from p2 ? only based on isHost ?

	// Send every player the data to init the game
	for u := range ch.manager.players {
		if err := c.SendPacket(ch.manager.players[u].Conn, c.GameStarting, &gameStartingMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
	}

	// do the first roll
	firstUsername := ch.manager.game.Data.PlayersOrder[0]
	secondUsername := ch.manager.game.Data.PlayersOrder[1]
	ch.manager.game.Data.Players[firstUsername].RollDice()

	var selectDiceMessage c.SelectDiceMessage
	// TODO: fill this selectDiceMessage struct with data from the game.
	// FIXME: how to identify p1 from p2 ? only based on isHost ?

	if err := c.SendPacket(ch.manager.players[firstUsername].Conn, c.SelectDice, &selectDiceMessage); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	ch.manager.players[firstUsername].ExpectedCommands = []c.Command{c.KeepDice}
	ch.manager.players[secondUsername].ExpectedCommands = []c.Command{}

	ch.manager.game.Rolls++

	return nil
}

func (ch *CommandHandler) handleKeepDice(packet *c.Packet) error {
	var message c.KeepDiceMessage
	if err := c.ParsePacketData(packet, &message); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	for dice_idx, dice_kept := range message.Kept {
		ch.manager.game.Data.Players[ch.Username].Dice[dice_idx].Kept = dice_kept
	}

	if ch.manager.game.Rolls >= 4 {
		for u := range ch.manager.players {
			ch.manager.game.Data.Players[u].RollDice()
		}

		ch.manager.game.Data.ComputeRound()
		ch.manager.game.Rolls = 0

		// TODO: new round ... we need to find a way to hydrate both clients after computation

		// gameData, err := ch.manager.game.String()
		// if err != nil {
		// 	return fmt.Errorf("error serializing game data: %w", err)
		// }

		// // Send every player the update of the game
		// for u := range ch.manager.players {
		// 	if err := c.SendPacket(ch.manager.players[u].Conn, &c.Packet{Command: c.GameInfo, Data: gameData}); err != nil {
		// 		return fmt.Errorf("error sending packet: %w", err)
		// 	}
		// }

		if ch.manager.game.Data.Players[ch.manager.game.Data.PlayersOrder[1]].Health <= 0 {
			// P1 won
			slog.Info("Congratulations P1, you won ! :)")
			os.Exit(1)
		} else if ch.manager.game.Data.Players[ch.manager.game.Data.PlayersOrder[0]].Health <= 0 {
			// P2 won
			slog.Info("Congratulations P2, you won ! :)")
			os.Exit(2)
		} else {
			ch.manager.game.Data.ChangePlayersPosition()

			firstUsername := ch.manager.game.Data.PlayersOrder[0]
			secondUsername := ch.manager.game.Data.PlayersOrder[1]
			ch.manager.game.Data.Players[firstUsername].RollDice()

			var selectDiceMessage c.SelectDiceMessage
			// TODO: fill this selectDiceMessage struct with data from the game.
			// FIXME: how to identify p1 from p2 ? only based on isHost ?

			if err := c.SendPacket(ch.manager.players[firstUsername].Conn, c.SelectDice, &selectDiceMessage); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}

			ch.manager.players[firstUsername].ExpectedCommands = []c.Command{c.KeepDice}
			ch.manager.players[secondUsername].ExpectedCommands = []c.Command{}

			ch.manager.game.Rolls++
		}
	} else {
		otherUsername := ch.manager.game.Data.PlayersOrder[slices.IndexFunc(ch.manager.game.Data.PlayersOrder, func(p string) bool {
			return p != ch.Username
		})]

		ch.manager.game.Data.Players[otherUsername].RollDice()

		var selectDiceMessage c.SelectDiceMessage
		// TODO: fill this selectDiceMessage struct with data from the game.
		// FIXME: how to identify p1 from p2 ? only based on isHost ?

		// Send otherPlayer the gameData after its roll.
		if err := c.SendPacket(ch.manager.players[otherUsername].Conn, c.SelectDice, &selectDiceMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}

		ch.manager.players[ch.Username].ExpectedCommands = []c.Command{}
		ch.manager.players[otherUsername].ExpectedCommands = []c.Command{c.KeepDice}

		ch.manager.game.Rolls++
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
