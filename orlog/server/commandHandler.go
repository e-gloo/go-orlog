package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/e-gloo/orlog/orlog/commons"
	"github.com/gorilla/websocket"
)

var games = sync.Map{}

type gameManager struct {
	game        *Game
	player1Conn *websocket.Conn
	player2Conn *websocket.Conn
}

func commandHandler(conn *websocket.Conn, packet *commons.Packet) {
	switch packet.Command {
	case commons.Create:
		log.Printf("Creating new game")
		createData := &commons.CreateData{}
		json.Unmarshal(packet.Data, createData)

		game, err := InitGame(createData.Player)
		var response []byte
		if err != nil {
			log.Printf("Error creating game: %v", err)
			response = []byte("Failed")
		} else {
            manager := &gameManager{
                game: game,
                player1Conn: conn,
            }
            games.Store(game.uuid, manager)
			response = []byte("Succeed")
		}
		newPacket := &commons.Packet{
			Command: commons.Create,
			Data:    response,
		}
		newPacketBuffer := new(bytes.Buffer)
		json.NewEncoder(newPacketBuffer).Encode(newPacket)
		log.Printf("Game created: %s", game.uuid)

		conn.WriteMessage(websocket.TextMessage, newPacketBuffer.Bytes())
	case commons.Join:
		log.Printf("Joining game")
		createData := &commons.CreateData{}
		json.Unmarshal(packet.Data, createData)
		value, ok := games.Load(createData.Uuid)
		manager := value.(*gameManager)
		if !ok {
			log.Printf("Error joining game")
		} else {
			log.Printf("Joining success")
		}
		manager.game.AddPlayer(createData.Player)
        manager.player2Conn = conn
		log.Printf("Player added to the game")
		log.Printf("P1 %s, P2 %s", manager.game.player1.Name, manager.game.player2.Name)
        manager.player1Conn.WriteMessage(websocket.TextMessage, []byte("Game Starting"))
        manager.player2Conn.WriteMessage(websocket.TextMessage, []byte("Game Starting"))
	}
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		// log.Printf("recv: %s", message)
		packet := &commons.Packet{}
		json.Unmarshal(message, packet)
		commandHandler(conn, packet)

		fmt.Println("Got Message")

		//		err = c.WriteMessage(mt, message)
		//		if err != nil {
		//			log.Println("write:", err)
		//			break
		//		}
	}
}
