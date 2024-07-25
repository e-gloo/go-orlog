package client

import (
	"encoding/json"
	"fmt"

	"net/url"
	"os"
	"os/signal"
	"time"

	g "github.com/e-gloo/orlog/internal/client/game"
	l "github.com/e-gloo/orlog/internal/client/lobby"
	c "github.com/e-gloo/orlog/internal/commands"
	"github.com/gorilla/websocket"
)

type State int

const (
	LobbyState State = iota
	GameState
)

type Client interface {
	Run(IOHandler) error
	CreateGame() error
	JoinGame(string) error
	AddPlayerName(string) error
	GetState() State
	GetLobby() *l.Lobby
	GetGame() *g.ClientGame
}

type client struct {
	conn  *websocket.Conn
	lobby *l.Lobby
	game  *g.ClientGame
	state State
}

func NewClient(serverAddr string) (Client, error) {
	u := &url.URL{Scheme: "ws", Host: serverAddr, Path: "/connect"}
	// slog.Info("connecting", "url", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %w", err)
	}

	lobby := &l.Lobby{ServerUrl: u.String()}

	return &client{conn: conn, lobby: lobby, state: LobbyState}, nil
}

func (cl *client) Run(ioh IOHandler) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ch := NewCommandHandler(cl.conn, cl.lobby, ioh)

	defer cl.conn.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := cl.conn.ReadMessage()
			if err != nil {
				// slog.Error("read", "err", err)
				break
			}

			packet := &c.Packet{}
			err = json.Unmarshal(message, packet)
			if err != nil {
				// slog.Error("Error unmarshalling packet", "err", err)
				return
			}

			// slog.Debug("New message", "packet", packet)
			err = ch.Handle(packet)
			if err != nil {
				// slog.Error("Error handling packet", "err", err)
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			return nil
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := cl.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				// slog.Warn("write close:", "error", err)
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

func (cl *client) GetState() State {
	return cl.state
}

func (cl *client) GetLobby() *l.Lobby {
	return cl.lobby
}

func (cl *client) GetGame() *g.ClientGame {
	return cl.game
}

func (cl *client) CreateGame() error {
	err := c.SendPacket(cl.conn, c.CreateGame, &c.CreateGameMessage{})

	if err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return err
}

func (cl *client) JoinGame(uuid string) error {
	err := c.SendPacket(cl.conn, c.JoinGame, &c.JoinGameMessage{Uuid: uuid})

	if err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return err
}

func (cl *client) AddPlayerName(name string) error {
	return nil
}
