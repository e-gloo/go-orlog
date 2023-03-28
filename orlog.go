package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func printDices(dices [6]Die) {
	for dice_nb, die := range dices {
		fmt.Print(1+dice_nb, die.Face().String())
	}
	fmt.Print("\n")
}

func parseArgs() (isServer bool, port int, host string) {
	isServer = false
	port = 8080
	host = "localhost"

	for argIdx, _ := range os.Args {
		if os.Args[argIdx] == "-s" || os.Args[argIdx] == "--server" {
			isServer = true
		}
		if os.Args[argIdx] == "-p" || os.Args[argIdx] == "--port" {
			i, err := strconv.ParseInt(os.Args[argIdx+1], 10, 64)
			if err != nil {
				continue
			}
			port = int(i)
		}
		if os.Args[argIdx] == "-h" || os.Args[argIdx] == "--host" {
			host = os.Args[argIdx+1]
		}
		if os.Args[argIdx] == "--help" {
			fmt.Println("Usage: orlog [OPTIONS]")
			fmt.Println("Options:")
			fmt.Println("  -s, --server\t\tStart a server")
			fmt.Println("  -p, --port\t\tSpecify the port to use")
			fmt.Println("  -h, --host\t\tSpecify the host to use")
			fmt.Println("  --help\t\tPrint this help")
			os.Exit(0)
		}
	}

	return isServer, port, host
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// isServer, port, host := parseArgs()
	// fmt.Println(isServer, port, host)

	game := InitGame()
	game.Play()
}
