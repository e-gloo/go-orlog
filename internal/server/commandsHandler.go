package server

import (
	"fmt"
	"log/slog"
	"slices"
	"sync"

	"github.com/e-gloo/orlog/internal/commands"
	og "github.com/e-gloo/orlog/internal/orlog"
	"github.com/gorilla/websocket"
)

var joinableGames = sync.Map{}

type gameManager struct {
	game    *og.Game
	players map[string]*websocket.Conn
}

type CommandHandler struct {
	manager          *gameManager
	isHost           bool
	expectedCommands []commands.Command
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		expectedCommands: []commands.Command{
			commands.CreateGame,
			commands.JoinGame,
		},
	}
}

func (ch *CommandHandler) Handle(conn *websocket.Conn, packet *commands.Packet) error {
	if !slices.Contains(ch.expectedCommands, packet.Command) {
		slog.Warn("Unexpected command", "command", packet.Command)
		if err := commands.SendPacket(conn, &commands.Packet{Command: commands.CommandError}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
		return nil
	}

	switch packet.Command {
	case commands.CreateGame:
		return ch.handleCreateGame(conn)
	case commands.JoinGame:
		return ch.handleJoinGame(conn, packet)
	case commands.AddPlayer:
		return ch.handleAddPlayer(conn, packet)
	case commands.CommandError:
		slog.Debug("Oops désolé :D")
		return nil
	default:
		return ch.handleDefaultCase(conn, packet.Command)
	}
}

func (ch *CommandHandler) handleCreateGame(conn *websocket.Conn) error {
	slog.Info("Creating new game")

	game, err := og.InitGame()
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

	if err := commands.SendPacket(conn, &commands.Packet{Command: commands.CommandOK, Data: game.Uuid}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	if err := commands.SendPacket(conn, &commands.Packet{Command: commands.AddPlayer}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	return nil
}

func (ch *CommandHandler) handleJoinGame(conn *websocket.Conn, packet *commands.Packet) error {
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

	if err := commands.SendPacket(conn, &commands.Packet{Command: commands.CommandOK, Data: packet.Data}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	if err := commands.SendPacket(conn, &commands.Packet{Command: commands.AddPlayer}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}

func (ch *CommandHandler) handleAddPlayer(conn *websocket.Conn, packet *commands.Packet) error {
	if ch.isHost {
		if err := ch.manager.game.SetPlayer1(packet.Data); err != nil {
			if err := commands.SendPacket(conn, &commands.Packet{Command: commands.CommandError, Data: err.Error()}); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}

			if err := commands.SendPacket(conn, &commands.Packet{Command: commands.AddPlayer}); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}
			return nil
		}

		ch.manager.players[packet.Data] = conn
		slog.Info("Player 1 added", "name", packet.Data)
	} else {
		if err := ch.manager.game.SetPlayer2(packet.Data); err != nil {
			if err := commands.SendPacket(conn, &commands.Packet{Command: commands.CommandError, Data: err.Error()}); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}

			if err := commands.SendPacket(conn, &commands.Packet{Command: commands.AddPlayer}); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}
			return nil
		}

		ch.manager.players[packet.Data] = conn
		slog.Info("Player 2 added", "name", packet.Data)
	}

	ch.expectedCommands = []commands.Command{}

	if err := commands.SendPacket(conn, &commands.Packet{Command: commands.CommandOK}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	if ch.manager.game.IsGameReady() {
		slog.Info("Game is starting...")
		for u := range ch.manager.players {
			if err := commands.SendPacket(ch.manager.players[u], &commands.Packet{Command: commands.GameStarting}); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}
		}
		return nil
	} else {
		slog.Debug("Player 1 ready", "status", ch.manager.game.Player1 != nil)
		slog.Debug("Player 2 ready", "status", ch.manager.game.Player2 != nil)
		return nil
	}
}

func (ch *CommandHandler) handleDefaultCase(conn *websocket.Conn, command commands.Command) error {
	slog.Debug("Unknown command", "command", command)
	if err := commands.SendPacket(conn, &commands.Packet{Command: commands.CommandError}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}
