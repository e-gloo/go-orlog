package server

import "math/rand"

type PlayerDice [6]*PlayerDie

type PlayerDie struct {
	kept      bool
	faceIndex int
	quantity  int
}

func NewPlayerDie() *PlayerDie {
	return &PlayerDie{
		kept:      false,
		faceIndex: 0,
		quantity:  1,
	}
}

func (pd *PlayerDie) Roll() {
	pd.faceIndex = rand.Intn(6)
}

func (pd *PlayerDie) Keep() {
	pd.kept = true
}

func (pd *PlayerDie) Unkeep() {
	pd.kept = false
}

func (pd *PlayerDie) IsKept() bool {
	return pd.kept
}

func (pd *PlayerDie) GetFaceIndex() int {
	return pd.faceIndex
}

func (pd *PlayerDie) GetQuantity() int {
	return pd.quantity
}
