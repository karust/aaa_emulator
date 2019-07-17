package main

import (
	"fmt"
	"io/ioutil"
	"time"

	SC "./opcodes/SC"

	"../common/packet"
)

// SCX2EnterWorldResponsePacket ... Provides RSA public key to client
func (sess *Session) SCX2EnterWorldResponsePacket(reason uint16, token uint32, port uint16) {
	// 2800000000005e7dbde0e20493652f5d000000004cffffff0401040100040000ba26f7dc17d1e5b614bc1d194c0cd4c4d010ebafde2960c89c3bcf3979cd5ade6d99634b69e24dce76e0b58b6f9e23c0734c0213b5c0291644636d24b9ff12cec858ece1faee4a050c94fa728fe01a7e50d315b7b2d444c6360cb677cd0a659f28a84108c876dc7aa59402bf6dde04405b1e57893978124e65cbaa67793604b100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100012fc7ff5c67d3 01000000
	w := packet.CreateEncWriter(SC.EnterWorldResponse, sess.conn.encSeq)
	fmt.Println("SCEnterWorldResponsePacket", token)
	w.Short(reason) // Reason
	//w.Bool(false)                     // GM, no such field in 3.5
	w.UInt(0x1)                       //w.UInt(token + 1)                 // SC
	w.Short(port)                     // SP
	w.Long(uint64(time.Now().Unix())) // WF
	w.UInt(0xffffff4c)                // TZ

	w.Short(260) // H, Public Key Size  0401 (Should be 260)
	w.Short(260) // H, Public Key Length (in pub key) 128*2 + 4 = 260
	w.UInt(1024)
	w.Bytes(gameServer.PubModulus)
	w.Bytes(gameServer.PubExponent)

	// TODO: IP Adress and Port of client??
	w.Byte(127) // natAddr (remote client address)
	w.Byte(0)
	w.Byte(0)
	w.Byte(1)
	w.Short(25375) // natPort (remote client port)
	w.Int(1)       // authority, present in 3.5
	w.Send(sess.conn)
}

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

// SCNationMemberAdd ...
func (sess *Session) SCNationMemberAdd(protectFaction byte, time int64) {
	// 760b 1502 01 705cba5a00000000 e2070000 03000000 1b000000 12000000 00000000
	// TODO: Parse time
	w := packet.CreateEncWriter(SC.NationMemberAdd, sess.conn.encSeq)
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

// SCRefreshInCharacterList ... TODO: Implement race congestion
func (sess *Session) SCRefreshInCharacterList() {
	w := packet.CreateEncWriter(SC.RefreshInCharacterList, sess.conn.encSeq)
	/*RACE_CONGESTION = {
	    LOW = 0,
	    MIDDLE = 1,
	    HIGH = 2,
	    FULL = 3,
	    PRE_SELECT_RACE_FULL = 9,
	    CHECK = 10
	}*/
	for i := 0; i < 9; i++ {
		w.Byte(0)
	}
	w.Send(sess.conn)
}

// SCReconnectAuth ... Sends cookie (token) to client to reconnect to login server
func (sess *Session) SCReconnectAuth(cookie uint32) {
	// 1e13 3102 f77f44ae
	w := packet.CreateEncWriter(SC.ReconnectAuth, sess.conn.encSeq)
	fmt.Println("SCReconnectAuth", cookie)
	w.UInt(cookie)
	w.Send(sess.conn)
}

//  SCChatMessage ... TODO: Implement arguments (character)
func (sess *Session) SCChatMessage(chtype int16, character byte, message string, ability int, langType byte) {
	w := packet.CreateEncWriter(SC.ChatMessage, sess.conn.encSeq)
	// a512 ad01 feff 0000 00000000 000000 00000000 00 00 00000000 0000 0900 476f6f642d62796521 00000000 00000000 00
	fmt.Println("SCChatMessage")
	//w.Short(uint16(chtype))
	//w.Short(0)       // chat
	w.Int(0xfffe)
	w.UInt(0)        // char.factionID
	w.UInt24(0)      // objID
	w.Int(0)         // char.ID
	w.Byte(langType) //
	w.Byte(0)        // race
	w.Int(0)         // type, factionID
	w.String("")     // Char name
	w.String(message)

	w.Byte(0)
	w.Byte(0)
	w.Byte(0)
	w.Byte(0)
	/*
		for i := 0; i < 4; i++ {
			linkType := 0
			w.Byte(0) // linkType
			if linkType > 0 {
				w.Short(0) // start
				w.Short(0) // length
				if linkType == 1 {

				} else if linkType == 3 {

				} else if linkType == 4 {

				}
			}
		}

		w.Long(0)
		w.UInt24(0) //senderObjId
		w.Int(0)    // characterId
		w.Byte(langType)
		w.Byte(0)        // CharRace
		w.Int(0)         // type
		w.String("Game") // Name
		w.String(message)


	*/

	w.Int(0)  // ability
	w.Byte(0) // option
	w.Send(sess.conn)
}
