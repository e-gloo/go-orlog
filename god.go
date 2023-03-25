package main

import (
	"fmt"
)

type Power struct {
	Description string
	TokenCost   int
	Quantity    int
}

type God struct {
	Name        string
	Description string
	Priority    int
	Levels      [3]Power
	Activate    func(self *Player, opponent *Player, god *God, level int)
}

func initThor() *God {
	return &God{
		Name:        "Thor's Strike",
		Description: "Deal damage to the opponent after the resolution phase.",
		Priority:    6,
		Levels: [3]Power{
			Power{
				Description: "Deal 2 damage",
				TokenCost:   4,
				Quantity:    2,
			},
			Power{
				Description: "Deal 5 damage",
				TokenCost:   8,
				Quantity:    5,
			},
			Power{
				Description: "Deal 8 damage",
				TokenCost:   12,
				Quantity:    8,
			},
		},
		Activate: ActivateThor,
	}
}

func ActivateThor(self *Player, opponent *Player, god *God, level int) {
	opponent.health -= god.Levels[level].Quantity
}

func initHel() *God {
	return &God{
		Name:        "Hel's Grip",
		Description: "Each 🪓 damage dealt to the opponent heals you.",
		Priority:    4,
		Levels: [3]Power{
			Power{
				Description: "+1 Health per damage",
				TokenCost:   6,
				Quantity:    1,
			},
			Power{
				Description: "+2 Health per damage",
				TokenCost:   12,
				Quantity:    2,
			},
			Power{
				Description: "+3 Health per damage",
				TokenCost:   18,
				Quantity:    3,
			},
		},
		Activate: ActivateHel,
	}
}

func ActivateHel(self *Player, opponent *Player, god *God, level int) {
	axeDamageDealt := AssertFaces(self.dices, func(face *Face) bool { return face.kind == Axe }) - AssertFaces(opponent.dices, func(face *Face) bool { return face.kind == Helmet })
	self.health += axeDamageDealt * god.Levels[level].Quantity
}

func initVidar() *God {
	return &God{
		Name:        "Vidar's Might",
		Description: "Removes 🪖 from the opponent.",
		Priority:    4,
		Levels: [3]Power{
			Power{
				Description: "-2 🪖",
				TokenCost:   2,
				Quantity:    2,
			},
			Power{
				Description: "-4 🪖",
				TokenCost:   4,
				Quantity:    4,
			},
			Power{
				Description: "-6 🪖",
				TokenCost:   6,
				Quantity:    6,
			},
		},
		Activate: ActivateVidar,
	}
}

func ActivateVidar(self *Player, opponent *Player, god *God, level int) {
	toRemove := god.Levels[level].Quantity
	for idx, _ := range opponent.dices {
		if opponent.dices[idx].Face().kind == Helmet {
			removed := Min(opponent.dices[idx].Face().quantity, toRemove)
			opponent.dices[idx].Face().quantity -= removed
			toRemove -= removed
		}
	}
}

func initHeimdall() *God {
	return &God{
		Name:        "Heimdall's Watch",
		Description: "Heal health for each attack you block.",
		Priority:    4,
		Levels: [3]Power{
			Power{
				Description: "+1 Health per block",
				TokenCost:   4,
				Quantity:    1,
			},
			Power{
				Description: "+2 Health per block",
				TokenCost:   7,
				Quantity:    2,
			},
			Power{
				Description: "+3 Health per block",
				TokenCost:   10,
				Quantity:    3,
			},
		},
		Activate: ActivateHeimdall,
	}
}

func ActivateHeimdall(self *Player, opponent *Player, god *God, level int) {
	axeBlocked := Min(AssertFaces(self.dices, func(face *Face) bool { return face.kind == Helmet }), AssertFaces(opponent.dices, func(face *Face) bool { return face.kind == Axe }))
	arrowBlocked := Min(AssertFaces(self.dices, func(face *Face) bool { return face.kind == Shield }), AssertFaces(opponent.dices, func(face *Face) bool { return face.kind == Arrow }))

	self.health += god.Levels[level].Quantity * (axeBlocked + arrowBlocked)
}

func PrintGods(gods []*God) {
	for idx, god := range gods {
		fmt.Printf("[%d] %s: %s\n", idx, god.Name, god.Description)
	}
}

func InitGods() []*God {
	// https://www.thegamer.com/assassins-creed-valhalla-orlog-god-favors/
	return []*God{
		initThor(),
		initHel(),
		initVidar(),
	}
}
