package server_game

import "math"

type Power struct {
	Description string
	TokenCost   int
	Quantity    int
}

type GodActivation func(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int)

type God struct {
	Emoji       string
	Name        string
	Description string
	Priority    int
	Levels      [3]Power
	Activate    GodActivation
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
	// FIXME: multiply by die quantity (in assertFaces ?)
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
	return &God{
		Emoji:       "🏹🏹",
		Name:        "Skadi's Hunt",
		Description: "Add 🏹 to each die that rolled 🏹.",
		Priority:    4,
		Levels: [3]Power{
			{
				Description: "+1 🏹 per die",
				TokenCost:   6,
				Quantity:    1,
			},
			{
				Description: "+2 🏹 per die",
				TokenCost:   10,
				Quantity:    2,
			},
			{
				Description: "+3 🏹 per die",
				TokenCost:   14,
				Quantity:    3,
			},
		},
		Activate: ActivateSkadi,
	}
}

func ActivateSkadi(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	for dieIdx, state := range self.dice {
		if game.Dice[dieIdx].faces[state.faceIndex].kind == Arrow {
			state.quantity += god.Levels[level].Quantity
		}
	}
}

func initUllr() *God {
	return &God{
		Emoji:       "🚫🛡️",
		Name:        "Ullr's Aim",
		Description: "🏹 will ignore your opponent's 🛡️.",
		Priority:    4,
		Levels: [3]Power{
			{
				Description: "2 🏹 ignore 🛡️",
				TokenCost:   2,
				Quantity:    2,
			},
			{
				Description: "3 🏹 ignore 🛡️",
				TokenCost:   3,
				Quantity:    3,
			},
			{
				Description: "6 🏹 ignore 🛡️",
				TokenCost:   4,
				Quantity:    6,
			},
		},
		Activate: ActivateUllr,
	}
}

func ActivateUllr(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	maxToRemove := god.Levels[level].Quantity
	for dieIdx, die := range opponent.dice {
		if game.Dice[dieIdx].faces[die.faceIndex].kind == Shield {
			removed := min(die.quantity, maxToRemove)
			die.quantity -= removed
			maxToRemove -= removed
		}
	}
}

func initBaldr() *God {
	return &God{
		Emoji:       "🪖🛡️",
		Name:        "Ullr's Invulnerability",
		Description: "Add 🪖 and 🛡️ based on the current roll.",
		Priority:    4,
		Levels: [3]Power{
			{
				Description: "+1 🪖 or 🛡️",
				TokenCost:   3,
				Quantity:    1,
			},
			{
				Description: "+2 🪖 or 🛡️",
				TokenCost:   6,
				Quantity:    2,
			},
			{
				Description: "+3 🪖 or 🛡️",
				TokenCost:   9,
				Quantity:    3,
			},
		},
		Activate: ActivateBaldr,
	}
}

func ActivateBaldr(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	for dieIdx, die := range self.dice {
		kind := game.Dice[dieIdx].faces[die.faceIndex].kind
		if kind == Shield || kind == Helmet {
			die.quantity += god.Levels[level].Quantity
		}
	}
}

// FIXME: cant implement this god for now,
// too many [6] hardcoded + the flow is not clear
func initFreyja() *God {
	return &God{
		Emoji:       "📈🎲",
		Name:        "Freyja's Plenty",
		Description: "Roll additional dice this round.",
		Priority:    2,
		Levels: [3]Power{
			{
				Description: "+1 🎲",
				TokenCost:   2,
				Quantity:    1,
			},
			{
				Description: "+2 🎲",
				TokenCost:   4,
				Quantity:    2,
			},
			{
				Description: "+3 🎲",
				TokenCost:   6,
				Quantity:    3,
			},
		},
		Activate: ActivateFreyja,
	}
}

func ActivateFreyja(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	panic("not implemented")
}

func initFreyr() *God {
	return &God{
		Emoji:       "📈🎲",
		Name:        "Freyr's Gift",
		Description: "Add to the total of whichever die face is in the majority.",
		Priority:    4,
		Levels: [3]Power{
			{
				Description: "+2",
				TokenCost:   4,
				Quantity:    2,
			},
			{
				Description: "+3",
				TokenCost:   6,
				Quantity:    3,
			},
			{
				Description: "+4",
				TokenCost:   8,
				Quantity:    4,
			},
		},
		Activate: ActivateFreyr,
	}
}

