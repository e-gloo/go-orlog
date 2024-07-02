package client

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/e-gloo/orlog/internal/commands"
	"github.com/gorilla/websocket"
)

func NewClient(url *url.URL) (*websocket.Conn, error) {
	slog.Info("connecting", "url", url.String())

	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %w", err)
	}

	// defer conn.Close()

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

	// TODO: player.ChooseGod()

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

func StartGame(conn *websocket.Conn) error {
	err := joinOrCreateGame(conn)
	if err != nil {
		return err
	}

	player, err := initPlayer(conn)
	if err != nil {
		return err
	}

	slog.Info("Welcome to the game", "player", player.Data.Name)

	return nil
}
