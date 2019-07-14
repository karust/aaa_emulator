package main

import (
	"encoding/hex"
	"io/ioutil"
	"time"

	SC "./opcodes/SC"

	"../common/packet"
)

// SCEnterWorldResponsePacket ... Provides RSA public key to client
func (sess *Session) SCEnterWorldResponsePacket(reason uint16, gm bool, token uint, port uint16) {
	w := packet.CreateEncWriter(SC.EnterWorldResponse, sess.conn.encSeq)

	w.Short(reason) // Reason
	//w.Bool(gm)      // GM, no such field in 3.5
	w.UInt(0x5933b51b)                // SC
	w.Short(port)                     // SP
	w.Long(uint64(time.Now().Unix())) // WF
	w.UInt(0xffffff4c)                // TZ

	w.Short(260) // H, Public Key Size  0401 (Should be 260)
	w.Short(260) // H, Public Key Length (in pub key) 128*2 + 4 = 260
	w.UInt(1024)
	w.Bytes(gameServer.PubModulus)
	w.Bytes(gameServer.PubExponent)

	// TODO: IP Adress and Port of client??
	w.Byte(92) // natAddr (remote client address)
	w.Byte(255)
	w.Byte(199)
	w.Byte(47)
	w.Short(49494) // natPort (remote client port)
	w.Int(1)       // authority, present in 3.5
	w.Send(sess.conn)
}

func (sess *Session) BeginGame() {
	sess.World_0x14f()
	sess.World_0x145()
}

// State Responses

// SCHackGuardRetAddrsRequestPacket ...
func (sess *Session) SCHackGuardRetAddrsRequestPacket(sendAddr bool, spMD5 bool, luaMD5 bool, modPack bool) {
	w := packet.CreateEncWriter(SC.HackGuardRetAddrsRequest, sess.conn.encSeq)
	w.Bool(sendAddr)     // send Address
	w.Bool(spMD5)        // sp Md5
	w.Bool(luaMD5)       // lua Md5
	w.String("x2ui/hud") // dir
	w.Bool(modPack)      // modPack
	w.Send(sess.conn)
}

// SCInitialConfigPacket ...
func (sess *Session) SCInitialConfigPacket() {
	// 6202b7011000617263686561676567616d652e636f6d1a007f37340f79087dcb376503dea4863c0e02e66fc7fbddae001f00 00000000 000000000000010101010100000000000000000090010001
	// 9c027b011000617263686561676567616d652e636f6d1a007f37340f79087dcb376503dea4863c0e02e66fc7fbddae001f00 00000000 000000000000010101010100000000000000000090010001
	//  00000000000000000000010101010100000000000000000090010001
	w := packet.CreateEncWriter(SC.InitialConfig, sess.conn.encSeq)
	w.String("archeagegame.com")
	fset := "7f37340f79087dcb376503dea4863c0e02e66fc7fbddae001f00" // Host
	w.HexStringL(fset)                                             // fset
	w.UInt(0)                                                      // count
	w.UInt(0)                                                      // initial Labor points

	// TODO: Initialization of this configs
	w.Bool(false) // can place house
	w.Bool(false) // can pay tax
	w.Bool(true)  // can use auction
	w.Bool(true)  // can trade
	w.Bool(true)  // can send mail
	w.Bool(true)  // can use bank
	w.Bool(true)  // can use copper

	w.Byte(0) // second  password max fail count
	w.UInt(0) // idle kick time

	w.Bool(false) // enable
	w.Byte(0)     // pcbang
	w.Byte(0)     // premium
	w.Byte(0)     // max characters
	w.Short(400)  // honorPointDuringWarPercent
	w.Byte(0)     // ucc ver
	w.Byte(1)     // member type
	w.Send(sess.conn)
}

// SCTrionConfigPacket ...
func (sess *Session) SCTrionConfigPacket(activate bool, platformURL, commerceURL string) {
	w := packet.CreateEncWriter(SC.TrionConfig, sess.conn.encSeq)
	w.Bool(activate) // Activate
	w.String(platformURL)
	w.String(commerceURL)
	// TODO: Below parameters should be URLs also
	w.Short(0) // HaveWikiUrl
	w.Short(0) // HaveCsUrl
	w.Send(sess.conn)
}

