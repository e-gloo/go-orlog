package main

import (
	"math/rand"
	"strconv"
)

const (
	Shield = "🛡"
	Helmet = "🪖"
	Arrow  = "🏹"
	Axe    = "🪓"
	Thief  = "👌"
)

type Face struct {
	kind     string
	magic    bool
	quantity int
}

type Die struct {
	faces        [6]Face
	current_face int
	kept         bool
}

func (f *Face) String() string {
	var res = ""

	res += f.kind

	if f.quantity != 1 {
		res += strconv.Itoa(f.quantity)
	} else {
		res += ""
	}

	if f.magic {
		res += "🔮"
	} else {
		res += " "
	}

	return res + " \t"
}

func (d *Die) Face() *Face {
	return &d.faces[d.current_face]
}

func (d *Die) Roll() {
	d.current_face = rand.Intn(6)
}

func AssertDices(dices [6]Die, assert func(d *Die) bool) int {
	count := 0
	for _, die := range dices {
		if assert(&die) {
			count += 1
		}
	}
	return count
}

func AssertFaces(dices [6]Die, assert func(f *Face) bool) int {
	count := 0
	for _, die := range dices {
		if assert(die.Face()) {
			count += die.Face().quantity
		}
	}
	return count
}

func (d *Die) ResetDie() {
	d.kept = false
	for faceIdx, _ := range d.faces {
		d.faces[faceIdx].quantity = 2
	}
}

func InitDices() [6]Die {
	// Based on https://boardgamegeek.com/thread/2541060/orlog-ac-valhalla-dice
	// https://cf.geekdo-images.com/0J1WjiWz1jpny63yiVQwKA__original/img/OXm6A6qUuSZ_x3vZVCH-xWvEtXM=/0x0/filters:format(png)/pic5791191.png
	return [6]Die{
		{
			faces: [6]Face{
				{kind: Axe, magic: false, quantity: 1},
				{kind: Shield, magic: false, quantity: 1},
				{kind: Arrow, magic: true, quantity: 1},
				{kind: Axe, magic: false, quantity: 1},
				{kind: Helmet, magic: false, quantity: 1},
				{kind: Thief, magic: true, quantity: 1},
			},
			current_face: 0,
			kept:         false,
		},
		{
			faces: [6]Face{
				{kind: Axe, magic: false, quantity: 1},
				{kind: Shield, magic: true, quantity: 1},
				{kind: Arrow, magic: false, quantity: 1},
				{kind: Axe, magic: false, quantity: 1},
				{kind: Thief, magic: true, quantity: 1},
				{kind: Helmet, magic: false, quantity: 1},
			},
			current_face: 0,
			kept:         false,
		},
		{
			faces: [6]Face{
				{kind: Axe, magic: false, quantity: 1},
				{kind: Arrow, magic: true, quantity: 1},
				{kind: Thief, magic: false, quantity: 1},
				{kind: Axe, magic: false, quantity: 1},
				{kind: Helmet, magic: true, quantity: 1},
				{kind: Shield, magic: false, quantity: 1},
			},
			current_face: 0,
			kept:         false,
		},
		{
			faces: [6]Face{
				{kind: Axe, magic: false, quantity: 1},
				{kind: Shield, magic: false, quantity: 1},
				{kind: Thief, magic: true, quantity: 1},
				{kind: Arrow, magic: false, quantity: 1},
				{kind: Helmet, magic: true, quantity: 1},
				{kind: Axe, magic: false, quantity: 1},
			},
			current_face: 0,
			kept:         false,
		},
		{
			faces: [6]Face{
				{kind: Axe, magic: false, quantity: 1},
				{kind: Shield, magic: true, quantity: 1},
				{kind: Thief, magic: false, quantity: 1},
				{kind: Axe, magic: false, quantity: 1},
				{kind: Helmet, magic: false, quantity: 1},
				{kind: Arrow, magic: true, quantity: 1},
			},
			current_face: 0,
			kept:         false,
		},
		{
			faces: [6]Face{
				{kind: Axe, magic: false, quantity: 1},
				{kind: Shield, magic: true, quantity: 1},
				{kind: Thief, magic: false, quantity: 1},
				{kind: Axe, magic: false, quantity: 1},
				{kind: Arrow, magic: false, quantity: 1},
				{kind: Helmet, magic: true, quantity: 1},
			},
			current_face: 0,
			kept:         false,
		},
	}
}
