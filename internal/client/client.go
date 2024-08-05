package client

import (
	"encoding/json"
	"fmt"

	"net/url"
	"os"
	"os/signal"
	"time"

	g "github.com/e-gloo/orlog/internal/client/game"
	"github.com/e-gloo/orlog/internal/commands"
	// "github.com/e-gloo/orlog/internal/pkg/logging"
	"github.com/gorilla/websocket"
)

type State int

const (
	LobbyState State = iota
	GameState
)

type Phase int

const (
	CreateOrJoin Phase = iota + 1
	ConfigPlayer
	GameStarting
	RollDice
	WaitingDiceRoll
	DiceRoll
	PickDice
	WaitingDicePick
)

type Client interface {
	Run(IOHandler) error
	ServerUrl() string
	CreateGame() error
	JoinGame(string) error
	GameUuid() string
	AddPlayerName(string) error
	GetGameGods() []g.ClientGod
	GetMyGods() [3]int
	GetOpponentGods() [3]int
	RollDice() error
	GetGameDice() [6]g.ClientDie
	GetMe() *g.ClientPlayer
	GetMyDice() g.PlayerDice
	ToggleDieState(int)
	GetOpponent() *g.ClientPlayer
	GetOpponentDice() g.PlayerDice
	KeepDice() error
	Error() string
}

type client struct {
	conn       *websocket.Conn
	game       *g.ClientGame
	state      State
	phase      Phase
	gameUuid   string
	playerName string
	err        string
}

func NewClient(serverAddr string) (Client, error) {
	u := &url.URL{Scheme: "ws", Host: serverAddr, Path: "/connect"}
	// slog.Info("connecting", "url", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %w", err)
	}

	return &client{conn: conn, state: LobbyState, phase: CreateOrJoin}, nil
}

func (cl *client) Run(ioh IOHandler) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ch := NewCommandHandler(cl.conn, cl, ioh)

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

			packet := &commands.Packet{}
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

func (cl *client) ServerUrl() string {
	return cl.conn.RemoteAddr().String()
}

func (cl *client) CreateGame() error {
	err := commands.SendPacket(cl.conn, commands.CreateGame, &commands.CreateGameMessage{})

	if err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return err
}

func (cl *client) JoinGame(uuid string) error {
	err := commands.SendPacket(cl.conn, commands.JoinGame, &commands.JoinGameMessage{Uuid: uuid})

	if err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return err
}

func (cl *client) GameUuid() string {
	return cl.gameUuid
}

func (cl *client) AddPlayerName(name string) error {
	err := commands.SendPacket(cl.conn, commands.AddPlayer, &commands.AddPlayerMessage{Username: name, GodIndexes: [3]int{0, 1, 2}})
	if err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return err
}

func (cl *client) GetGame() *g.ClientGame {
	return cl.game
}

func (cl *client) GetGameGods() []g.ClientGod {
	return cl.game.Gods
}

func (cl *client) GetMyGods() [3]int {
	me := cl.GetMe()
	return me.GetGods()
}

func (cl *client) ToggleDieState(idx int) {
	me := cl.GetMe()
	die := me.GetDice()[idx]
	die.SetKept(!die.IsKept())
}

func (cl *client) GetOpponentGods() [3]int {
	opponent := cl.GetOpponent()
	return opponent.GetGods()
}

func (cl *client) RollDice() error {
	err := commands.SendPacket(cl.conn, commands.RollDice, &commands.RollDiceMessage{})
	if err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return err
}

func (cl *client) GetGameDice() [6]g.ClientDie {
	return cl.game.Dice
}

func (cl *client) GetMe() *g.ClientPlayer {
	me := cl.game.MyUsername
	return cl.game.Players[me]
}

func (cl *client) GetMyDice() g.PlayerDice {
	me := cl.GetMe()
	return me.GetDice()
}

func (cl *client) GetOpponent() *g.ClientPlayer {
	me := cl.game.MyUsername
	for username := range cl.game.Players {
		if username != me {
			return cl.game.Players[username]
		}
	}
	return nil
}

func (cl *client) GetOpponentDice() g.PlayerDice {
	opponent := cl.GetOpponent()
	return opponent.GetDice()
}

func (cl *client) KeepDice() error {
	var keepDiceMessage commands.KeepDiceMessage
	dice := cl.GetMe().GetDice()
	
	for idx := range dice {
		keepDiceMessage.Kept[idx] = dice[idx].IsKept()
	}

	if err := commands.SendPacket(cl.conn, commands.KeepDice, &keepDiceMessage); err != nil {
		return fmt.Errorf("Error sending packet: %w", err)
	}
	return nil
}

func (cl *client) Error() string {
	return cl.err
}
