package main

import "fmt"

func main() {
	dices := Init()
	for dice_nb, dice := range dices {
		fmt.Println("Dice number", 1 + dice_nb)
		for _, face := range dice.faces {
			fmt.Println(face.kind)
		}
	}
}
