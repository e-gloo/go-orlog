package main

import "fmt"
import "time"
import "strings"
import "strconv"
import "math/rand"

func print_dices(dices [6]Dice) {
	for dice_nb, dice := range dices {
		fmt.Println("Dice number", 1 + dice_nb, dice.Face())
	}
}

func roll_dices(dices *[6]Dice) {06:49 PM
	for idx, _ := range dices {
		if dices[idx].kept == false {
			dices[idx].Roll()
		}
		dices[idx].kept = false
	}
	print_dices(*dices)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	dices := Init()

	roll_dices(&dices)
	for i := 0; i < 3; i++ {
		input := ""
		fmt.Scanln(&input)

		to_keep := strings.Split(input, ",")

		for _, dice_nb := range to_keep {
			i, err := strconv.ParseInt(dice_nb, 10, 32)
			if err != nil {
				continue
			}
			dices[i - 1].kept = true
		}

		roll_dices(&dices)
	}

}
