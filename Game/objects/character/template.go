package character

import "../world"

type CharacterTemplate struct {
	Race                   Race
	Gender                 Gender
	ModelID                uint
	ZoneID                 uint
	FactionID              uint
	ReturnDictrictID       uint
	ResurrectionDictrictID uint
	Position               *world.Point
	Items                  []uint
	Buffs                  []uint
	NumInventorySlot       byte
	NumBankSlot            int16
}

// NewCharacterTemplate ... Creates CharacterTemplate object
func NewCharacterTemplate() {
	chTemp := CharacterTemplate{}
	chTemp.Position = world.NewPoint(0, 0, 0)
	chTemp.Items = make([]uint, 7)
	chTemp.Buffs = make([]uint, 0)
}
