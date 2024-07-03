package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/e-gloo/orlog/internal/commands"
	"github.com/gorilla/websocket"
)

func NewClient(url *url.URL) (*websocket.Conn, error) {
	slog.Info("connecting", "url", url.String())

	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %w", err)
	}

	return conn, nil
}

func initPlayer(conn *websocket.Conn) (*ClientPlayer, error) {
	username := "Player"
	fmt.Println("Enter your name : ")
	_, err := fmt.Scanln(&username)
	if err != nil && err.Error() != "unexpected newline" {
		return nil, err
	}

	player := NewClientPlayer(username)

	if err := commands.SendPacket(conn, &commands.Packet{
		Command: commands.AddPlayer,
		Data:    username,
	}); err != nil {
		return nil, fmt.Errorf("error sending packet: %w", err)
	}

	return player, nil
}

func joinOrCreateGame(conn *websocket.Conn) error {
	gameUuid := ""
	fmt.Println("Enter the game UUID (empty for new): ")
	_, err := fmt.Scanln(&gameUuid)
	if err != nil && err.Error() != "unexpected newline" {
		return err
	}

	var command commands.Command
	if gameUuid != "" {
		command = commands.JoinGame
	} else {
		command = commands.CreateGame
	}

	if err := commands.SendPacket(conn, &commands.Packet{
		Command: command,
		Data:    gameUuid,
	}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	return nil
}

func ListenForServer(conn *websocket.Conn) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ch := NewCommandHandler()

	defer conn.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				slog.Error("read", "err", err)
				break
			}

			packet := &commands.Packet{}
			err = json.Unmarshal(message, packet)
			if err != nil {
				slog.Error("Error unmarshalling packet", "err", err)
				return
			}

			slog.Debug("New message", "packet", packet)
			ch.Handle(conn, packet)
		}
	}()

	joinOrCreateGame(conn)

	for {
		select {
		case <-done:
			return nil
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				slog.Warn("write close:", "error", err)
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}
