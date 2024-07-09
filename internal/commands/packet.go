package commands

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Message interface {
	CreatedOrJoinedMessage | CreateOrJoinMessage | ConfigurePlayerMessage | SelectDiceMessage | WantToPlaysGodsMessage | GameStartingMessage | CreateGameMessage | JoinGameMessage | AddPlayerMessage | PlayGodsMessage | KeepDiceMessage | CommandErrorMessage
}

type Command string

type Packet struct {
	Command Command `json:"command"`
	Data    string  `json:"data"`
}

func newPacket[MessagePayload Message](command Command, data *MessagePayload) (*Packet, error) {
	packetData, err := stringifyPacketData(data)
	if err != nil {
		return nil, fmt.Errorf("error encoding data: %w", err)
	}

	return &Packet{
		Command: command,
		Data:    packetData,
	}, nil
}

func stringifyPacketData[MessagePayload Message](data *MessagePayload) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error encoding data: %w", err)
	}
	return string(jsonData), nil
}

func ParsePacketData[MessagePayload Message](packet *Packet, message *MessagePayload) error {
	err := json.Unmarshal([]byte(packet.Data), message)
	if err != nil {
		return fmt.Errorf("error decoding data: %w", err)
	}
	return nil
}

func SendPacket[MessagePayload Message](conn *websocket.Conn, command Command, data *MessagePayload) error {
	packet, err := newPacket(command, data)
	if err != nil {
		return fmt.Errorf("error constructing packet: %w", err)
	}

	newPacketBuffer := new(bytes.Buffer)
	err = json.NewEncoder(newPacketBuffer).Encode(packet)
	if err != nil {
		return fmt.Errorf("error encoding data: %w", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, newPacketBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("error writing message: %w", err)
	}

	return nil
}
