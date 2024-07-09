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

type SelectDiceMessagePlayerDie struct {
	FaceIndex int
	FaceKept  bool
}

type SelectDiceMessagePlayer []SelectDiceMessagePlayerDie

type SelectDiceMessage struct {
	Turn    int
	Players map[string]SelectDiceMessagePlayer
}

type WantToPlaysGodsMessage struct {
	// maybe the list of selected gods,
	// or nothing if we trust client
}

type GameStartingMessageDieFace struct {
	Face  string
	Magic bool
}

type GameStartingMessageDie struct {
	Faces []GameStartingMessageDieFace
}

type GameStartingMessagePlayer struct {
	Username   string
	Health     int
	GodIndexes [3]int
}

type GameStartingMessage struct {
	YourUsername string
	Dice         []GameStartingMessageDie
	Players      map[string]GameStartingMessagePlayer
}

type CommandErrorMessage struct {
	Reason string
}
