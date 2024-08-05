package server

import (
	"fmt"
	"log/slog"
	"math"
	"sync"

	c "github.com/e-gloo/orlog/internal/commands"
	g "github.com/e-gloo/orlog/internal/server/game"
	"github.com/gorilla/websocket"
)

var joinableGames = sync.Map{}

type CommandHandler struct {
	Conn     *websocket.Conn
	Username string
	game     *g.ServerGame
	isHost   bool
}

func NewCommandHandler(conn *websocket.Conn) *CommandHandler {
	return &CommandHandler{
		Conn:     conn,
		isHost:   false,
		Username: "Player",
	}
}

func (ch *CommandHandler) Handle(packet *c.Packet) error {
	// if !slices.Contains(ch.game.Players[ch.Username].GetExpectedCommands(), packet.Command) {
	// 	return ch.handleUnexpectedCommand(packet.Command)
	// }

	switch packet.Command {
	case c.CreateGame:
		return ch.handleCreateGame()
	case c.JoinGame:
		return ch.handleJoinGame(packet)
	case c.AddPlayer:
		return ch.handleAddPlayer(packet)
	case c.RollDice:
		return ch.handleRollDice()
	case c.KeepDice:
		return ch.handleKeepDice(packet)
	case c.PlayGod:
		return ch.handlePlayGod(packet)
	default:
		return ch.handleDefaultCase(packet.Command)
	}
}

func (ch *CommandHandler) handleCreateGame() error {
	slog.Info("Creating new game")

	game, err := g.NewServerGame()
	if err != nil {
		return fmt.Errorf("error initializing game: %w", err)
	}

	joinableGames.Store(game.Uuid, game)
	ch.isHost = true
	ch.game = game
	// ch.ExpectedCommands = []c.Command{c.AddPlayer}
	slog.Info("Game created", "uuid", game.Uuid)

	var createdOrJoinedMessage c.CreatedOrJoinedMessage
	createdOrJoinedMessage.Uuid = game.Uuid
	if err := c.SendPacket(ch.Conn, c.CreatedOrJoined, &createdOrJoinedMessage); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	var configurePlayerMessage c.ConfigurePlayerMessage
	configurePlayerMessage.Gods = ch.game.GetGodsDefinition()
	if err := c.SendPacket(ch.Conn, c.ConfigurePlayer, &configurePlayerMessage); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	return nil
}

