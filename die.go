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

type Die struct {
	faces        [6]Face
	current_face int
	kept         bool
}

func (f *Face) String() string {
	// if f.magic {
	// 	return f.kind + "🔮\t"
	// }
	return f.kind + " \t"
}

func (d *Die) Face() *Face {
	return &d.faces[d.current_face]
}

func (d *Die) Roll() {
	d.current_face = rand.Intn(6)
}

func InitDices() [6]Die {
	// Based on https://boardgamegeek.com/thread/2541060/orlog-ac-valhalla-dice
	// https://cf.geekdo-images.com/0J1WjiWz1jpny63yiVQwKA__original/img/OXm6A6qUuSZ_x3vZVCH-xWvEtXM=/0x0/filters:format(png)/pic5791191.png
	return [6]Die{
		{
			faces: [6]Face{
				{kind: Axe, magic: false},
				{kind: Shield, magic: false},
				{kind: Arrow, magic: true},
				{kind: Axe, magic: false},
				{kind: Helmet, magic: false},
				{kind: Thief, magic: true},
			},
			current_face: 0,
			kept:         false,
		},
		{
			faces: [6]Face{
				{kind: Axe, magic: false},
				{kind: Shield, magic: true},
				{kind: Arrow, magic: false},
				{kind: Axe, magic: false},
				{kind: Thief, magic: true},
				{kind: Helmet, magic: false},
			},
			current_face: 0,
			kept:         false,
		},
		{
			faces: [6]Face{
				{kind: Axe, magic: false},
				{kind: Arrow, magic: true},
				{kind: Thief, magic: false},
				{kind: Axe, magic: false},
				{kind: Helmet, magic: true},
				{kind: Shield, magic: false},
			},
			current_face: 0,
			kept:         false,
		},
		{
			faces: [6]Face{
				{kind: Axe, magic: false},
				{kind: Shield, magic: false},
				{kind: Thief, magic: true},
				{kind: Arrow, magic: false},
				{kind: Helmet, magic: true},
				{kind: Axe, magic: false},
			},
			current_face: 0,
			kept:         false,
		},
		{
			faces: [6]Face{
				{kind: Axe, magic: false},
				{kind: Shield, magic: true},
				{kind: Thief, magic: false},
				{kind: Axe, magic: false},
				{kind: Helmet, magic: false},
				{kind: Arrow, magic: true},
			},
			current_face: 0,
			kept:         false,
		},
		{
			faces: [6]Face{
				{kind: Axe, magic: false},
				{kind: Shield, magic: true},
				{kind: Thief, magic: false},
				{kind: Axe, magic: false},
				{kind: Arrow, magic: false},
				{kind: Helmet, magic: true},
			},
			current_face: 0,
			kept:         false,
		},
	}
}
