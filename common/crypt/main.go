// Thanks https://github.com/NL0bP for clarifying client encryption offset problem

package crypt

import (
	"crypto/aes"
	"crypto/cipher"
)

// Crc8 ... Checksum for CS & SC encrypted packets
// TODO: check how it works
func Crc8(packet []byte) byte {
	checksum := byte(0)
	for i := 0; i < len(packet); i++ {
		checksum *= 0x13
		checksum += packet[i]
	}
	return checksum
}

// toClientEncr help function
func add(cry *uint) byte {
	*cry += 0x2FCBD5
	n := (*cry >> 0x10) & 0xF7
	if n == 0 {
		return 0xFE
	}
	return byte(n)
}

func makeSeq(seq *uint) byte {
	*seq += 0x2FA245
	n := byte(*seq>>0xE) & 0x73
	if n == 0 {
		return 0xFE
	}
	return n
}

//ToClientEncr ... encrypt message to client
func ToClientEncr(packet []byte) []byte {
	length := len(packet)
	array := make([]byte, length)
	cry := uint(length ^ 522286496)
	n := 4 * int(length/4)

	for i := n - 1; i >= 0; i-- {
		val := add(&cry)
		array[i] = packet[i] ^ val
	}
	for i := n; i < length; i++ {
		val := add(&cry)
		array[i] = packet[i] ^ val
	}
	return array
}

// CryptAES ... decrypt packets encrypted with AES
type CryptAES struct {
	aesKey []byte
	xorKey uint
	msgKey map[uint8]uint8
	Seq    byte
	mSeq   uint
	mode   cipher.BlockMode
}

// DecXor ... Decrypts 1st layer of XOR encryption
func (cr *CryptAES) DecXor(packet []byte, mkey uint) []byte {
	length := len(packet)
	array := make([]byte, length)
	mul := cr.xorKey * mkey

	//var cry = mul ^ ((uint)MakeSeq() + 0x75A02411) ^ 0xCE24CEE0;         // 5.3kr
	//var cry = mul ^ ((uint)MakeSeq(ref seq) + 0x75B5BA52) ^ 0x7F7D9778;  // 3.5.5.3ru
	cry := mul ^ (uint(makeSeq(&cr.mSeq)) + 0x75A02435) ^ 0x28308228 // 3.5.0.3 NA
	//cry := mul ^ (uint(makeSeq(&cr.mSeq)) + 0x75a024a4) ^ 0xC3903b6a //3.0.3.0ru

	offset := 4
	if cr.Seq != 0 {
		if cr.Seq%3 != 0 {
			if cr.Seq%5 != 0 {
				if cr.Seq%7 != 0 {
					if cr.Seq%9 != 0 {
						if !(cr.Seq%11 != 0) {
							offset = 7
						}
					} else {
						offset = 3
					}
				} else {
					offset = 11
				}
			} else {
				offset = 2
			}
		} else {
			offset = 5
		}
	} else {
		offset = 9
	}

	n := offset * (length / offset)
	//fmt.Println("off, len, n:", offset, length, n)
	for i := n - 1; i >= 0; i-- {
		array[i] = packet[i] ^ add(&cry)
	}
	for i := n; i < length; i++ {
		array[i] = packet[i] ^ add(&cry)
	}
	cr.Seq += makeSeq(&cr.mSeq)
	cr.Seq++
	return array
}

// ClientCrypt ... Decrypt packets from client
func ClientCrypt(aesKey []byte, xorKey uint) *CryptAES {
	_aes := new(CryptAES)
	_aes.aesKey = aesKey
	_aes.xorKey = xorKey * xorKey & 0xffffffff
	_aes.msgKey = map[uint8]uint8{0x30: 1, 0x31: 2, 0x32: 3, 0x33: 4, 0x34: 5, 0x35: 6, 0x36: 7, 0x37: 8, 0x38: 9, 0x39: 0xa, 0x3a: 0xb, 0x3b: 0xc, 0x3c: 0xd, 0x3d: 0xe, 0x3e: 0xf, 0x3f: 0x10}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(err)
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	_aes.mode = cipher.NewCBCDecrypter(block, iv)
	return _aes
}

// Decrypt ...
func (cr *CryptAES) Decrypt(data []byte, size int) []byte {
	if (len(data) - 1) < aes.BlockSize {
		panic("[CRYPT] ciphertext too short")
	}
	// CBC mode always works in whole blocks.
	if (len(data)-1)%aes.BlockSize != 0 {
		panic("[CRYPT] ciphertext is not a multiple of the block size")
	}

	//if num == 0 {
	//	cr.Seq = 0
	//	cr.seq = 0
	//}
	//xored := make([]byte, size)
	//if _, ok := cr.msgKey[data[0]]; !ok {
	//	fmt.Println("[CRYPT] No cry in map:", data[0])
	//}

	mkey := uint(size/16-1) << 4
	//fmt.Println("Mkey1", mkey)
	mkey += uint(cr.msgKey[data[0]])
	//fmt.Println("Mkey2", mkey)

	msg := data[1 : size-2]
	//fmt.Print("Mkey: ", mkey, " ")

	xored := cr.DecXor(msg, mkey)

	decr := make([]byte, size)
	cr.mode.CryptBlocks(decr, xored)
	return decr
}
