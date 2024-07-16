package commons

type InitGamePlayer struct {
	Username   string
	Health     int
	GodIndexes [3]int
}

type PlayerMap[T any] map[string]T
