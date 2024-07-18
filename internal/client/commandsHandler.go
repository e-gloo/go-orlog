package client

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	g "github.com/e-gloo/orlog/internal/client/game"
	c "github.com/e-gloo/orlog/internal/commands"
	"github.com/gorilla/websocket"
)

type CommandHandler struct {
	ioh  IOHandler
	conn *websocket.Conn
	game *g.ClientGame
}

func NewCommandHandler(ioh IOHandler, conn *websocket.Conn) *CommandHandler {
	return &CommandHandler{
		ioh:  ioh,
		conn: conn,
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
	case c.TurnFinished:
		return ch.handleTurnFinished(packet)
	case c.GameStarting:
		return ch.handleGameStarting(packet)
	case c.GameFinished:
		return ch.handleGameFinished(packet)
	case c.DiceRoll:
		return ch.handleDiceRoll(packet)
	case c.SelectDice:
		return ch.handleSelectDice(packet)
	case c.CommandError:
		return ch.handleErrorCommand(packet)
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

	input := "Player"
	err := ch.ioh.ReadInput(&input)
	if err != nil {
		return err
	}

	if err = c.SendPacket(ch.conn, c.AddPlayer, &c.AddPlayerMessage{Username: input, GodIndexes: [3]int{0, 0, 0}}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	return nil
}

func (ch *CommandHandler) handleTurnFinished(packet *c.Packet) error {
	var turnFinishedMessage c.TurnFinishedMessage
	if err := c.ParsePacketData(packet, &turnFinishedMessage); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	ch.game.UpdatePlayers(
		turnFinishedMessage.Players,
	)

	return nil
}

func (ch *CommandHandler) handleGameStarting(packet *c.Packet) error {
	var gameStartingMessage c.GameStartingMessage
	if err := c.ParsePacketData(packet, &gameStartingMessage); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	ch.game = g.NewClientGame(
		gameStartingMessage.YourUsername,
		gameStartingMessage.Dice,
		gameStartingMessage.Players,
	)

	ch.ioh.DisplayMessage(ch.game.MyUsername + ": GET READY FOR VALHALLA !")

	return nil
}

func (ch *CommandHandler) handleGameFinished(packet *c.Packet) error {
	var gameFinishedMessage c.GameFinishedMessage
	if err := c.ParsePacketData(packet, &gameFinishedMessage); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	ch.ioh.DisplayMessage("Game finished, the winner is: " + gameFinishedMessage.Winner)

	return nil
}

func (ch *CommandHandler) handleDiceRoll(packet *c.Packet) error {
	var diceRollMessage c.DiceRollMessage
	if err := c.ParsePacketData(packet, &diceRollMessage); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	ch.game.UpdatePlayersDice(
		diceRollMessage.Players,
	)

	ch.ioh.DisplayMessage(
		ch.game.FormatGame(),
	)

	return nil
}

func (ch *CommandHandler) handleSelectDice(packet *c.Packet) error {
	var selectDiceMessage c.SelectDiceMessage
	if err := c.ParsePacketData(packet, &selectDiceMessage); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	ch.ioh.DisplayMessage("Choose your dice to keep (1-6, separated by commas, * to keep all): ")
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
		return fmt.Errorf("error while validating chosen dice: %w", err)
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

func (ch *CommandHandler) handleErrorCommand(packet *c.Packet) error {
	var errorMessage c.CommandErrorMessage
	if err := c.ParsePacketData(packet, &errorMessage); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	slog.Info("Command did not work", "reason", errorMessage.Reason)
	return nil
}

func (ch *CommandHandler) handleDefaultCase(command c.Command) error {
	slog.Warn("Server sent an unknown command", "command", command)
	return nil
}
