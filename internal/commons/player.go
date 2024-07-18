package commons

type InitGamePlayer struct {
	Username   string
	Health     int
	GodIndexes [3]int
}

type UpdateGamePlayer struct {
	Health              int
	Tokens              int
	AxeDamageReceived   int
	ArrowDamageReceived int
	TokensGained        int
	TokensStolen        int
}

type PlayerMap[T any] map[string]T
