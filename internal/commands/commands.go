package commands

type Command string

const (
	CreateGame   Command = "create"
	JoinGame     Command = "join"
	AddPlayer    Command = "add_player"
	ChooseGods   Command = "choose_gods"
	GameStarting Command = "starting"
	PlayGods     Command = "play_gods"
	KeepDices    Command = "keep_dices"
)

const (
	SelectDices    Command = "select_dices"
	WantToPlayGods Command = "want_to_play_gods"
)

const (
	CommandOK    Command = "ok"
	CommandError Command = "error"
)
