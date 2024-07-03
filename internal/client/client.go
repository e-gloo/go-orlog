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

func ListenForServer(conn *websocket.Conn, ioh IOHandler) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ch := NewCommandHandler(ioh, conn)

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
