package client_game

type ClientGodPower struct {
	Description string
	TokenCost   int
}

type ClientGod struct {
	Emoji       string
	Name        string
	Description string
	Priority    int
	Levels      [3]ClientGodPower
}
