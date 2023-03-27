package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/e-gloo/orlog/orlog/commons"
	"github.com/gorilla/websocket"
)

func initPlayer() *commons.Player {
	player := &commons.Player{
		Name:     "Player",
		Health:   15,
		Token:    0,
		Dices:    commons.InitDices(),
		Position: 1,
	}

	fmt.Println("Enter your name : ")
	fmt.Scanln(&player.Name)

	// TODO: Choose gods
	// https://www.thegamer.com/assassins-creed-valhalla-orlog-god-favors/

	return player
}

func StartGame(c *websocket.Conn, join string) error {
	player := initPlayer()
	createData := &commons.CreateData{
		Uuid:   "",
		Player: player,
	}

	dataBuffer := new(bytes.Buffer)

	var command string
	if join != "" {
		command = commons.Join
		createData.Uuid = join
	} else {
		command = commons.Create
	}

	json.NewEncoder(dataBuffer).Encode(createData)
	packet := &commons.Packet{
		Command: command,
		Data:    dataBuffer.Bytes(),
	}

	packetBuffer := new(bytes.Buffer)
	json.NewEncoder(packetBuffer).Encode(packet)

	err := c.WriteMessage(websocket.TextMessage, packetBuffer.Bytes())
	if err != nil {
		fmt.Println("write:", err)
		return err
	}
	return nil
}
