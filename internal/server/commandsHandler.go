package server

import (
	"fmt"
	"log/slog"
	"slices"
	"sync"

	"github.com/e-gloo/orlog/internal/commands"
	"github.com/gorilla/websocket"
)

var joinableGames = sync.Map{}

type gameManager struct {
	game    *ServerGame
	players map[string]*websocket.Conn
}

type CommandHandler struct {
	conn             *websocket.Conn
	manager          *gameManager
	isHost           bool
	expectedCommands []commands.Command
}

func NewCommandHandler(conn *websocket.Conn) *CommandHandler {
	return &CommandHandler{
		conn:   conn,
		isHost: false,
		expectedCommands: []commands.Command{
			commands.CreateGame,
			commands.JoinGame,
		},
	}
}

func (ch *CommandHandler) Handle(packet *commands.Packet) error {
	if !slices.Contains(ch.expectedCommands, packet.Command) {
		slog.Warn("Unexpected command", "command", packet.Command)
		if err := commands.SendPacket(ch.conn, &commands.Packet{Command: commands.CommandError}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
		return nil
	}

	switch packet.Command {
	case commands.CreateGame:
		return ch.handleCreateGame()
	case commands.JoinGame:
		return ch.handleJoinGame(packet)
	case commands.AddPlayer:
		return ch.handleAddPlayer(packet)
	case commands.CommandError:
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
		players: make(map[string]*websocket.Conn),
	}
	joinableGames.Store(game.Uuid, manager)
	ch.isHost = true
	ch.manager = manager
	ch.expectedCommands = []commands.Command{commands.AddPlayer}
	slog.Info("Game created", "uuid", game.Uuid)

	if err := commands.SendPacket(ch.conn, &commands.Packet{Command: commands.CommandOK, Data: game.Uuid}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	if err := commands.SendPacket(ch.conn, &commands.Packet{Command: commands.AddPlayer}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	return nil
}

func (ch *CommandHandler) handleJoinGame(packet *commands.Packet) error {
	slog.Info("Trying to join...", "uuid", packet.Data)
	value, ok := joinableGames.Load(packet.Data)
	if !ok {
		slog.Debug("Error joining game, uuid not found", "uuid", packet.Data)
		return nil
	}

	manager := value.(*gameManager)
	ch.manager = manager
	ch.expectedCommands = []commands.Command{commands.AddPlayer}
	slog.Info("Joined game", "uuid", packet.Data)
	joinableGames.Delete(packet.Data)

	if err := commands.SendPacket(ch.conn, &commands.Packet{Command: commands.CommandOK, Data: packet.Data}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	if err := commands.SendPacket(ch.conn, &commands.Packet{Command: commands.AddPlayer}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}

func (ch *CommandHandler) handleAddPlayer(packet *commands.Packet) error {

	if err := ch.manager.game.Data.AddPlayer(packet.Data); err != nil {
		if err := commands.SendPacket(ch.conn, &commands.Packet{Command: commands.CommandError, Data: err.Error()}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}

		if err := commands.SendPacket(ch.conn, &commands.Packet{Command: commands.AddPlayer}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
		return nil
	}

	ch.manager.players[packet.Data] = ch.conn

	if ch.isHost {
		slog.Info("Player 1 added", "name", packet.Data)
	} else {
		slog.Info("Player 2 added", "name", packet.Data)
	}

	ch.expectedCommands = []commands.Command{}

	if err := commands.SendPacket(ch.conn, &commands.Packet{Command: commands.CommandOK}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	if ch.manager.game.Data.IsGameReady() {
		ch.manager.game.Data.SelectFirstPlayer()
		slog.Info("Game is starting...")

		gameData, err := ch.manager.game.String()
		if err != nil {
			return fmt.Errorf("error serializing game data: %w", err)
		}

		for u := range ch.manager.players {
			if err := commands.SendPacket(ch.manager.players[u], &commands.Packet{Command: commands.GameStarting, Data: gameData}); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}
		}
	}
	return nil
}

func (ch *CommandHandler) handleDefaultCase(command commands.Command) error {
	slog.Debug("Unknown command", "command", command)
	if err := commands.SendPacket(ch.conn, &commands.Packet{Command: commands.CommandError}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}
