package main

import "fmt"
import "time"
import "math/rand"

func printDices(dices [6]Die) {
	for dice_nb, die := range dices {
		fmt.Print(1+dice_nb, die.Face().String())
	}
	fmt.Print("\n")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := InitGame()
	game.Play()
}
