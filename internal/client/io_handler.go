package client

import (
	"fmt"
)

type IOHandler interface {
	// DisplayMessage(string)
	// ReadInput(*string) error
	Send(any)
}

type TermHandler struct {
}

func (th *TermHandler) DisplayMessage(msg string) {
	fmt.Println(msg)
}

func (th *TermHandler) ReadInput(buff *string) error {
	_, err := fmt.Scanln(buff)
	if err != nil && err.Error() != "unexpected newline" {
		return err
	}
	return nil
}
