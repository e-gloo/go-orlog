package commands

// Server to Client
const (
	CreateOrJoin    Command = "create_or_join"
	CreatedOrJoined Command = "created_or_joined"
	ConfigurePlayer Command = "configure_player"
	SelectDice      Command = "select_dice"
	WantToPlayGods  Command = "want_to_play_gods"
	GameStarting    Command = "starting"
	CommandError    Command = "error"
)

type CreateOrJoinMessage struct {
	Welcome string
	// Lobbies []string
}

type CreatedOrJoinedMessage struct {
	Uuid string
}

type ConfigurePlayerMessageGod struct {
	Name        string
	Description string
	Priority    int
}

type ConfigurePlayerMessage struct {
	Gods []ConfigurePlayerMessageGod
}

type SelectDiceMessagePlayer struct {
	FaceIndexes [6]int
	FacesKept   [6]bool
}

type SelectDiceMessage struct {
	Turn int
	P1   SelectDiceMessagePlayer
	P2   SelectDiceMessagePlayer
}

type WantToPlaysGodsMessage struct {
	// maybe the list of selected gods,
	// or nothing if we trust client
}

type GameStartingMessageDieFace struct {
	Face string
}

type GameStartingMessageDie struct {
	Faces [6]GameStartingMessageDieFace
}

type GameStartingMessagePlayer struct {
	Username   string
	Health     int
	GodIndexes [3]int
}

type GameStartingMessage struct {
	Dice [6]GameStartingMessageDie
	P1   GameStartingMessagePlayer
	P2   GameStartingMessagePlayer
}

type CommandErrorMessage struct {
	Reason string
}