func ActivateFreyr(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	// count each kind of face
	counts := map[string]int{
		"Arrow":  0,
		"Axe":    0,
		"Helmet": 0,
		"Shield": 0,
		"Thief":  0,
	}

	for dieIdx, die := range self.dice {
		kind := game.Dice[dieIdx].faces[die.faceIndex].kind
		counts[kind] += die.quantity
	}

	maxKind := ""
	maxCount := 0
	for kind, count := range counts {
		if count > maxCount {
			maxKind = kind
			maxCount = count
		}
	}

	for dieIdx, die := range self.dice {
		kind := game.Dice[dieIdx].faces[die.faceIndex].kind
		if kind == maxKind {
			die.quantity += god.Levels[level].Quantity
			break
		}
	}
}

func initIdun() *God {
	return &God{
		Emoji:       "🍎🍏",
		Name:        "Idun's Rejuvenation",
		Description: "Heal Health after the Resolution phase.",
		Priority:    7,
		Levels: [3]Power{
			{
				Description: "Heal 2 HP",
				TokenCost:   4,
				Quantity:    2,
			},
			{
				Description: "Heal 4 HP",
				TokenCost:   7,
				Quantity:    4,
			},
			{
				Description: "Heal 6 HP",
				TokenCost:   10,
				Quantity:    6,
			},
		},
		Activate: ActivateIdun,
	}
}

func ActivateIdun(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	self.health += god.Levels[level].Quantity
}

func initBrunhild() *God {
	return &God{
		Emoji:       "🪓🔥",
		Name:        "Brunhild's Fury",
		Description: "Multiply 🪓, rounded up.",
		Priority:    4,
		Levels: [3]Power{
			{
				Description: "x 1.5 🪓",
				TokenCost:   6,
				Quantity:    15,
			},
			{
				Description: "x 2 🪓",
				TokenCost:   10,
				Quantity:    20,
			},
			{
				Description: "x 3 🪓",
				TokenCost:   18,
				Quantity:    30,
			},
		},
		Activate: ActivateBrunhild,
	}
}

func ActivateBrunhild(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	countAxes := 0
	for dieIdx, die := range self.dice {
		kind := game.Dice[dieIdx].faces[die.faceIndex].kind
		if kind == Axe {
			countAxes += die.quantity
			die.quantity = 0
		}
	}

	for dieIdx, die := range self.dice {
		kind := game.Dice[dieIdx].faces[die.faceIndex].kind
		if kind == Axe {
			die.quantity = int(math.Ceil(float64(countAxes) * float64(god.Levels[level].Quantity) / float64(10.0)))
		}
	}
}

func initSkuld() *God {
	return &God{
		Emoji:       "🏹🔮",
		Name:        "Skuld's Claim",
		Description: "Destroy opponent's 🔮 for each die with an 🏹.",
		Priority:    3,
		Levels: [3]Power{
			{
				Description: "Destroy 2 🔮 per 🏹",
				TokenCost:   4,
				Quantity:    2,
			},
			{
				Description: "Destroy 3 🔮 per 🏹",
				TokenCost:   6,
				Quantity:    3,
			},
			{
				Description: "Destroy 4 🔮 per 🏹",
				TokenCost:   8,
				Quantity:    4,
			},
		},
		Activate: ActivateSkuld,
	}

}

func ActivateSkuld(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	countArrows := 0
	for dieIdx, die := range self.dice {
		kind := game.Dice[dieIdx].faces[die.faceIndex].kind
		if kind == Arrow {
			countArrows++ // no quantity here
		}
	}

	tokensToRemove := min(countArrows*god.Levels[level].Quantity, opponent.tokens)
	opponent.tokens -= tokensToRemove
}

