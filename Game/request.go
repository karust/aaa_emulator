package main

import (
	//"encoding/hex"
	"encoding/hex"
	"fmt"

	"../common/crypt"
	"../common/packet"
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
	ttype := reader.Short()
	pFrom := reader.Int()
	pTo := reader.Int()
	accountID := reader.Int()
	cookie := reader.Long()
	zoneID := reader.Int()
	tb := reader.Short()
	revision := reader.Int()
	index := reader.Int()

	fmt.Println("[GAME, X2EnterWorld]:", ttype, pFrom, pTo, accountID, cookie, zoneID, tb, revision, index)

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
