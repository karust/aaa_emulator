package ability

type AbilityType byte

const (
	general AbilityType = 0 + iota
	fight
	illusion
	adamant
	will
	death
	wild
	magic
	vocation
	romance
	love
	none
)

type Ability struct {
	ID    byte
	Order byte
	Exp   int
}

func (a *Ability) NewAbility() {
	a.Order = 255
}

func (a *Ability) NewAbilityID(id byte) {
	a.ID = id
	a.Order = 255
}
