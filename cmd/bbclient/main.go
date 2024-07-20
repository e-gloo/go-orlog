package main

import (
	"flag"
	"log/slog"

	"github.com/e-gloo/orlog/internal/bbtea"
	"github.com/e-gloo/orlog/internal/pkg/logging"
)

func main() {
	// Client run config
	dev := flag.Bool("dev", false, "Running in development mode")
	flag.Parse()

	logging.InitLogger(*dev)

	client := bbtea.NewClient()
	if _, err := client.Run(); err != nil {
		slog.Error(err.Error())
	}

}
