package client

type ClientGame struct {
	MyUsername string
}

func NewClientGame(playerUsername string) *ClientGame {
	return &ClientGame{
		MyUsername: playerUsername,
	}
}
