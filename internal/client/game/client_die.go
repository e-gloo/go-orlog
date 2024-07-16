package client_game

import "strconv"

type ClientFace struct {
	kind  string
	magic bool
}

func (f ClientFace) GetKind() string {
	return f.kind
}

func (f ClientFace) IsMagic() bool {
	return f.magic
}

type ClientDie struct {
	faces [6]ClientFace
}

func (d ClientDie) GetFaces() [6]ClientFace {
	return d.faces
}

func (d ClientDie) FormatDie(state *PlayerDie) string {
	face := d.faces[state.face]

	var res = ""

	res += face.kind

	if state.quantity != 1 {
		res += strconv.Itoa(state.quantity)
	} else {
		res += ""
	}

	if face.magic {
		res += "🔮"
	} else {
		res += " "
	}

	return res + " \t"
}
