package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/e-gloo/orlog/internal/client"
	"github.com/e-gloo/orlog/internal/pkg/logging"
	"github.com/gorilla/websocket"
)

func main() {
	// Client run config
	dev := flag.Bool("dev", false, "Running in development mode")
	flag.Parse()

	logging.InitLogger(*dev)

	var conn *websocket.Conn

	var ioh client.IOHandler = &client.TermHandler{}

	for conn == nil {
		url, err := readServerUrl(ioh)
		if err != nil {
			slog.Error("Could not read server URL", "err", err)
			return
		}

		conn, err = client.NewClient(url)
		if err != nil {
			slog.Error("Error connecting to server", "err", err)
			ioh.DisplayMessage(fmt.Sprintf("Failed to connect to %s", url))
			continue
		}
	}

	if err := client.ListenForServer(conn, ioh); err != nil {
		slog.Error("Error starting game", "err", err)
	}
}

func readServerUrl(ioh client.IOHandler) (*url.URL, error) {
	serverInput := "localhost:8080"
	ioh.DisplayMessage("Enter the server URL:")
	err := ioh.ReadInput(&serverInput)
	if err != nil {
		return nil, err
	}

	u := &url.URL{Scheme: "ws", Host: serverInput, Path: "/connect"}

	return u, nil
}
