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
	for conn == nil {
		url, err := readServerUrl()
		if err != nil {
			slog.Error("Could not read server URL", "err", err)
			return
		}

		conn, err = client.NewClient(url)
		if err != nil {
			slog.Error("Error connecting to server", "err", err)
			continue
		}
	}

	if err := client.ListenForServer(conn); err != nil {
		slog.Error("Error starting game", "err", err)
	}
}

func readServerUrl() (*url.URL, error) {
	serverInput := "localhost:8080"
	fmt.Println("Enter the server URL:")
	_, err := fmt.Scanln(&serverInput)
	if err != nil && err.Error() != "unexpected newline" {
		return nil, err
	}

	u := &url.URL{Scheme: "ws", Host: serverInput, Path: "/connect"}

	return u, nil
}
