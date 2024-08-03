package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/e-gloo/orlog/internal/bbtea"
	"github.com/e-gloo/orlog/internal/pkg/logging"
)

func main() {
	// Client run config
	dev := flag.Bool("dev", false, "Running in development mode")
	flag.Parse()

	if *dev {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()

		logging.InitLogger(*dev, f, f)
	}

	client := bbtea.NewClient()
	if _, err := client.Run(); err != nil {
		slog.Error(err.Error())
	}

}
