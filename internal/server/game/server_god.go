package server_game

type Power struct {
	Description string
	TokenCost   int
	Quantity    int
}

type God struct {
	Emoji       string
	Name        string
	Description string
	Priority    int
	Levels      [3]Power
	Activate    func(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int)
}

func initThor() *God {
	return &God{
		Emoji:       "⚡⚡",
		Name:        "Thor's Strike",
		Description: "Deal damage to the opponent after the resolution phase.",
		Priority:    6,
		Levels: [3]Power{
			{
				Description: "Deal 2 damage",
				TokenCost:   4,
				Quantity:    2,
			},
			{
				Description: "Deal 5 damage",
				TokenCost:   8,
				Quantity:    5,
			},
			{
				Description: "Deal 8 damage",
				TokenCost:   12,
				Quantity:    8,
			},
		},
		Activate: ActivateThor,
	}
}

func ActivateThor(_ *ServerGame, _ *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	opponent.health -= god.Levels[level].Quantity
}

func initHel() *God {
	return &God{
		Emoji:       "🪓❤️",
		Name:        "Hel's Grip",
		Description: "Each 🪓 damage dealt to the opponent heals you.",
		Priority:    4,
		Levels: [3]Power{
			{
				Description: "+1 Health per damage",
				TokenCost:   6,
				Quantity:    1,
			},
			{
				Description: "+2 Health per damage",
				TokenCost:   12,
				Quantity:    2,
			},
			{
				Description: "+3 Health per damage",
				TokenCost:   18,
				Quantity:    3,
			},
		},
		Activate: ActivateHel,
	}
}

func ActivateHel(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	// FIXME: multiply by die quantity
	yourAxes := self.assertFaces(game.Dice, func(face *ServerFace) bool { return face.kind == Axe })
	theirHelmets := opponent.assertFaces(game.Dice, func(face *ServerFace) bool { return face.kind == Helmet })
	axeDamageDealt := yourAxes - theirHelmets
	self.health += axeDamageDealt * god.Levels[level].Quantity
}

func initVidar() *God {
	return &God{
		Emoji:       "🚫🪖",
		Name:        "Vidar's Might",
		Description: "Removes 🪖 from the opponent.",
		Priority:    4,
		Levels: [3]Power{
			{
				Description: "-2 🪖",
				TokenCost:   2,
				Quantity:    2,
			},
			{
				Description: "-4 🪖",
				TokenCost:   4,
				Quantity:    4,
			},
			{
				Description: "-6 🪖",
				TokenCost:   6,
				Quantity:    6,
			},
		},
		Activate: ActivateVidar,
	}
}

func ActivateVidar(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	maxToRemove := god.Levels[level].Quantity
	for dieIdx, die := range opponent.dice {
		if game.Dice[dieIdx].faces[die.faceIndex].kind == Helmet {
			removed := min(die.quantity, maxToRemove)
			die.quantity -= removed
			maxToRemove -= removed
		}
	}
}

func initHeimdall() *God {
	return &God{
		Emoji:       "🛡️🪖",
		Name:        "Heimdall's Watch",
		Description: "Heal health for each attack you block.",
		Priority:    4,
		Levels: [3]Power{
			{
				Description: "+1 Health per block",
				TokenCost:   4,
				Quantity:    1,
			},
			{
				Description: "+2 Health per block",
				TokenCost:   7,
				Quantity:    2,
			},
			{
				Description: "+3 Health per block",
				TokenCost:   10,
				Quantity:    3,
			},
		},
		Activate: ActivateHeimdall,
	}
}

func ActivateHeimdall(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	// FIXME: multiply by die quantity
	theirArrows := opponent.assertFaces(game.Dice, func(face *ServerFace) bool { return face.kind == Arrow })
	yourShields := self.assertFaces(game.Dice, func(face *ServerFace) bool { return face.kind == Shield })
	arrowBlocked := min(yourShields, theirArrows)

	theirAxes := opponent.assertFaces(game.Dice, func(face *ServerFace) bool { return face.kind == Axe })
	yourHelmets := self.assertFaces(game.Dice, func(face *ServerFace) bool { return face.kind == Helmet })
	axeBlocked := min(yourHelmets, theirAxes)

	self.health += god.Levels[level].Quantity * (axeBlocked + arrowBlocked)
}

func initSkadi() *God {
	return nil
}

func initUllr() *God {
	return nil
}

func initBaldr() *God {
	return nil
}

func initFreyja() *God {
	return nil
}

func initFreyr() *God {
	return nil
}

func initIdun() *God {
	return nil
}

func initBrunhild() *God {
	return nil
}

func initSkuld() *God {
	return nil
}

func initFrigg() *God {
	return nil
}

func initLoki() *God {
	return nil
}

func initMimir() *God {
	return nil
}

func initBragi() *God {
	return nil
}

func initOdin() *God {
	return nil
}

func initVar() *God {
	return nil
}

func initThrymr() *God {
	return nil
}

func initTyr() *God {
	return nil
}

func InitGods() []*God {
	// https://www.thegamer.com/assassins-creed-valhalla-orlog-god-favors/
	return []*God{
		initThor(),
		initHel(),
		initVidar(),
		initHeimdall(),
		// initSkadi(),
		// initUllr(),
		// initBaldr(),
		// initFreyja(),
		// initFreyr(),
		// initIdun(),
		// initBrunhild(),
		// initSkuld(),
		// initFrigg(),
		// initLoki(),
		// initMimir(),
		// initBragi(),
		// initOdin(),
		// initVar(),
		// initThrymr(),
		// initTyr(),
	}
}
