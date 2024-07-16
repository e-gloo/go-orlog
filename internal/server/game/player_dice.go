package server_game

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

func (die *PlayerDie) Roll() {
	die.faceIndex = rand.Intn(6)
}

func (die *PlayerDie) Keep() {
	die.kept = true
}

func (die *PlayerDie) Unkeep() {
	die.kept = false
}

func (die *PlayerDie) Reset() {
	die.kept = false
	// die.faceIndex = 0
	die.quantity = 1
}

func (die *PlayerDie) IsKept() bool {
	return die.kept
}

func (die *PlayerDie) GetFaceIndex() int {
	return die.faceIndex
}

func (die *PlayerDie) GetQuantity() int {
	return die.quantity
}