// FIXME: cant implement this god for now,
// we need an addional user input to choose which dice to reroll ...
func initFrigg() *God {
	return &God{
		Emoji:       "🔄🎲",
		Name:        "Frigg's Sight",
		Description: "Reroll any of your opponent's dice.",
		Priority:    2,
		Levels: [3]Power{
			{
				Description: "Reroll 2 dice",
				TokenCost:   2,
				Quantity:    2,
			},
			{
				Description: "Reroll 3 dice",
				TokenCost:   3,
				Quantity:    3,
			},
			{
				Description: "Reroll 4 dice",
				TokenCost:   4,
				Quantity:    4,
			},
		},
		Activate: ActivateFrigg,
	}
}

func ActivateFrigg(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	// we need to know which dice to reroll, which is not possible yet.
	panic("not implemented")
}

// FIXME: cant implement this god for now,
// we need an addional user input to choose which dice to ban ...
func initLoki() *God {
	return &God{
		Emoji:       "🚫🎲",
		Name:        "Loki's Trick",
		Description: "Ban opponent's dice for the round.",
		Priority:    2,
		Levels: [3]Power{
			{
				Description: "Ban 1 die",
				TokenCost:   3,
				Quantity:    1,
			},
			{
				Description: "Ban 2 dice",
				TokenCost:   6,
				Quantity:    2,
			},
			{
				Description: "Ban 3 dice",
				TokenCost:   9,
				Quantity:    3,
			},
		},
		Activate: ActivateLoki,
	}
}

func ActivateLoki(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	// we need to know which dice to ban, which is not possible yet.
	panic("not implemented")
}

func initMimir() *God {
	return &God{
		Emoji:       "💔🔮",
		Name:        "Mimir's Wisdom",
		Description: "Gain 🔮 for each damage dealt to you this round.",
		Priority:    4,
		Levels: [3]Power{
			{
				Description: "+1 🔮 per 💔",
				TokenCost:   3,
				Quantity:    1,
			},
			{
				Description: "+2 🔮 per 💔",
				TokenCost:   5,
				Quantity:    2,
			},
			{
				Description: "+3 🔮 per 💔",
				TokenCost:   7,
				Quantity:    3,
			},
		},
		Activate: ActivateMimir,
	}
}

func ActivateMimir(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	// FIXME: multiply by die quantity (in assertFaces ?)
	theirAxes := opponent.assertFaces(game.Dice, func(face *ServerFace) bool { return face.kind == Axe })
	yourHelmets := self.assertFaces(game.Dice, func(face *ServerFace) bool { return face.kind == Helmet })
	axeDamageReceived := theirAxes - yourHelmets

	theirArrows := opponent.assertFaces(game.Dice, func(face *ServerFace) bool { return face.kind == Arrow })
	yourShields := self.assertFaces(game.Dice, func(face *ServerFace) bool { return face.kind == Shield })
	arrowDamageReceived := theirArrows - yourShields

	self.tokens += (axeDamageReceived + arrowDamageReceived) * god.Levels[level].Quantity
}

func initBragi() *God {
	return &God{
		Emoji:       "👌🔮",
		Name:        "Bragi's Verve",
		Description: "Gain extra 🔮 for each die that rolled 👌.",
		Priority:    4,
		Levels: [3]Power{
			{
				Description: "+2 🔮 per die",
				TokenCost:   4,
				Quantity:    2,
			},
			{
				Description: "+3 🔮 per die",
				TokenCost:   8,
				Quantity:    3,
			},
			{
				Description: "+4 🔮 per die",
				TokenCost:   12,
				Quantity:    4,
			},
		},
		Activate: ActivateBragi,
	}
}

func ActivateBragi(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	countSteals := 0
	for dieIdx, die := range self.dice {
		if game.Dice[dieIdx].faces[die.faceIndex].kind == Thief {
			countSteals++ // no quantity here
		}
	}

	self.tokens += countSteals * god.Levels[level].Quantity
}

// FIXME: hardcoded 5HP, but should be any number
func initOdin() *God {
	return &God{
		Emoji:       "🔄🔮",
		Name:        "Odin's Sacrifice",
		Description: "Sacrifice 5 of your ❤️ and gain 🔮 per HP sacrificed.",
		Priority:    7,
		Levels: [3]Power{
			{
				Description: "+3 🔮 per HP",
				TokenCost:   6,
				Quantity:    3,
			},
			{
				Description: "+4 🔮 per HP",
				TokenCost:   8,
				Quantity:    4,
			},
			{
				Description: "+5 🔮 per HP",
				TokenCost:   10,
				Quantity:    5,
			},
		},
		Activate: ActivateOdin,
	}
}

