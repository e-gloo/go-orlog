package main

import "fmt"
import "time"
import "strings"
import "strconv"
import "math/rand"

func printDices(dices [6]Dice) {
	for dice_nb, dice := range dices {
		fmt.Print(1+dice_nb, dice.Face().String())
	}
	fmt.Print("\n")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	players := InitPlayers()

	for i := 0; i < 3; i++ {
		for player_idx, _ := range players {
			fmt.Println("Turn", i+1, players[player_idx].name)
			players[player_idx].RollDices()
			printDices(players[player_idx].dices)

			// We dont reroll on the last turn
			if i < 2 {
				input := ""
				fmt.Scanln(&input)

				to_keep := strings.Split(input, ",")

				for _, dice_nb := range to_keep {
					i, err := strconv.ParseInt(dice_nb, 10, 32)
					if err != nil {
						continue
					}
					players[player_idx].dices[i-1].kept = true
				}
			}
		}
	}

}
