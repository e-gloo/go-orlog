package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	og "github.com/e-gloo/orlog/internal/orlog"
	"github.com/gorilla/websocket"
)

var joinableGames = sync.Map{}

type Packet struct {
	Command Command `json:"command"`
	Data    string  `json:"data"`
}

type gameManager struct {
	game        *og.Game
	player1Conn *websocket.Conn
	player2Conn *websocket.Conn
}

type CommandHandler struct {
	Manager *gameManager
	IsHost  bool
}

func (ch *CommandHandler) Handle(conn *websocket.Conn, packet *Packet) error {
	switch packet.Command {
	case CreateGame:
		return ch.handleCreateGame(conn)
	case JoinGame:
		return ch.handleJoinGame(conn, packet)
	case AddPlayer:
		return ch.handleAddPlayer(conn, packet)
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
		game:        game,
		player1Conn: conn,
	}
	joinableGames.Store(game.Uuid, manager)
	ch.IsHost = true
	ch.Manager = manager
	slog.Info("Game created", "uuid", game.Uuid)

	if err := SendPacket(conn, &Packet{Command: CommandOK, Data: game.Uuid}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	if err := SendPacket(conn, &Packet{Command: AddPlayer}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	return nil
}

func (ch *CommandHandler) handleJoinGame(conn *websocket.Conn, packet *Packet) error {
	slog.Info("Trying to join...", "uuid", packet.Data)
	value, ok := joinableGames.Load(packet.Data)
	if !ok {
		slog.Debug("Error joining game, uuid not found", "uuid", packet.Data)
		return nil
	}

	manager := value.(*gameManager)
	manager.player2Conn = conn
	ch.Manager = manager
	slog.Info("Joined game", "uuid", packet.Data)
	joinableGames.Delete(packet.Data)

	if err := SendPacket(conn, &Packet{Command: AddPlayer}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}

func (ch *CommandHandler) handleAddPlayer(conn *websocket.Conn, packet *Packet) error {
	if ch.IsHost {
		ch.Manager.game.SetPlayer1(packet.Data)
		slog.Info("Player 1 added", "name", packet.Data)
	} else {
		ch.Manager.game.SetPlayer2(packet.Data)
		slog.Info("Player 2 added", "name", packet.Data)
	}

	if err := SendPacket(conn, &Packet{Command: CommandOK}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	if ch.Manager.game.IsGameReady() {
		slog.Info("Game is starting...")
		if err := SendPacket(ch.Manager.player1Conn, &Packet{Command: GameStarting}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
		if err := SendPacket(ch.Manager.player2Conn, &Packet{Command: GameStarting}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
		return nil
	} else {
		slog.Info("Player 1 ready", "status", ch.Manager.game.Player1 != nil)
		slog.Info("Player 2 ready", "status", ch.Manager.game.Player2 != nil)
		return nil
	}
}

func (ch *CommandHandler) handleDefaultCase(conn *websocket.Conn, command Command) error {
	slog.Debug("Unknown command", "command", command)
	if err := SendPacket(conn, &Packet{Command: CommandError}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}

func SendPacket(conn *websocket.Conn, packet *Packet) error {
	newPacketBuffer := new(bytes.Buffer)
	err := json.NewEncoder(newPacketBuffer).Encode(packet)
	if err != nil {
		return fmt.Errorf("error encoding data: %w", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, newPacketBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("error writing message: %w", err)
	}
	return nil
}
