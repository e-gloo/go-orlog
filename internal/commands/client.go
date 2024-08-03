package commands

// Client to Server
const (
	CreateGame Command = "create"
	JoinGame   Command = "join"
	AddPlayer  Command = "add_player"
	RollDice   Command = "roll_dice"
	PlayGod    Command = "play_god"
	KeepDice   Command = "keep_dice"
)

type CreateGameMessage struct {
}

type JoinGameMessage struct {
	Uuid string
}

type RollDiceMessage struct{}

type AddPlayerMessage struct {
	Username   string
	GodIndexes [3]int
}

type PlayGodMessage struct {
	GodIndex int // -1 if no god, or any index in game.gods
	GodLevel int // -1 if no god, or 0-2
}

type KeepDiceMessage struct {
	Kept [6]bool
}
