package client_game

type PlayerDice [6]*PlayerDie

type PlayerDie struct {
	kept     bool
	face     int
	quantity int
}

func NewPlayerDie() *PlayerDie {
	return &PlayerDie{
		kept:     false,
		face:     0,
		quantity: 1,
	}
}

func (d *PlayerDie) IsKept() bool {
	return d.kept
}

func (d *PlayerDie) GetFaceId() int {
	return d.face
}

func (d *PlayerDie) SetKept(state bool) {
	d.kept = state
}

func (d *PlayerDie) SetFaceId(face int) {
	d.face = face
}

// func (d *PlayerDie) GetQuantity() int {
// 	return d.quantity
// }
