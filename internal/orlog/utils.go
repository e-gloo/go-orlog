package orlog

import "fmt"

func PrintDices(dices [6]Die) {
	for dice_nb, die := range dices {
		fmt.Print(1+dice_nb, die.Face().String())
	}
	fmt.Print("\n")
}