func ActivateOdin(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	// FIXME: we need an additional user input to choose how many HP to sacrifice
	hpToSacrifice := 5

	self.health -= hpToSacrifice
	self.tokens += hpToSacrifice * god.Levels[level].Quantity
}

// FIXME: cant implement this god for now,
// we need a callback mecanism to do something when the opponent spends tokens
func initVar() *God {
	return &God{
		Emoji:       "🔄🍏",
		Name:        "Var's bound",
		Description: "Each God Token spent by your opponent heals you.",
		Priority:    1,
		Levels: [3]Power{
			{
				Description: "Heal for 1HP per 🔮",
				TokenCost:   10,
				Quantity:    1,
			},
			{
				Description: "Heal for 2HP per 🔮",
				TokenCost:   14,
				Quantity:    2,
			},
			{
				Description: "Heal for 3HP per 🔮",
				TokenCost:   18,
				Quantity:    3,
			},
		},
		Activate: ActivateVar,
	}
}

func ActivateVar(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	// FIXME: we need to know how many tokens the opponent spent this round ... callback mecanism ??
	// TODO: if callback is implemented, we should to the same for each "damage this turn" effect
	panic("not implemented")
}

func initThrymr() *God {
	return &God{
		Emoji:       "📉⚡️",
		Name:        "Thrymr's Theft",
		Description: "Reduce the effect level of a God Favor invoked by the opponent this round.",
		Priority:    1,
		Levels: [3]Power{
			{
				Description: "-1 Level",
				TokenCost:   3,
				Quantity:    1,
			},
			{
				Description: "-2 Level",
				TokenCost:   6,
				Quantity:    2,
			},
			{
				Description: "-3 Level",
				TokenCost:   9,
				Quantity:    3,
			},
		},
		Activate: ActivateThrymr,
	}
}

func ActivateThrymr(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	if opponent.godChoice == nil {
		return
	}

	// FIXME: he still need to pay the initial cost ...
	if god.Levels[level].Quantity >= opponent.godChoice.level {
		opponent.godChoice = nil
	} else {
		opponent.godChoice.level -= god.Levels[level].Quantity
	}
}

// FIXME: hardcoded 5HP, but should be any number
func initTyr() *God {
	return &God{
		Emoji:       "🔒🔒",
		Name:        "Tyr's Pledge",
		Description: "Sacrifice 5 of your ❤️ to destroy an opponent's 🔮 per HP sacrificed.",
		Priority:    1,
		Levels: [3]Power{
			{
				Description: "-2 🔮 per HP",
				TokenCost:   4,
				Quantity:    2,
			},
			{
				Description: "-3 🔮 per HP",
				TokenCost:   6,
				Quantity:    3,
			},
			{
				Description: "-4 🔮 per HP",
				TokenCost:   8,
				Quantity:    4,
			},
		},
		Activate: ActivateTyr,
	}
}

func ActivateTyr(game *ServerGame, self *ServerPlayer, opponent *ServerPlayer, god *God, level int) {
	// FIXME: we need an additional user input to choose how many HP to sacrifice
	hpToSacrifice := 5

	self.health -= hpToSacrifice
	toRemove := min(hpToSacrifice*god.Levels[level].Quantity, opponent.tokens)
	opponent.tokens -= toRemove
}

func InitGods() []*God {
	// https://www.thegamer.com/assassins-creed-valhalla-orlog-god-favors/
	return []*God{
		initThor(),
		initHel(),
		initVidar(),
		initHeimdall(),
		initSkadi(),
		initUllr(),
		initBaldr(),
		// initFreyja(),
		initFreyr(),
		initIdun(),
		initBrunhild(),
		initSkuld(),
		// initFrigg(),
		// initLoki(),
		initMimir(),
		initBragi(),
		initOdin(),
		// initVar(),
		initThrymr(),
		initTyr(),
	}
}
