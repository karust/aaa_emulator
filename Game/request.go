package main

import (
	//"encoding/hex"
	"encoding/hex"
	"fmt"

	"../common/crypt"
	"../common/packet"
	"./objects"
)

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// CSX2EnterWorld ... opc=0, type(H), pFrom(I), pTo(I), accID(I), cookie(Q)
func (sess *Session) CSX2EnterWorld(reader *packet.Reader) {
	pFrom := reader.UInt()
	pTo := reader.UInt()
	accountID := reader.Long()
	cookie := reader.UInt()
	zoneID := reader.Int()
	tb := reader.Short()
	revision := reader.UInt()
	index := reader.UInt()

	fmt.Println("[GAME, X2EnterWorld]:", pFrom, pTo, accountID, cookie, zoneID, tb, revision, index)
	if connID, ok := accounts.Get(uint64(accountID)); ok {
		fmt.Println(connID, cookie, connID == cookie)
	}

	// TODO: Check if authorized here
	// TODO: Load char data from DB

	sess.accountID = accountID

	sess.SCEnterWorldResponsePacket(0, false, 0, 1250)
	sess.ChangeState(0)
}

// CSGetRsaAesKeys ...
func (sess *Session) CSGetRsaAesKeys(reader *packet.Reader, rsa crypt.CryptRSA) {
	// 0001 0000 0001 3bfff4dd97f4eb56537378978569c023b6966549928165b502ab96d6054ec1ce1fc0b099c0cbccf351a2409aef7118911b2ee6ba30186d49260958ba3df81c79f0192da5321a95687e6035854da36cf98e239955cb9785ffae7ac1be00625c355b90b1f3324048e20a4b3571eaec0dfbdb47cea1b2949739f451cf8f22aedb37 26bc9065e4b3e4ac1ef5b8b4eaf79e3ee522b139e78bd62527986a1e535506bc49bd9c5f44fede05391122f5f9604470843a50d1048affd9feaf47e910054f7a0ddc6f2418d5b206e54ef15dec5cb58c9af6265b23c4c60e1d91a796cc0f856f32d77e41ed879040e16acda68f46de719a3d6d38ed7b8a7be47aebb2b723f1a0
	//reader.Short() // Unknown, always = 355 (0x6301)
	reader.Short() // lenXOR ?
	reader.Short() // ?
	reader.Short() // lenAES ?

	encAES := reader.BytesLen(128)
	encXOR := reader.BytesLen(128)

	aesKey := rsa.GetAesKey(encAES)
	xorKey := rsa.GetXorKey(encXOR)

	sess.cr = crypt.ClientCrypt(aesKey, xorKey)

	fmt.Println("[GAME, getKeys]: AES: ", hex.EncodeToString(aesKey), ", XOR: ", xorKey)

	// Following sequence of responses are for characters in menu?
	sess.SCGetSlotCountPacket(0)
	// TODO: Changing time in `SCAccountInfoPacket`
	sess.SCAccountInfoPacket(1, 1, 0, 0x6ff6835c)
	// TODO: Load account here

	sess.SCAccountAttendancePacket(31)
	// da10 5e00 0000000000000000 0000
	sess.SCRaceCongestionPacket()

	sess.SCCharacterListPacket(true)
	//sess.SCCharacterListPacket()

	// a512 ad01 feff0000000000000000000000000000000000000000000900476f6f642d62796521000000000000000000
	h1, _ := hex.DecodeString("3100dd05f5325dc06f9e3101d2a2724212e3b3835323f4c494643405d5a5754c16a1d9e9330a95bef2514010e0b0815180b0e7")
	sess.conn.Write(h1)
	h2, _ := hex.DecodeString("3100dd0598335dc06f9e3101d2a2724212e3b3835323f4c494643405d5a5754c16a1d9e9330a95bef2514010e0b0815180b0e7")
	sess.conn.Write(h2)
}

// CSCreateCharacter ... Creation of character
func (sess *Session) CSCreateCharacter(reader *packet.Reader) {
	// 5e62 af01 06007177653132330601954e0000cd3d00000000000000000000320200002f0200000000000003ba02000000000000000000000000000000000000000000000900000000000000000000000000803f000000000000803f0000803f00000000000000001000000c0300000000000000000000803f0000803f0000803f0000803f0000803f0000803f0000803fff9538ffc9bc01ffc9bc01ff240005ff00000000800000d6e2d4c83c34a5dfe3641c64b99c649c6400000005f40e00eb000cef07a4fedbd2dc649cf61bc73dd7f2dcd5009cdea99c0000df0d24649c0500003fe99cca5b009c00319c6400f2646464d0000000002b64000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060d0d01ffffffff1fbf000000
	name := reader.String()
	race := reader.Byte()
	gender := reader.Byte()
	items := make([]uint32, 7)
	for i := 0; i < 7; i++ {
		items[i] = reader.UInt()
	}

	// Custom model
	charModel := objects.CharacterModel{}
	charModel.Parse(reader)

	ability1 := reader.Byte()
	ability2 := reader.Byte()
	ability3 := reader.Byte()
	level := reader.Byte()
	introZoneID := reader.Int()

	fmt.Println(name, race, gender, items, ability1, ability2, ability3, level, introZoneID)
}

