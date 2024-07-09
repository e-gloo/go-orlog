package commands

// Client to Server
const (
	CreateGame Command = "create"
	JoinGame   Command = "join"
	AddPlayer  Command = "add_player"
	PlayGods   Command = "play_gods"
	KeepDice   Command = "keep_dice"
)

type CreateGameMessage struct {
}

type JoinGameMessage struct {
	Uuid string
}

type AddPlayerMessage struct {
	Username   string
	GodIndexes [3]int
}

type PlayGodsMessage struct {
	GodIndex int // -1 if no god, or 0,1,2
}

type KeepDiceMessage struct {
	Kept [6]bool
}
