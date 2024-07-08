package server

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/e-gloo/orlog/internal/commands"
	"github.com/gorilla/websocket"
)

var joinableGames = sync.Map{}

type gameManager struct {
	game    *ServerGame
	players map[string]*CommandHandler
}

type CommandHandler struct {
	Conn             *websocket.Conn
	Username         string
	manager          *gameManager
	isHost           bool
	ExpectedCommands []commands.Command
}

func NewCommandHandler(conn *websocket.Conn) *CommandHandler {
	return &CommandHandler{
		Conn:     conn,
		isHost:   false,
		Username: "Player",
		ExpectedCommands: []commands.Command{
			commands.CreateGame,
			commands.JoinGame,
		},
	}
}

func (ch *CommandHandler) Handle(packet *commands.Packet) error {
	if !slices.Contains(ch.ExpectedCommands, packet.Command) {
		slog.Warn("Unexpected command", "command", packet.Command)
		if err := commands.SendPacket(ch.Conn, &commands.Packet{Command: commands.CommandError}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
		return nil
	}

	switch packet.Command {
	case commands.CreateGame:
		return ch.handleCreateGame()
	case commands.JoinGame:
		return ch.handleJoinGame(packet)
	case commands.AddPlayer:
		return ch.handleAddPlayer(packet)
	case commands.KeepDices:
		return ch.handleKeepDices(packet)
	case commands.CommandError:
		slog.Debug("Oops désolé :D")
		return nil
	default:
		return ch.handleDefaultCase(packet.Command)
	}
}

func (ch *CommandHandler) handleCreateGame() error {
	slog.Info("Creating new game")

	game, err := NewServerGame()
	if err != nil {
		return fmt.Errorf("error initializing game: %w", err)
	}

	manager := &gameManager{
		game:    game,
		players: make(map[string]*CommandHandler),
	}
	joinableGames.Store(game.Uuid, manager)
	ch.isHost = true
	ch.manager = manager
	ch.ExpectedCommands = []commands.Command{commands.AddPlayer}
	slog.Info("Game created", "uuid", game.Uuid)

	if err := commands.SendPacket(ch.Conn, &commands.Packet{Command: commands.CommandOK, Data: game.Uuid}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	if err := commands.SendPacket(ch.Conn, &commands.Packet{Command: commands.AddPlayer}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	return nil
}

func (ch *CommandHandler) handleJoinGame(packet *commands.Packet) error {
	slog.Info("Trying to join...", "uuid", packet.Data)
	value, ok := joinableGames.Load(packet.Data)
	if !ok {
		slog.Debug("Error joining game, uuid not found", "uuid", packet.Data)
		return nil
	}

	manager := value.(*gameManager)
	ch.manager = manager
	ch.ExpectedCommands = []commands.Command{commands.AddPlayer}
	slog.Info("Joined game", "uuid", packet.Data)
	joinableGames.Delete(packet.Data)

	if err := commands.SendPacket(ch.Conn, &commands.Packet{Command: commands.CommandOK, Data: packet.Data}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	if err := commands.SendPacket(ch.Conn, &commands.Packet{Command: commands.AddPlayer}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}

func (ch *CommandHandler) handleAddPlayer(packet *commands.Packet) error {

	if err := ch.manager.game.Data.AddPlayer(packet.Data); err != nil {
		if err := commands.SendPacket(ch.Conn, &commands.Packet{Command: commands.CommandError, Data: err.Error()}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}

		if err := commands.SendPacket(ch.Conn, &commands.Packet{Command: commands.AddPlayer}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
		return nil
	}

	ch.manager.players[packet.Data] = ch
	ch.Username = packet.Data

	if ch.isHost {
		slog.Info("Player 1 added", "name", packet.Data)
	} else {
		slog.Info("Player 2 added", "name", packet.Data)
	}

	ch.ExpectedCommands = []commands.Command{}

	if err := commands.SendPacket(ch.Conn, &commands.Packet{Command: commands.CommandOK}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	if ch.manager.game.Data.IsGameReady() {
		if err := ch.handleStartingGame(); err != nil {
			return fmt.Errorf("error starting the game: %w", err)
		}
	}
	return nil
}

func (ch *CommandHandler) handleStartingGame() error {
	ch.manager.game.Data.SelectFirstPlayer()
	slog.Info("Game is starting...")

	gameData, err := ch.manager.game.String()
	if err != nil {
		return fmt.Errorf("error serializing game data: %w", err)
	}

	// Send every player the data to init the game
	for u := range ch.manager.players {
		if err := commands.SendPacket(ch.manager.players[u].Conn, &commands.Packet{Command: commands.GameStarting, Data: gameData}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}
	}

	firstUsername := ch.manager.game.Data.PlayersOrder[0]
	secondUsername := ch.manager.game.Data.PlayersOrder[1]
	ch.manager.game.Data.Players[firstUsername].RollDices()

	gameData, err = ch.manager.game.String()
	if err != nil {
		return fmt.Errorf("error serializing game data: %w", err)
	}

	// Send P1 the gameData after its first roll.
	if err := commands.SendPacket(ch.manager.players[firstUsername].Conn, &commands.Packet{Command: commands.SelectDices, Data: gameData}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	ch.manager.players[firstUsername].ExpectedCommands = []commands.Command{commands.KeepDices}
	ch.manager.players[secondUsername].ExpectedCommands = []commands.Command{}

	ch.manager.game.Rolls++

	return nil
}

func (ch *CommandHandler) handleKeepDices(packet *commands.Packet) error {
	// TODO: validate Packet.Data formating using regexp ?

	to_keep := strings.Split(packet.Data, ",")

	for _, dice_nb := range to_keep {
		i, err := strconv.ParseInt(dice_nb, 10, 32)
		if err != nil {
			continue
		}
		ch.manager.game.Data.Players[ch.Username].Dices[i-1].Kept = true
	}

	if ch.manager.game.Rolls >= 4 {
		for u := range ch.manager.players {
			ch.manager.game.Data.Players[u].RollDices()
		}

		ch.manager.game.Data.ComputeRound()
		ch.manager.game.Rolls = 0

		gameData, err := ch.manager.game.String()
		if err != nil {
			return fmt.Errorf("error serializing game data: %w", err)
		}

		// Send every player the update of the game
		for u := range ch.manager.players {
			if err := commands.SendPacket(ch.manager.players[u].Conn, &commands.Packet{Command: commands.GameInfo, Data: gameData}); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}
		}

		if ch.manager.game.Data.Players[ch.manager.game.Data.PlayersOrder[1]].Health <= 0 {
			// P1 won
			slog.Info("Bravo P1 :)")
			panic("Bravo P1 :)")
		} else if ch.manager.game.Data.Players[ch.manager.game.Data.PlayersOrder[0]].Health <= 0 {
			// P2 won
			slog.Info("Bravo P2 :)")
			panic("Bravo P2 :)")
		} else {
			ch.manager.game.Data.ChangePlayersPosition()

			firstUsername := ch.manager.game.Data.PlayersOrder[0]
			secondUsername := ch.manager.game.Data.PlayersOrder[1]
			ch.manager.game.Data.Players[firstUsername].RollDices()

			gameData, err := ch.manager.game.String()
			if err != nil {
				return fmt.Errorf("error serializing game data: %w", err)
			}

			// Send P1 the gameData after its first roll.
			if err := commands.SendPacket(ch.manager.players[firstUsername].Conn, &commands.Packet{Command: commands.SelectDices, Data: gameData}); err != nil {
				return fmt.Errorf("error sending packet: %w", err)
			}

			ch.manager.players[firstUsername].ExpectedCommands = []commands.Command{commands.KeepDices}
			ch.manager.players[secondUsername].ExpectedCommands = []commands.Command{}

			ch.manager.game.Rolls++
		}
	} else {
		otherUsername := ch.manager.game.Data.PlayersOrder[slices.IndexFunc(ch.manager.game.Data.PlayersOrder, func(p string) bool {
			return p != ch.Username
		})]

		ch.manager.game.Data.Players[otherUsername].RollDices()

		gameData, err := ch.manager.game.String()
		if err != nil {
			return fmt.Errorf("error serializing game data: %w", err)
		}

		// Send otherPlayer the gameData after its roll.
		if err := commands.SendPacket(ch.manager.players[otherUsername].Conn, &commands.Packet{Command: commands.SelectDices, Data: gameData}); err != nil {
			return fmt.Errorf("error sending packet: %w", err)
		}

		ch.manager.players[ch.Username].ExpectedCommands = []commands.Command{}
		ch.manager.players[otherUsername].ExpectedCommands = []commands.Command{commands.KeepDices}

		ch.manager.game.Rolls++
	}

	return nil
}

func (ch *CommandHandler) handleDefaultCase(command commands.Command) error {
	slog.Debug("Unknown command", "command", command)
	if err := commands.SendPacket(ch.Conn, &commands.Packet{Command: commands.CommandError}); err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}
