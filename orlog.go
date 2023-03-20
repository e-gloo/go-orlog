package main

import "fmt"
import "time"
import "math/rand"

func printDices(dices [6]Dice) {
	for dice_nb, dice := range dices {
		fmt.Print(1+dice_nb, dice.Face().String())
	}
	fmt.Print("\n")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := InitGame()
	game.Play()
}
