package server_game

const (
	Shield = "🛡"
	Helmet = "🪖"
	Arrow  = "🏹"
	Axe    = "🪓"
	Thief  = "👌"
)

type ServerFace struct {
	kind  string
	magic bool
}

func (f ServerFace) GetKind() string {
	return f.kind
}

func (f ServerFace) IsMagic() bool {
	return f.magic
}

type ServerDie struct {
	faces [6]ServerFace
}

func (d ServerDie) GetFaces() [6]ServerFace {
	return d.faces
}

func InitDice() [6]ServerDie {
	return [6]ServerDie{
		{
			faces: [6]ServerFace{
				{kind: Axe, magic: false},
				{kind: Shield, magic: false},
				{kind: Arrow, magic: true},
				{kind: Axe, magic: false},
				{kind: Helmet, magic: false},
				{kind: Thief, magic: true},
			},
		},
		{
			faces: [6]ServerFace{
				{kind: Axe, magic: false},
				{kind: Shield, magic: true},
				{kind: Arrow, magic: false},
				{kind: Axe, magic: false},
				{kind: Thief, magic: true},
				{kind: Helmet, magic: false},
			},
		},
		{
			faces: [6]ServerFace{
				{kind: Axe, magic: false},
				{kind: Arrow, magic: true},
				{kind: Thief, magic: false},
				{kind: Axe, magic: false},
				{kind: Helmet, magic: true},
				{kind: Shield, magic: false},
			},
		},
		{
			faces: [6]ServerFace{
				{kind: Axe, magic: false},
				{kind: Shield, magic: false},
				{kind: Thief, magic: true},
				{kind: Arrow, magic: false},
				{kind: Helmet, magic: true},
				{kind: Axe, magic: false},
			},
		},
		{
			faces: [6]ServerFace{
				{kind: Axe, magic: false},
				{kind: Shield, magic: true},
				{kind: Thief, magic: false},
				{kind: Axe, magic: false},
				{kind: Helmet, magic: false},
				{kind: Arrow, magic: true},
			},
		},
		{
			faces: [6]ServerFace{
				{kind: Axe, magic: false},
				{kind: Shield, magic: true},
				{kind: Thief, magic: false},
				{kind: Axe, magic: false},
				{kind: Arrow, magic: false},
				{kind: Helmet, magic: true},
			},
		},
	}
}
