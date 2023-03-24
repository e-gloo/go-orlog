package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/e-gloo/orlog/orlog/commons"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var games = make(map[string]*Game)

func commandHandler(conn *websocket.Conn, packet *commons.Packet) {
	switch packet.Command {
	case commons.Create:
		log.Printf("Creating new game")
		player := &commons.Player{}
		json.Unmarshal(packet.Data, player)

		game, err := InitGame(player)
		var response []byte
		if err != nil {
			log.Printf("Error creating game: %v", err)
			response = []byte("Failed")
		} else {
			games[game.uuid] = game
			response = []byte("Succeed")
		}
		newPacket := &commons.Packet{
			Command: commons.Create,
			Data:    response,
		}
		newPacketBuffer := new(bytes.Buffer)
		json.NewEncoder(newPacketBuffer).Encode(newPacket)
		log.Printf("Game created: %v", game.player1.Name)
		conn.WriteMessage(websocket.TextMessage, newPacketBuffer.Bytes())
	case commons.Join:
		log.Printf("Joining game")
		newuuid, _ := uuid.NewUUID()
		value, ok := games[newuuid.String()]
		if !ok {
			log.Printf("Error joining game")
		} else {
			log.Printf("Joining success", value)
		}
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
