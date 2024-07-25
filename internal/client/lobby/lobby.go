package lobby

type Phase int

const (
	CreateOrJoin Phase = iota + 1
	AddPlayerName
	Game
)

type Lobby struct {
	ServerUrl  string
	GameUuid   string
	PlayerName string
	Phase      Phase
	Err        string
}
