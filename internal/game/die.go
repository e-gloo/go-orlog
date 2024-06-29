package game

import "math/rand"

const (
	Shield = "🛡"
	Helmet = "🪖"
	Arrow  = "🏹"
	Axe    = "🪓"
	Thief  = "👌"
)

type Face struct {
	Kind  string
	Magic bool
}

type Die struct {
	Faces        [6]Face `json:"faces"`
	Current_face int     `json:"current_faces"`
	Kept         bool    `json:"kept"`
}

func (f *Face) String() string {
	// if f.magic {
	// 	return f.kind + "🔮\t"
	// }
	return f.Kind + " \t"
}

func (d *Die) Face() *Face {
	return &d.Faces[d.Current_face]
}

func (d *Die) Roll() {
	d.Current_face = rand.Intn(6)
}

func InitDices() [6]Die {
	// Based on https://boardgamegeek.com/thread/2541060/orlog-ac-valhalla-dice
	// https://cf.geekdo-images.com/0J1WjiWz1jpny63yiVQwKA__original/img/OXm6A6qUuSZ_x3vZVCH-xWvEtXM=/0x0/filters:format(png)/pic5791191.png
	return [6]Die{
		{
			Faces: [6]Face{
				{Kind: Axe, Magic: false},
				{Kind: Shield, Magic: false},
				{Kind: Arrow, Magic: true},
				{Kind: Axe, Magic: false},
				{Kind: Helmet, Magic: false},
				{Kind: Thief, Magic: true},
			},
			Current_face: 0,
			Kept:         false,
		},
		{
			Faces: [6]Face{
				{Kind: Axe, Magic: false},
				{Kind: Shield, Magic: true},
				{Kind: Arrow, Magic: false},
				{Kind: Axe, Magic: false},
				{Kind: Thief, Magic: true},
				{Kind: Helmet, Magic: false},
			},
			Current_face: 0,
			Kept:         false,
		},
		{
			Faces: [6]Face{
				{Kind: Axe, Magic: false},
				{Kind: Arrow, Magic: true},
				{Kind: Thief, Magic: false},
				{Kind: Axe, Magic: false},
				{Kind: Helmet, Magic: true},
				{Kind: Shield, Magic: false},
			},
			Current_face: 0,
			Kept:         false,
		},
		{
			Faces: [6]Face{
				{Kind: Axe, Magic: false},
				{Kind: Shield, Magic: false},
				{Kind: Thief, Magic: true},
				{Kind: Arrow, Magic: false},
				{Kind: Helmet, Magic: true},
				{Kind: Axe, Magic: false},
			},
			Current_face: 0,
			Kept:         false,
		},
		{
			Faces: [6]Face{
				{Kind: Axe, Magic: false},
				{Kind: Shield, Magic: true},
				{Kind: Thief, Magic: false},
				{Kind: Axe, Magic: false},
				{Kind: Helmet, Magic: false},
				{Kind: Arrow, Magic: true},
			},
			Current_face: 0,
			Kept:         false,
		},
		{
			Faces: [6]Face{
				{Kind: Axe, Magic: false},
				{Kind: Shield, Magic: true},
				{Kind: Thief, Magic: false},
				{Kind: Axe, Magic: false},
				{Kind: Arrow, Magic: false},
				{Kind: Helmet, Magic: true},
			},
			Current_face: 0,
			Kept:         false,
		},
	}
}