// CSSecurityReport ... Some message related to securuty violation?
func (sess *Session) CSSecurityReport(reader *packet.Reader) {
	// 02 00000000 00000000 cc994900 0000
	// TODO: What tells this packet?
	srType := reader.Byte()

	if srType == 1 {
		unkUInt1 := reader.UInt()
		unkLong := reader.Long()
		str := reader.String()
		unkUInt2 := reader.UInt()
		unkByte := reader.Byte()
		fmt.Println("CSSecurityReportPacket 1:", unkUInt1, unkLong, str, unkUInt2, unkByte)
	} else if srType == 2 {
		unkUInt1 := reader.UInt()
		unkUInt2 := reader.UInt()
		fmt.Println("CSSecurityReportPacket 2:", unkUInt1, unkUInt2)
	} else if srType == 3 {
		objID := reader.Byte()
		unkSh := reader.Short()
		fmt.Println("CSSecurityReportPacket 3:", objID, unkSh)
	}
}

/*
<packet type="0x09B" level="0x05" desc="CS_PREMIUM_SERVICE_MSG">
<chunk type="d" name="stage"/>
</packet>
*/
// CSPremiumServiceMSG ... TODO: ?
func (sess *Session) CSPremiumServiceMSG(reader *packet.Reader) {
	// 01000000 000052cc000053cc000000
	stage := reader.Int()
	//Connection.SendPacket(new SCAccountWarnedPacket(2, "Premium ..."));
	fmt.Println("CSPremiumServiceMSG:", stage)
}

/*
<packet type="0x0BF" level="0x05" desc="CS_LEAVE_WORLD">
<chunk type="w" name="pSize"/>
<chunk type="w" name="pLevel"/>
<chunk type="w" name="pHash"/>
<chunk type="w" name="pType"/>
</packet>
*/
// CSLeaveWorld ... Report kind of leaving from game
func (sess *Session) CSLeaveWorld(reader *packet.Reader) {
	leaveType := reader.Byte()

	switch leaveType {
	case 0:
		fmt.Println("CSLeaveWorld, Exit game:", leaveType)
	case 1:
		fmt.Println("CSLeaveWorld, Choose characters:", leaveType)
		// connection.SendPacket(new SCPrepareLeaveWorldPacket(10000, type, false));
		// connection.LeaveTask = new LeaveWorldTask(connection, type);
		// TaskManager.Instance.Schedule(connection.LeaveTask, TimeSpan.FromSeconds(10));
	case 2:
		// if (connection.State == GameState.Lobby)
		// {
		// 	var gsId = AppConfiguration.Instance.Id;
		// 	LoginNetwork
		// 		.Instance
		// 		.GetConnection()
		// 		.SendPacket(new GLPlayerReconnectPacket(gsId, connection.AccountId, connection.Id));
		// }
		sess.SCReconnectAuth(0x0)
		fmt.Println("CSLeaveWorld, Choose server:", leaveType)
	default:
		fmt.Println("CSLeaveWorld, Unknown type:", leaveType)
	}
}

// CSRefreshInCharacterList ... TODO: ?
func (sess *Session) CSRefreshInCharacterList(reader *packet.Reader) {
	fmt.Println("CSRefreshInCharacterList")
	sess.SCRefreshInCharacterList()
}

/*
func (sess *Session) OnMovement(pack []byte) {
	reader := packet.CreateReader(pack)
	reader.Byte()
	reader.Byte()
	reader.Short() // op :=
	reader.Int24() // bc :=
	reader.Byte()  // _type :=
	reader.Int()   // time :=
	reader.Byte()  // flags :=

	posX := reader.Int24()
	posY := reader.Int24()
	posZ := reader.Int24()
	velX := reader.Short()
	velY := reader.Short()
	velZ := reader.Short()
	rotX := reader.Byte()
	rotY := reader.Byte()
	rotZ := reader.Byte()
	aDmX := reader.Byte()
	aDmY := reader.Byte()
	aDmZ := reader.Byte()
	reader.Byte() // aStace :=
	reader.Byte() // aAlertness :=
	reader.Byte() // aFlags :=
	//fmt.Println(posX, posY, posZ)
	//fmt.Println(rotX, rotY, rotZ)
	// Escaping compiling error
	aDmX = aDmX
	aDmY = aDmY
	aDmZ = aDmZ
	velX = velX
	velY = velY
	velZ = velZ

	go sess.MovementReply(pack, uint32(posX), uint32(posY), uint32(posZ), uint16(rotX), uint16(rotY), uint16(rotZ))
}


func (sess *Session) MovementReply(pack []byte, x, y, z uint32, rx, ry, rz uint16) {
	for i := range sessions {
		if sessions[i].alive && sessions[i].ingame && sess != sessions[i] {
			if !intInSlice(sess.accountID, sessions[i].visibleChars) {
				sessions[i].visibleChars = append(sessions[i].visibleChars, sess.accountID)
				sessions[i].UnitState0x8d(x, y, z, rx, ry, rz, sess)
			}
			sessions[i].World_dd01_0x162(pack, sess)
		}
	}
}
*/
