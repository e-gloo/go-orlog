package commons

import "fmt"

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func PrintDices(dices [6]Die) {
	for dice_nb, die := range dices {
		fmt.Print(1+dice_nb, die.Face().String())
	}
	fmt.Print("\n")
}
