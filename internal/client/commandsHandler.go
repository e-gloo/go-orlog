package client

import (
	"fmt"
	"log/slog"

	"github.com/e-gloo/orlog/internal/commands"
	"github.com/gorilla/websocket"
)

type CommandHandler struct {
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{}
}

func (ch *CommandHandler) Handle(conn *websocket.Conn, packet *commands.Packet) error {
	switch packet.Command {
	case commands.AddPlayer:
		return ch.handleAddPlayer(conn)
	case commands.CommandOK:
		if packet.Data != "" {
			fmt.Printf("%s\n", packet.Data)
		}
		return nil
	case commands.CommandError:
		slog.Debug("Oops désolé :D")
		return nil
	default:
		return ch.handleDefaultCase(packet.Command)
	}
}

func (ch *CommandHandler) handleAddPlayer(conn *websocket.Conn) error {
	initPlayer(conn)

	return nil
}

func (ch *CommandHandler) handleDefaultCase(command commands.Command) error {
	slog.Debug("Unknown command", "command", command)

	return nil
}
