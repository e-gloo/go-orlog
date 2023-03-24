package commons

const (
	Create = "create"
	Join   = "join"
)

type Packet struct {
	Command string `json:"command"`
	Data    []byte `json:"data"`
}
