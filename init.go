package main

const (
	Shield = "shield"
	Helmet = "helmet"
	Arrow = "arrow"
	Axe = "axe"
	Thief = "thief"
)

type Face struct {
	kind string
	magic bool
}

type Dice struct {
	faces [6]Face
	kept bool
}

func Init() [6]Dice {
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
			kept: false,
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
			kept: false,
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
			kept: false,
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
			kept: false,
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
			kept: false,
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
			kept: false,
		},
	}
}
