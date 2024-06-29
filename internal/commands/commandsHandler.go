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
		slog.Info("Creating new game")

		game, err := og.InitGame()

		if err != nil {
			err = fmt.Errorf("Error initializing game: %w", err)
			newErr := SendPacket(conn, &Packet{Command: CommandError, Data: err.Error()})
			if newErr != nil {
				err = newErr
			}
			return err
		}

		manager := &gameManager{
			game:        game,
			player1Conn: conn,
		}
		joinableGames.Store(game.Uuid, manager)
		ch.IsHost = true
		ch.Manager = manager
		slog.Info("Game created", "uuid", game.Uuid)

		err = SendPacket(conn, &Packet{Command: CommandOK, Data: game.Uuid})
		if err != nil {
			err = fmt.Errorf("Error sending packet: %w", err)
			return err
		}

		err = SendPacket(conn, &Packet{Command: AddPlayer})
		if err != nil {
			err = fmt.Errorf("Error sending packet: %w", err)
			return err
		}

	case JoinGame:
		slog.Info("Trying to join...", "uuid", packet.Data)
		value, ok := joinableGames.Load(packet.Data)
		if !ok {
			slog.Warn("Error joining game, uuid not found", "uuid", packet.Data)
			err := SendPacket(conn, &Packet{Command: CommandError, Data: "UUID not found"})
			if err != nil {
				err = fmt.Errorf("Error sending packet: %w", err)
				return err
			}
			return nil
		}
		manager := value.(*gameManager)
		manager.player2Conn = conn
		ch.Manager = manager
		slog.Info("Joined game", "uuid", packet.Data)
		joinableGames.Delete(packet.Data)

		err := SendPacket(conn, &Packet{Command: AddPlayer})
		if err != nil {
			err = fmt.Errorf("Error sending packet: %w", err)
			return err
		}

	case AddPlayer:
		if ch.IsHost {
			ch.Manager.game.SetPlayer1(packet.Data)
			slog.Info("Player 1 added", "name", packet.Data)
		} else {
			ch.Manager.game.SetPlayer2(packet.Data)
			slog.Info("Player 2 added", "name", packet.Data)
		}
		err := SendPacket(conn, &Packet{Command: CommandOK})
		if err != nil {
			err = fmt.Errorf("Error sending packet: %w", err)
			return err
		}

		if ch.Manager.game.IsGameReady() {
			slog.Info("Game is starting...")
			err = SendPacket(ch.Manager.player1Conn, &Packet{Command: GameStarting})
			if err != nil {
				err = fmt.Errorf("Error sending packet: %w", err)
				return err
			}
			err = SendPacket(ch.Manager.player2Conn, &Packet{Command: GameStarting})
			if err != nil {
				err = fmt.Errorf("Error sending packet: %w", err)
				return err
			}
		} else {
			slog.Info("Player 1 nil", "p1", ch.Manager.game.Player1 == nil)
			slog.Info("Player 2 nil", "p2", ch.Manager.game.Player2 == nil)
		}

	default:
		err := SendPacket(conn, &Packet{Command: CommandError, Data: "Command not found"})
		if err != nil {
			err = fmt.Errorf("Error sending packet: %w", err)
			return err
		}
	}
	return nil
}

func SendPacket(conn *websocket.Conn, packet *Packet) error {
	newPacketBuffer := new(bytes.Buffer)
	err := json.NewEncoder(newPacketBuffer).Encode(packet)
	if err != nil {
		slog.Error("Error encoding data", "err", err)
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, newPacketBuffer.Bytes())
	if err != nil {
		slog.Error("Error writing message", "err", err)
		return err
	}
	return nil
}