func (ch *CommandHandler) handleJoinGame(packet *c.Packet) error {
	var joinGameMessage c.JoinGameMessage
	if err := c.ParsePacketData(packet, &joinGameMessage); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	slog.Info("Trying to join...", "uuid", joinGameMessage.Uuid)
	var commandErrorMessage c.CommandErrorMessage

	value, ok := joinableGames.Load(joinGameMessage.Uuid)
	if !ok {
		slog.Debug("Error joining game, uuid not found", "uuid", joinGameMessage.Uuid)
		commandErrorMessage.Reason = "Game not found, try again."
		if err := c.SendPacket(ch.Conn, c.CommandError, &commandErrorMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
		return nil
	}

	game, ok := value.(*g.ServerGame)
	if !ok {
		err := fmt.Errorf("could not retrieve a valid game")
		slog.Debug(err.Error(), "uuid", joinGameMessage.Uuid)
		commandErrorMessage.Reason = err.Error()
		if err := c.SendPacket(ch.Conn, c.CreateOrJoin, &commandErrorMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
		return err
	}
	ch.game = game

	// ch.ExpectedCommands = []c.Command{c.AddPlayer}
	slog.Info("Joined game", "uuid", joinGameMessage.Uuid)
	joinableGames.Delete(joinGameMessage.Uuid)

	var createdOrJoinedMessage c.CreatedOrJoinedMessage
	createdOrJoinedMessage.Uuid = joinGameMessage.Uuid
	if err := c.SendPacket(ch.Conn, c.CreatedOrJoined, &createdOrJoinedMessage); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	var configurePlayerMessage c.ConfigurePlayerMessage
	configurePlayerMessage.Gods = ch.game.GetGodsDefinition()
	if err := c.SendPacket(ch.Conn, c.ConfigurePlayer, &configurePlayerMessage); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}

func (ch *CommandHandler) handleAddPlayer(packet *c.Packet) error {
	var message c.AddPlayerMessage
	if err := c.ParsePacketData(packet, &message); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	slog.Debug("addPlayer", "game", ch.game, "username", message.Username)

	if err := ch.game.AddPlayer(ch.Conn, message.Username, message.GodIndexes); err != nil {
		var errorMessage c.CommandErrorMessage
		errorMessage.Reason = err.Error()
		if err := c.SendPacket(ch.Conn, c.CommandError, &errorMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}

		var configurePlayerMessage c.ConfigurePlayerMessage
		configurePlayerMessage.Gods = ch.game.GetGodsDefinition()
		if err := c.SendPacket(ch.Conn, c.ConfigurePlayer, &configurePlayerMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
		return nil
	}

	ch.Username = message.Username

	if ch.isHost {
		slog.Info("Player 1 added", "name", message.Username)
	} else {
		slog.Info("Player 2 added", "name", message.Username)
	}

	// ch.ExpectedCommands = []c.Command{}

	if ch.game.IsGameReady() {
		if err := ch.handleStartingGame(); err != nil {
			return fmt.Errorf("error starting the game: %w", err)
		}
	}
	return nil
}

func (ch *CommandHandler) handleStartingGame() error {
	ch.game.Restart()
	slog.Info("Game is starting...")

	var gameStartingMessage c.GameStartingMessage
	gameStartingMessage.Players, gameStartingMessage.Dice = ch.game.GetStartingDefinition()
	gameStartingMessage.Gods = ch.game.GetGodsDefinition()
	for u := range ch.game.Players {
		gameStartingMessage.YourUsername = u
		if err := c.SendPacket(ch.game.Players[u].Conn, c.GameStarting, &gameStartingMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
	}

	// ask first player to roll
	firstUsername := ch.game.PlayersOrder[0]

	// ch.players[firstUsername].ExpectedCommands = []c.Command{c.RollDice}

	var askRollDiceMessage c.AskRollDiceMessage
	askRollDiceMessage.Player = firstUsername
	for u := range ch.game.Players {
		if err := c.SendPacket(ch.game.Players[u].Conn, c.AskRollDice, &askRollDiceMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
	}

	return nil
}

func (ch *CommandHandler) handleRollDice() error {
	ch.game.Players[ch.Username].RollDice()

	var diceRollMessage c.DiceRollMessage
	diceRollMessage.DiceState = ch.game.GetPlayerRollDiceState(ch.Username)
	diceRollMessage.Player = ch.Username
	for u := range ch.game.Players {
		if err := c.SendPacket(ch.game.Players[u].Conn, c.DiceRoll, &diceRollMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
	}

	var selectDiceMessage c.SelectDiceMessage
	selectDiceMessage.Turn = int(math.Ceil(float64(ch.game.Rolls) / 2))
	selectDiceMessage.Player = ch.Username
	for u := range ch.game.Players {
		if err := c.SendPacket(ch.game.Players[u].Conn, c.SelectDice, &selectDiceMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
	}

	// ch.players[firstUsername].ExpectedCommands = []c.Command{c.KeepDice}
	// ch.players[secondUsername].ExpectedCommands = []c.Command{}

	ch.game.Rolls++

	return nil
}

func (ch *CommandHandler) handleKeepDice(packet *c.Packet) error {
	var message c.KeepDiceMessage
	if err := c.ParsePacketData(packet, &message); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	for dice_idx, dice_kept := range message.Kept {
		if dice_kept {
			ch.game.Players[ch.Username].GetDice()[dice_idx].Keep()
		} else {
			ch.game.Players[ch.Username].GetDice()[dice_idx].Unkeep()
		}
	}
	//
	var diceStateMessage c.DiceStateMessage
	diceStateMessage.DiceState = ch.game.GetRollDiceState()
	for u := range ch.game.Players {
		if err := c.SendPacket(ch.game.Players[u].Conn, c.DiceState, &diceStateMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
	}

	if ch.game.Rolls >= 4 {
		for u := range ch.game.Players {
			ch.game.Players[u].RollDice()

		}

		var diceStateMessage c.DiceStateMessage
		diceStateMessage.DiceState = ch.game.GetRollDiceState()
		for u := range ch.game.Players {
			if err := c.SendPacket(ch.game.Players[u].Conn, c.DiceState, &diceStateMessage); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}
		}

		// var askToPlayGodMessage c.AskToPlayGodMessage
		// if err := c.SendPacket(ch.game.Players[ch.game.PlayersOrder[0]].Conn, c.AskToPlayGod, &askToPlayGodMessage); err != nil {
		// 	return fmt.Errorf("error sending packet: %w", err)
		// }
	} else {
		otherUsername := ch.game.GetOpponentName(ch.Username)

		var askRollDiceMessage c.AskRollDiceMessage
		askRollDiceMessage.Player = otherUsername
		for u := range ch.game.Players {
			if err := c.SendPacket(ch.game.Players[u].Conn, c.AskRollDice, &askRollDiceMessage); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}
		}

		// ch.game.Players[otherUsername].RollDice()
		//
		// var rollDiceMessage c.DiceRollMessage
		// rollDiceMessage.Players = ch.game.GetRollDiceState()
		//
		// for u := range ch.game.Players {
		// 	if err := c.SendPacket(ch.game.Players[u].Conn, c.DiceRoll, &rollDiceMessage); err != nil {
		// 		return fmt.Errorf("error sending packet: %w", err)
		// 	}
		// }
		//
		// var selectDiceMessage c.SelectDiceMessage
		// selectDiceMessage.Turn = int(math.Ceil(float64(ch.game.Rolls) / 2))
		// if err := c.SendPacket(ch.game.Players[otherUsername].Conn, c.SelectDice, &selectDiceMessage); err != nil {
		// 	return fmt.Errorf("error sending packet: %w", err)
		// }

		// ch.players[ch.Username].ExpectedCommands = []c.Command{}
		// ch.players[otherUsername].ExpectedCommands = []c.Command{c.KeepDice}

		// ch.game.Rolls++
	}
	//
	return nil
}

func (ch *CommandHandler) handlePlayGod(packet *c.Packet) error {
	var message c.PlayGodMessage
	if err := c.ParsePacketData(packet, &message); err != nil {
		return fmt.Errorf("error parsing packet data: %w", err)
	}

	ch.game.Players[ch.Username].SelectGod(message.GodIndex, message.GodLevel)

	if ch.Username == ch.game.PlayersOrder[0] {
		// ask P2 to play god
		var askToPlayGodMessage c.AskToPlayGodMessage
		if err := c.SendPacket(ch.game.Players[ch.game.PlayersOrder[1]].Conn, c.AskToPlayGod, &askToPlayGodMessage); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
	} else {
		// round is over
		slog.Info("Round is over, computing ...")

		var turnFinishedMessage c.TurnFinishedMessage
		turnFinishedMessage.Turn = ch.game.GetTurn()

		turnFinishedMessage.Players = ch.game.ComputeRound()

		for u := range ch.game.Players {
			if err := c.SendPacket(ch.game.Players[u].Conn, c.TurnFinished, &turnFinishedMessage); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}
		}

		if ch.game.IsGameFinished() {
			ch.handleGameFinished()
		} else {
			// ch.game.ChangePlayersPosition()
			//
			// firstUsername := ch.game.PlayersOrder[0]
			// // secondUsername := ch.game.PlayersOrder[1]
			// ch.game.Players[firstUsername].RollDice()
			//
			// var diceRollMessage c.DiceRollMessage
			// diceRollMessage.Players = ch.game.GetRollDiceState()
			// for u := range ch.game.Players {
			// 	if err := c.SendPacket(ch.game.Players[u].Conn, c.DiceRoll, &diceRollMessage); err != nil {
			// 		return fmt.Errorf("error sending packet: %w", err)
			// 	}
			// }
			//
			// var selectDiceMessage c.SelectDiceMessage
			// selectDiceMessage.Turn = int(math.Ceil(float64(ch.game.Rolls) / 2))
			// if err := c.SendPacket(ch.game.Players[firstUsername].Conn, c.SelectDice, &selectDiceMessage); err != nil {
			// 	return fmt.Errorf("error sending packet: %w", err)
			// }
			//
			// // ch.players[firstUsername].ExpectedCommands = []c.Command{c.KeepDice}
			// // ch.players[secondUsername].ExpectedCommands = []c.Command{}
			//
			// ch.game.Rolls++
		}
	}

	return nil
}

func (ch *CommandHandler) handleGameFinished() error {
	if ch.game.Players[ch.game.PlayersOrder[1]].GetHealth() <= 0 {
		// P1 won
		slog.Info("Congratulations P1, you won ! :)")

		var gameFinishedMessage c.GameFinishedMessage
		gameFinishedMessage.Winner = ch.game.PlayersOrder[0]

		for u := range ch.game.Players {
			if err := c.SendPacket(ch.game.Players[u].Conn, c.GameFinished, &gameFinishedMessage); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}
		}

		ch.handleStartingGame()
	} else if ch.game.Players[ch.game.PlayersOrder[0]].GetHealth() <= 0 {
		// P2 won
		slog.Info("Congratulations P2, you won ! :)")

		var gameFinishedMessage c.GameFinishedMessage
		gameFinishedMessage.Winner = ch.game.PlayersOrder[1]

		for u := range ch.game.Players {
			if err := c.SendPacket(ch.game.Players[u].Conn, c.GameFinished, &gameFinishedMessage); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}
		}

		ch.handleStartingGame()
	} else {
		return fmt.Errorf("game is not finished")
	}

	return nil
}

func (ch *CommandHandler) HandleRagequit() error {
	if ch.game == nil {
		return fmt.Errorf("no game to ragequit")
	}

	opponent := ch.game.GetOpponentName(ch.Username)

	var gameFinishedMessage c.GameFinishedMessage
	gameFinishedMessage.Winner = opponent

	if err := c.SendPacket(ch.game.Players[opponent].Conn, c.GameFinished, &gameFinishedMessage); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	return nil
}

func (ch *CommandHandler) handleUnexpectedCommand(command c.Command) error {
	slog.Warn("Unexpected command", "command", command)
	var errorMessage c.CommandErrorMessage
	errorMessage.Reason = "unexpected command"
	if err := c.SendPacket(ch.Conn, c.CommandError, &errorMessage); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}

func (ch *CommandHandler) handleDefaultCase(command c.Command) error {
	slog.Debug("Unknown command", "command", command)
	var errorMessage c.CommandErrorMessage
	errorMessage.Reason = "unknown command"
	if err := c.SendPacket(ch.Conn, c.CommandError, &errorMessage); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}
