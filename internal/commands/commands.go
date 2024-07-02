package commands

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Command string

type Packet struct {
	Command Command `json:"command"`
	Data    string  `json:"data"`
}

const (
	CreateGame   Command = "create"
	JoinGame     Command = "join"
	AddPlayer    Command = "add_player"
	ChooseGods   Command = "choose_gods"
	GameStarting Command = "starting"
	PlayGods     Command = "play_gods"
	KeepDices    Command = "keep_dices"
)

const (
	SelectDices    Command = "select_dices"
	WantToPlayGods Command = "want_to_play_gods"
)

const (
	CommandOK    Command = "ok"
	CommandError Command = "error"
)

func SendPacket(conn *websocket.Conn, packet *Packet) error {
	newPacketBuffer := new(bytes.Buffer)
	err := json.NewEncoder(newPacketBuffer).Encode(packet)
	if err != nil {
		return fmt.Errorf("error encoding data: %w", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, newPacketBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("error writing message: %w", err)
	}
	return nil
}
