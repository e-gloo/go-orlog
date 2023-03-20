package main

import "math/rand"

const (
	Shield = "🛡"
	Helmet = "🪖"
	Arrow  = "🏹"
	Axe    = "🪓"
	Thief  = "👌"
)

type Face struct {
	kind  string
	magic bool
}

type Dice struct {
	faces        [6]Face
	current_face int
	kept         bool
}

func (f *Face) String() string {
	if f.magic {
		return f.kind + "🔮\t"
	}
	return f.kind + " \t"
}

func (d *Dice) Face() *Face {
	return &d.faces[d.current_face]
}

func (d *Dice) Roll() {
	d.current_face = rand.Intn(6)
}

func InitDices() [6]Dice {
	// Based on https://boardgamegeek.com/thread/2541060/orlog-ac-valhalla-dice
	// https://cf.geekdo-images.com/0J1WjiWz1jpny63yiVQwKA__original/img/OXm6A6qUuSZ_x3vZVCH-xWvEtXM=/0x0/filters:format(png)/pic5791191.png
	return [6]Dice{
		Dice{
			faces: [6]Face{
				Face{kind: Axe, magic: false},
				Face{kind: Shield, magic: false},
				Face{kind: Arrow, magic: true},
				Face{kind: Axe, magic: false},
				Face{kind: Helmet, magic: false},
				Face{kind: Thief, magic: true},
			},
			current_face: 0,
			kept:         false,
		},
		Dice{
			faces: [6]Face{
				Face{kind: Axe, magic: false},
				Face{kind: Shield, magic: true},
				Face{kind: Arrow, magic: false},
				Face{kind: Axe, magic: false},
				Face{kind: Thief, magic: true},
				Face{kind: Helmet, magic: false},
			},
			current_face: 0,
			kept:         false,
		},
		Dice{
			faces: [6]Face{
				Face{kind: Axe, magic: false},
				Face{kind: Arrow, magic: true},
				Face{kind: Thief, magic: false},
				Face{kind: Axe, magic: false},
				Face{kind: Helmet, magic: true},
				Face{kind: Shield, magic: false},
			},
			current_face: 0,
			kept:         false,
		},
		Dice{
			faces: [6]Face{
				Face{kind: Axe, magic: false},
				Face{kind: Shield, magic: false},
				Face{kind: Thief, magic: true},
				Face{kind: Arrow, magic: false},
				Face{kind: Helmet, magic: true},
				Face{kind: Axe, magic: false},
			},
			current_face: 0,
			kept:         false,
		},
		Dice{
			faces: [6]Face{
				Face{kind: Axe, magic: false},
				Face{kind: Shield, magic: true},
				Face{kind: Thief, magic: false},
				Face{kind: Axe, magic: false},
				Face{kind: Helmet, magic: false},
				Face{kind: Arrow, magic: true},
			},
			current_face: 0,
			kept:         false,
		},
		Dice{
			faces: [6]Face{
				Face{kind: Axe, magic: false},
				Face{kind: Shield, magic: true},
				Face{kind: Thief, magic: false},
				Face{kind: Axe, magic: false},
				Face{kind: Arrow, magic: false},
				Face{kind: Helmet, magic: true},
			},
			current_face: 0,
			kept:         false,
		},
	}
}
