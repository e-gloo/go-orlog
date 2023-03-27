package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
    "encoding/json"

	"github.com/gorilla/websocket"
    "github.com/e-gloo/orlog/orlog/client"
    "github.com/e-gloo/orlog/orlog/commons"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var join = flag.String("join", "", "uuid of the game to join")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/connect"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

    client.StartGame(conn, *join)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
            packet := &commons.Packet{}
            json.Unmarshal(message, packet)
			log.Printf("packet: %s", packet)

		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
