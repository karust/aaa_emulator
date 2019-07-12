package main

import (
	"encoding/hex"
	"fmt"

	"../aaa_emulator/common/crypt"
)

//LOBYTE ... #define LOBYTE(w) ((BYTE)(w))
func LOBYTE(word int) byte {
	return byte(word)
}

//HIBYTE ... #define HIBYTE(w) ((BYTE)(((WORD)(w) >> 8) & 0xFF))
func HIBYTE(word int) byte {
	return byte(word>>8) & 0xff
}

func addInit(key *uint) byte {
	*key += 3121733
	n := (*key >> 14) & 0x73
	if n == 0 {
		n = 254
	}
	return byte(n)
}

func getHash(msg []byte, unkLen int) byte {
	var result int
	for i := 0; i < unkLen; i++ {
		result += int(msg[i+5]) // hash starts from Seq
		result *= 0x13
	}
	return byte(result)
}

/*
00 00 00 00 00 00 00 00 00 00 00 00 00
[17 05 00 2C 84 00 E7 13 01 01 B2 FA 07 00 00 07 62 7A FF
4C 76 1F 60 03 00 00 00 00 00 00 00 00 00 00 00 00 01 00 04 00 00]
 75 6F 3A FB 51 88 3B 5F D5 F0
*/

func obfHeader(hash byte, seq byte) (byte, byte) {
	//LOBYTE(seqAdr) = 4 * ((*(_WORD *)hashAdr >> 8) & 1 | 2 * (~(unsigned __int8)*(_WORD *)hashAdr & 1))
	resSeq := 4 * (seq&1 | 2*(^hash&1)) // ^ = flip bits

	//HIBYTE(a2) = *(_WORD *)hashAdr;
	//HIBYTE(a2) >>= 5;
	a2 := hash >> 5

	//BYTE1(a2) = HIBYTE(a2) & 2 | ~(HS >> 1) & 1 | ~(S >> 1) & 0x40 | 2 * (S & 0xC0 | 2 * (S & 2 | ~HIBYTE(a2) & 1 | seqAdr));
	a1 := a2&2 | ^(hash>>1)&1 | ^(seq>>1)&0x40 | 2*(seq&0xc0|2*(seq&2|^a2&1|resSeq))

	//LOBYTE(a2) = HS & 8 | ~(HS >> 1) & 0x40 | ((HS & 4 | (S >> 1) & 2) >> 1) | 2 * (S & 0x10 | 2 * (~(S >> 5) & 1 | 4 * (8 * ~(HS >> 4) | ~(S >> 3) & 1)));
	a2 = hash&8 | ^(hash>>1)&0x40 | ((hash&4 | (seq>>1)&2) >> 1) | 2*(seq&0x10|2*(^(seq>>5)&1|4*(8*^(hash>>4)|^(seq>>3)&1)))

	//fmt.Println(resSeq, a1, a2)
	return a2, a1 // obf HASH, SEQ
}

func main() {
	msgLen := byte(0x2c)
	//message := make([]byte, msgLen)
	//message[2] = 0x17 // seq
	//message[3] = 0x5  // header
	message := []byte{0x0, 0x0, 0x17, 0x05, 0x00, 0x2C, 0x84, 0x00, 0xE7, 0x13, 0x01, 0x01, 0xB2, 0xFA,
		0x07, 0x00, 0x00, 0x07, 0x62, 0x7A, 0xFF, 0x4C, 0x76, 0x1F, 0x60, 0x03, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x04, 0x00, 0x00}

	aesKey, _ := hex.DecodeString("42de5786e63000a307dc43861653f139")
	xorKey := uint(3797019028)
	cr := crypt.ClientCrypt(aesKey, xorKey)
	fmt.Println(cr)

	mKey := uint(5)
	someKey := uint(mKey * xorKey * xorKey & 0xffffffff)

	//For AA-5.3 HEADER = someKey ^ (ADD_INIT(thisObj, v6, someKey) + 0x75A0240D) ^ 0xB610D60;
	someByte := someKey ^ (uint(addInit(&someKey)) + 0x75A024A4) ^ 0xC3903B6A&0xffffffff
	fmt.Printf("SomeByte: %x \n", someByte)

	// Choose const
	undefConst := uint(4)
	if msgLen == 1 {
		if msgLen%3 != 0 {
			if msgLen%5 != 0 {
				if msgLen%7 != 0 {
					if msgLen%9 != 0 {
						if msgLen%11 != 0 {
							undefConst = 7
						}
						undefConst = 3
					}
					undefConst = 11
				}
				undefConst = 2
			}
			undefConst = 5
		}
		undefConst = 9
	}

	message[5] = msgLen // seq

	aint := addInit(&someKey) + 1
	aint = aint
	message[4] = getHash(message, 0x27-1)
	fmt.Println("Hash:", getHash(message, 0x27-1))

	obfHash, obfSeq := obfHeader(0x3c, 0x2c)
	fmt.Printf("Obf Hash=%x, Seq=%x \n", obfHash, obfSeq)

	//
	//AES
	//

	// Pre XOR
	bufLen := ((((mKey + 0xF) >> 31) & 0xF) + mKey + 0xF) & 0xFFFFFFF0
	fmt.Println("bufLen", bufLen)

	msgStart := (bufLen - mKey) ^ 0x3f
	fmt.Println("msgStart", msgStart)

	v20 := bufLen / undefConst
	fmt.Println("v20", v20)

	_packLen := undefConst * (bufLen / undefConst)
	fmt.Println("_PackLen:", _packLen)

	packPointer := _packLen - 1
	fmt.Println("packPointer:", packPointer)

	fmt.Println(message)

	//
	// XOR
	//
	//cr.DecXor(message, mKey, 0)
}
