package commons

type InitGameDieFace struct {
	Kind  string
	Magic bool
}

type InitGameDie struct {
	Faces []InitGameDieFace
}

type DieState struct {
	Index int
	Kept  bool
}

type DiceState []DieState
