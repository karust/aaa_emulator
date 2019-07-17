package sc

// OPCODES                   3.5.0.3    3.5.5.3
const (
	EnterWorldResponse       = 0x0
	HackGuardRetAddrsRequest = 0x26a
	InitialConfig            = 0x1C1 // 0x17b
	TrionConfig              = 0x0E3 // 0x24d
	AccountInfo              = 0x2A9 // 0x1bd
	ChatSpamConfig           = 0x0D4 // 0x2a8
	AccountAttributeConfig   = 0x138 // 0x6b
	LevelRestrictionConfig   = 0x113 // 0x181
	TaxItemConfig            = 0x227 // 0x050
	InGameShopConfig         = 0x0B9 // 0x190
	GameRuleConfig           = 0x136 // 0x14d
	NationMemberAdd          = 0x215
	TaxItemConfig2           = 0x00F // 0x13f
	AccountAttendance        = 0x0FD // 0x211
	GetSlotCount             = 0x10D // 0x276
	CharacterList            = 0x25B // 0x15f
	RefreshInCharacterList   = 0x19e
	ReconnectAuth            = 0x1da
	ChatMessage              = 0x1db // 0x1ad
)

/*
//                            3.0.3.0
const (
	EnterWorldResponse       = 0x0
	HackGuardRetAddrsRequest = 0x094
	InitialConfig            = 0x34
	TrionConfig              = 0x2c3
	AccountInfo              = 0x0ec
	ChatSpamConfig           = 0x281
	AccountAttributeConfig   = 0x0ba
	LevelRestrictionConfig   = 0x18a
	TaxItemConfig            = 0x1cc
	InGameShopConfig         = 0x30
	GameRuleConfig           = 0x1af
	NationMemberAdd          = 0x12d
	TaxItemConfig2           = 0x29C
	AccountAttendance        = 0x8c
	GetSlotCount             = 0x272
	CharacterList            = 0x79
	RefreshInCharacterList   = 0x1da
	ReconnectAuth            = 0x1e5
	ChatMessage              = 0x1f1
)
*/