// SCAccountInfoPacket ...
func (sess *Session) SCAccountInfoPacket(payMethod, payLocation int32, payStart, payEnd uint64) {
	// 6c0e bd01 01000000 01000000 0000000000000000 5c83f66f00000000 0000000000000000 00000000
	w := packet.CreateEncWriter(SC.AccountInfo, sess.conn.encSeq)
	w.Int(payMethod)   // payMethod
	w.Int(payLocation) // payLocation
	w.Long(payStart)   // payStart
	w.Long(payEnd)     // payEnd
	w.Long(0)          // realPayTime
	w.UInt(0)          // buyPremiumCount
	w.Send(sess.conn)
}

// SCChatSpamConfigPacket ...
func (sess *Session) SCChatSpamConfigPacket() {
	//02 0100 0101010100000100000000010000010000 0000000000004040000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
	//00 0100001300 000070420600000000001644cdcc4c3f0ac803
	w := packet.CreateEncWriter(SC.ChatSpamConfig, sess.conn.encSeq)
	applyConfig := "0100001300"
	detectConfig := "000070420600000000001644cdcc4c3f0ac803"

	w.Byte(2)  // version
	w.Short(1) // report delay

	// TODO: Loop over `chatTypeGroup` (17 chats) bytes
	w.HexString("0101010100000100000000010000010000")

	// TODO: Loop over `chatGroupDelay` (17 chats) longs
	w.HexString("0000000000004040000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	w.Byte(0) // whisperChatGroup
	w.HexString(applyConfig)
	w.HexString(detectConfig)
	w.Send(sess.conn)
}

// SCAccountAttributeConfigPacket ...
func (sess *Session) SCAccountAttributeConfigPacket() {
	w := packet.CreateEncWriter(SC.AccountAttributeConfig, sess.conn.encSeq)
	w.Byte(0) //
	w.Byte(1) //
	w.Byte(0) //
	w.Send(sess.conn)
}

// SCLevelRestrictionConfigPacket ...
func (sess *Session) SCLevelRestrictionConfigPacket(searchLevel, bidLevel, postLevel, trade, mail byte) {
	w := packet.CreateEncWriter(SC.LevelRestrictionConfig, sess.conn.encSeq)
	w.Byte(searchLevel) // searchLevel
	w.Byte(bidLevel)    // bidLevel
	w.Byte(postLevel)   // postLevel
	w.Byte(trade)       // trade
	w.Byte(mail)        // mail
	// TODO: loop over `limitLevels` (15 items)
	w.HexString("0028282800002800000000000000280000") // limitLevels
	w.Send(sess.conn)
}

// SCTaxItemConfigPacket ...
func (sess *Session) SCTaxItemConfigPacket(convertRatioToAAPoint uint64) {
	// 0f08 2702 0000000000000000
	w := packet.CreateEncWriter(SC.TaxItemConfig, sess.conn.encSeq)
	w.Long(0)
	w.Send(sess.conn)
}

// SCInGameShopConfigPacket ...
func (sess *Session) SCInGameShopConfigPacket(ingameShopVersio, secondPriceType, askBuyLaborPowerPotion byte) {
	// 4509 9001 010200
	w := packet.CreateEncWriter(SC.InGameShopConfig, sess.conn.encSeq)
	w.Byte(ingameShopVersio)       // ingameShopVersion
	w.Byte(secondPriceType)        // secondPriceType
	w.Byte(askBuyLaborPowerPotion) // askBuyLaborPowerPotion
	w.Send(sess.conn)
}

// SCGameRuleConfigPacket ...
func (sess *Session) SCGameRuleConfigPacket(indunCount, conflictCount uint32) {
	// 120a 4d01 00000000 00000000
	w := packet.CreateEncWriter(SC.GameRuleConfig, sess.conn.encSeq)
	w.UInt(indunCount)
	// TODO: What does this packet do?
	/*
		for (var i = 0; i < _indunCount; i++)
		{
			stream.Write(_type); // type
			stream.Write(_pvp); // pvp
			stream.Write(_duel); // duel
		}
	*/
	w.UInt(conflictCount)
	/*
		for (var i = 0; i < _conflictCount; i++)
		{
			stream.Write(_type2); // type
			stream.Write(_peaceMin); // peaceMin
		}
	*/
	w.Send(sess.conn)
}

// SCUnknownPacket0x215 ...
func (sess *Session) SCUnknownPacket0x215(protectFaction byte, time int64) {
	// 760b 1502 01 705cba5a00000000 e2070000 03000000 1b000000 12000000 00000000
	// TODO: Parse time
	w := packet.CreateEncWriter(SC.Unknown0x215, sess.conn.encSeq)
	w.Byte(protectFaction) //protectFaction
	w.Long(1522162800)     //time
	w.UInt(2018)           //Year
	w.UInt(3)              //Month
	w.UInt(27)             //Day
	w.UInt(18)             //Hour
	w.UInt(0)              //Min
	w.Send(sess.conn)
}

// SCTaxItemConfig2Packet ...
func (sess *Session) SCTaxItemConfig2Packet(count uint32) {
	// 3a0c 3f01 00000000
	w := packet.CreateEncWriter(SC.TaxItemConfig2, sess.conn.encSeq)
	w.UInt(count)
	for i := uint32(0); i < count; i++ {
		w.UInt(0) // Type
		w.Byte(0) // declareDominion
	}
	w.Send(sess.conn)
}

func (sess *Session) World_6_BigPacket() {
	var (
		data []byte
		err  error
	)
	if sess.accountID == 1 {
		data, err = ioutil.ReadFile("etc/big_bad")
	} else {
		data, err = ioutil.ReadFile("etc/big_bad3")
	}
	if err != nil {
		panic(err)
	}
	sess.conn.Write(data)
}

// SCGetSlotCountPacket ...
func (sess *Session) SCGetSlotCountPacket(sc byte) {
	w := packet.CreateEncWriter(SC.GetSlotCount, sess.conn.encSeq)
	w.Byte(sc)
	w.Send(sess.conn)
}

// SCAccountAttendancePacket ...
func (sess *Session) SCAccountAttendancePacket(count uint) {
	// ec0f 1102 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
	w := packet.CreateEncWriter(SC.AccountAttendance, sess.conn.encSeq)
	//data := strings.Repeat("00", 248)
	//w.HexString(data)
	for i := uint(0); i < count; i++ {
		w.Long(0)
	}
	w.Send(sess.conn)
}

// SCRaceCongestionPacket ...
func (sess *Session) SCRaceCongestionPacket() {
	// // da10 5e00 0000000000000000 0000
	//w := packet.CreateEncWriter(0x14d, sess.conn.encSeq)
	w := packet.CreateEncWriter(0x5e, sess.conn.encSeq)
	w.Long(0)
	//w.Byte(0)
	w.Short(0)
	w.Send(sess.conn)
}

// SCCharacterListPacket ...
// TODO: Add character in argument
func (sess *Session) SCCharacterListPacket(last bool) {
	// f211 5f01 0100
	w := packet.CreateEncWriter(SC.CharacterList, sess.conn.encSeq)

	w.Bool(last) //LastChar
	w.Byte(0)    //TotalCount
	/*
		var charName string
		if sess.accountID == 2 {
			charName = "Rivestshamiradlemn"
		} else {
			charName = "Diffiehellman"
		}
		msg := "Hello"

		w.UInt(0x2938)     //CharID 2938
		w.String(charName) //CharName
		w.Byte(1)          //Race
		w.Byte(2)          //Gender
		w.Byte(1)          //Level
		w.UInt(370)        //HP
		w.UInt(320)        //MP
		w.UInt(179)        //zone_id
		w.UInt(101)        //F(r)actionId
		w.String(msg)      //msg
		w.UInt(0)          //type
		w.UInt(0)          //family
		w.UInt(0x1180000)  //validFlags

		//Appearance
		w.UInt(0x4d7f)                     //
		w.UInt(0x631c)                     //
		w.UInt(0x21b)                      // HairColor
		w.UInt(0)                          // twoToneHair
		w.UInt(0xd0d01)                    // twoToneFirstWidth
		w.UInt(0x4000000)                  // twoToneSecondWidth
		w.UInt(0x3da12)                    //
		w.UInt(0xc8000000)                 //
		w.UInt(0x3c1b5)                    //
		w.UInt(math.Float32bits(0x342c54)) // Float???????????

		w.HexString("cb10000000000000000000000000000000000000000000000400000000000000000000000000803f000000000000803f0000803f00000000000000000400bc01aa00000000000000000000803f0000803f0000803f8fc2353f0000803f0000803f0000803fe37b8bffafecefffafecefff584838ff00000000800000ef00ef00ee000103000000000000110000000000fe00063bb900d800ee00d400281bebe100e700f037230000000000640000000000000064000000f0000000000000002bd50000006400000000f9000000e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

		//Pers Info
		w.Short(0x1234)    //LaborPower
		w.Long(0x5bb45ff6) //lastLaborPowerModified
		w.Short(0)         //DeathCount
		w.Long(0x5bb45ff6) //deadTime
		w.UInt(0)          //rezWaitDuration
		w.Long(0x5bb45ff6) //rezTime
		w.UInt(0)          //rezPenaltyDuration
		w.Long(0x5bb45ff6) //lastWorldLeaveTime
		w.Long(0)          //moneyAmount
		w.Long(0)          //moneyAmount
		w.Short(0)         //crimePoint
		w.UInt(0)          //crimeRecord
		w.Short(0)         //crimeScore
		w.Long(0)          //deleteRequestedTime
		w.Long(0)          //transferRequestedTime
		w.Long(0)          //deleteDelay
		w.UInt(0)          //consumedLp
		w.Long(0)          //bmPoint
		w.Long(0)          //moneyAmount
		w.Long(0)          //moneyAmount
		w.Byte(0)          //autoUseAApoint
		w.UInt(1)          //prevPoint
		w.UInt(1)          //point
		w.UInt(0)          //gift
		w.Long(0x5bc39811) // updated
		w.Byte(0)          //forceNameChange
		w.UInt(0)          //highAbilityRsc
	*/
	w.Send(sess.conn)
}

func (sess *Session) World_0x14f() {
	w := packet.CreateEncWriter(0x14f, sess.conn.encSeq)
	w.UInt(1)
	w.Byte(1)
	w.UInt(0)
	w.Byte(255)
	w.UInt(0)
	w.Long(0)
	w.Long(0)
	w.Send(sess.conn)
}

func (sess *Session) World_0x145() {
	w := packet.CreateEncWriter(0x145, sess.conn.encSeq)
	w.UInt(0x2938) // charID
	w.Short(1)
	w.String("version 1\r\n")
	w.UInt(0xc)
	w.Send(sess.conn)
}

//Movement Packet
func (sess *Session) World_dd01_0x162(pack []byte, senderSess *Session) {
	print(hex.EncodeToString(pack[4:37]))
	w := packet.CreateWriter(0x1dd)
	w.Short(0)
	w.Short(0x162)
	w.Short(1)
	w.UInt24(0x66db + uint32(senderSess.accountID))
	w.Bytes(pack[7:38])

	w.Send(sess.conn)
}

func (sess *Session) MovePlayer(bc, x, y, z uint32) {
	w := packet.CreateWriter(0x1dd)
	w.Short(0)
	w.Short(0x162)
	w.Short(1)
	w.UInt24(bc)
	w.Byte(1)   // type
	w.UInt(0)   //tine?
	w.Byte(0)   //flags
	w.UInt24(x) //pos
	w.UInt24(y)
	w.UInt24(z)
	w.Short(0) //vel xyz
	w.Short(0)
	w.Short(0)
	w.Byte(0) //rot
	w.Byte(0)
	w.Byte(0)
	w.UInt24(0) // a.dm.xyz
	w.Byte(2)   //a.stace
	w.UInt24(0) //
	//w.Bytes(pack[7:38])

	w.Send(sess.conn)
}

//?        id     type time     fg pos XYZ              vel XYZ         rot XYZ   a.dmxyz   a.stace,alertness,flags   ???
//7d148400 a52b01 01   3f400800 00 011d7b c4db77 ae0703 0000 0000 1efd  00 00 39  00 00 00  02 00 00

//Display Unit
func (sess *Session) UnitState0x8d(x, y, z uint32, rx, ry, rz uint16, senderSess *Session) {
	w := packet.CreateEncWriter(0x8d, sess.conn.encSeq)
	w.UInt24(0x66db + uint32(senderSess.accountID)) // LiveID

	if senderSess.accountID == 1 {
		w.String("Diffiehellman") // name
	} else {
		w.String("RivestShamirAdlemn") // name
	}
	w.Byte(0) // type 0 - player
	if senderSess.accountID == 1 {
		w.UInt(0x2938) //charID
	} else {
		w.UInt(0x2b086) // charID
	}
	w.Long(0)   //something... "V"
	w.Short(0)  // String "master"
	w.UInt24(x) // coords
	w.UInt24(y)
	w.UInt24(z)
	w.UInt(0x3f800000) //Scale
	w.Byte(1)          //Level
	/*
		w.UInt(0x0B000000) // ModelRef

		//Inventory
		w.HexString("62450000000000000000000000000000005363000000000000000000000000000000E0600000000000000000000000514900000000000000000000004863000000000000000000000000000000000000000000000000000000000000000000000000000000000000002607000000000000000000000000000000D9360000000000000000000000000000007E4D0000425E00000000000000000000000000001802000000000000000000000000000003AA0E00000100000000000000000000000000803F0000803F0000000000000000000000000000803F000000000000803F350200000000803F000000000000803F0000000021000000000000003CDA3C3FFFCDC2FFA25F42FFA25F42FF2B250DFF4B4756FF800000FAFDE6F7DFE4553AF82622176437F5009CD934D8FE090800EBF06220BA2325F30E14FDFF02F0DA0FF325D7F516EB0A25C141E1B0D3159CCE0F0315001EFEF545E601043C1427FFED430DD5272A140023FCCB000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		w.UInt(10)      //HP
		w.UInt(10)      //MP
		w.Short(0xffff) //Points?
		w.Byte(0)       //isLooted
		w.Byte(0)       //activeWeapon
		w.Byte(0)       //learned SkillCount
		//w.UInt(0x28AB)  //type
		//w.Byte(1)       // level
		//w.UInt(0x2A00)  // type
		//w.Byte(1)       // level
		w.UInt(0)   // learnedBuffs
		w.Short(rx) //rotation xyz
		w.Short(ry)
		w.Short(rz)
		w.HexString("0800A662" +
			"01000000")
		//factionID confirm
		//ns.Write(npc.FactionId);
		w.UInt(0x65)

		w.HexString("0000000000000000")
	*/

	//	if senderSess.accountID == 1 {
	//		w.HexString(strings.Replace("00ffffffff0a000000000018017e4d0000455e0000180200000000000003dd02000000000000000000000000000000000000000000000100000000000000000000000000803f000000000000803f0000803f00000000000000005000003002aa0200000000001d000000803f0000803f0000803f0000803f0000803f0000803f0000803f000000005ab5f8ff5ab5f8ff3c2300ff603e48ff800000f5000011dc000b00000000170000000000f323000000003d0000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000088900000007d0000ffff00000001000000000001d4460000eb110000000000006500000000000000000030000000000000000000ff00000000ff00000000ff00000000ff00000000ff00000000ff00000000ff00000000ff00000000ff00000000ff00000000ff01010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001e3c32002864070001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010100000000db66000000000001010008007264310000020b127a01000000011500000000db660086b00200010100040086410000012539010000", "86b00200", "38290000", 1))
	//	} else {
	w.HexString("00ffffffff0a000000000018017e4d0000455e0000180200000000000003dd02000000000000000000000000000000000000000000000100000000000000000000000000803f000000000000803f0000803f00000000000000005000003002aa0200000000001d000000803f0000803f0000803f0000803f0000803f0000803f0000803f000000005ab5f8ff5ab5f8ff3c2300ff603e48ff800000f5000011dc000b00000000170000000000f323000000003d0000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000088900000007d0000ffff00000001000000000001d4460000eb110000000000006500000000000000000030000000000000000000ff00000000ff00000000ff00000000ff00000000ff00000000ff00000000ff00000000ff00000000ff00000000ff00000000ff01010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001e3c32002864070001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010100000000db66000000000001010008007264310000020b127a01000000011500000000db660086b00200010100040086410000012539010000")
	//	}

	w.Send(sess.conn)
}
