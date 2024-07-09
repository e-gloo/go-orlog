package client

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	c "github.com/e-gloo/orlog/internal/commands"
	"github.com/gorilla/websocket"
)

type CommandHandler struct {
	ioh        IOHandler
	conn       *websocket.Conn
	game       *ClientGame
	myUsername string
}

func NewCommandHandler(ioh IOHandler, conn *websocket.Conn) *CommandHandler {
	return &CommandHandler{
		ioh:        ioh,
		conn:       conn,
		myUsername: "Player",
	}
}

func (ch *CommandHandler) Handle(conn *websocket.Conn, packet *c.Packet) error {
	switch packet.Command {
	case c.CreateOrJoin:
		return ch.handleCreateOrJoin(packet)
	case c.CreatedOrJoined:
		return ch.handleCreatedOrJoined(packet)
	case c.ConfigurePlayer:
		return ch.handleConfigurePlayer(packet)
	case c.GameStarting:
		return ch.handleGameStarting(packet)
	case c.SelectDice:
		return ch.handleSelectDice(packet)
	case c.CommandError:
		return ch.handleErrorCommand(packet.Command)
	default:
		return ch.handleDefaultCase(packet.Command)
	}
}

func (ch *CommandHandler) handleCreateOrJoin(packet *c.Packet) error {
	var createOrJoinMessage c.CreateOrJoinMessage
	if err := c.ParsePacketData(packet, &createOrJoinMessage); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	ch.ioh.DisplayMessage(createOrJoinMessage.Welcome)

	ch.ioh.DisplayMessage("Enter the game UUID (empty for new): ")
	gameUuid := ""
	err := ch.ioh.ReadInput(&gameUuid)
	if err != nil {
		return err
	}

	if gameUuid == "" {
		err = c.SendPacket(ch.conn, c.CreateGame, &c.CreateGameMessage{})
	} else {
		err = c.SendPacket(ch.conn, c.JoinGame, &c.JoinGameMessage{Uuid: gameUuid})
	}

	if err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	return nil
}

func (ch *CommandHandler) handleCreatedOrJoined(packet *c.Packet) error {
	var createdOrJoinedMessage c.CreatedOrJoinedMessage
	if err := c.ParsePacketData(packet, &createdOrJoinedMessage); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	ch.ioh.DisplayMessage("You have joined the game with UUID: " + createdOrJoinedMessage.Uuid)

	return nil
}

func (ch *CommandHandler) handleConfigurePlayer(packet *c.Packet) error {
	// var configurePlayerMessage c.ConfigurePlayerMessage
	// if err := c.ParsePacketData(packet, &configurePlayerMessage); err != nil {
	// 	return fmt.Errorf("error parsing packet data: %w", err)
	// }

	ch.ioh.DisplayMessage("Enter your name : ")

	err := ch.ioh.ReadInput(&ch.myUsername)
	if err != nil {
		return err
	}

	if err = c.SendPacket(ch.conn, c.AddPlayer, &c.AddPlayerMessage{Username: ch.myUsername, GodIndexes: [3]int{0, 0, 0}}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	return nil
}

func (ch *CommandHandler) handleGameStarting(packet *c.Packet) error {
	var gameStartingMessage c.GameStartingMessage
	if err := c.ParsePacketData(packet, &gameStartingMessage); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	// TODO: create a game object and properly hydrate it.
	// ch.game = NewClientGame()
	// ch.game.Hydrate(gameStartingMessage)

	return nil
}

func (ch *CommandHandler) handleSelectDice(packet *c.Packet) error {
	var selectDiceMessage c.SelectDiceMessage
	if err := c.ParsePacketData(packet, &selectDiceMessage); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	// TODO: properly propagate the dice status to the game object

	ch.ioh.DisplayMessage(ch.game.Data.Players[ch.myUsername].FormatDice())
	input := ""
	if err := ch.ioh.ReadInput(&input); err != nil {
		return fmt.Errorf("error while choosing dice: %w", err)
	}

	input = strings.ReplaceAll(input, " ", "")
	if input == "*" {
		input = "1,2,3,4,5,6"
	}

	_, err := regexp.MatchString("^([1-6],?){0,6}$", input)
	if err != nil {
		return fmt.Errorf("error while choosing dice: %w", err)
	}

	var keep [6]bool
	for i := 0; i < 6; i++ {
		keep[i] = strings.Contains(input, fmt.Sprintf("%d", i+1))
	}

	if err := c.SendPacket(ch.conn, c.KeepDice, &c.KeepDiceMessage{Kept: keep}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	return nil
}

func (ch *CommandHandler) handleErrorCommand(command c.Command) error {
	slog.Debug("Command did not work", "command", command)
	return nil
}

func (ch *CommandHandler) handleDefaultCase(command c.Command) error {
	slog.Warn("Server sent an unknown command", "command", command)
	return nil
}
