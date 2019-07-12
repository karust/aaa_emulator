// Thanks https://github.com/NL0bP for clarifying client encryption offset problem

package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

// Crc8 ... Checksum for server-side encrypted packets
func Crc8(packet []byte) byte {
	checksum := byte(0)
	for i := 0; i < len(packet); i++ {
		checksum = checksum * 19
		checksum += packet[i]
	}
	return checksum & 255
}

// toClientEncr help function
func add(cry *uint) byte {
	*cry += 0x2FCBD5
	n := (*cry >> 0x10) & 0xF7
	if n == 0 {
		n = 0xFE
	}
	return byte(n)
}

func makeSeq(mSeq *uint) byte {
	*mSeq += 0x2FA245
	n := byte(*mSeq>>0xE) & 0x73
	if n == 0 {
		n = 0xFE
	}
	//fmt.Println("makeSeq:", strconv.FormatInt(int64(n), 16), "  mSeq:", strconv.FormatInt(int64(*mSeq), 16))
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

// CryptRSA ...
type CryptRSA struct {
	pubKey  *rsa.PublicKey
	privKey *rsa.PrivateKey
}

// LoadRSA ... loads predifined rsa keys and returns CryptRSA object
func LoadRSA() *CryptRSA {
	pemPub := "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCjirTvOfa4UrhpApiFWitJTSGv\nSLQ4IoUk20C1q7G+Zep3PyEWtlt00RP9w/fPkaAukMuFjiDG2VTEaQe5Oczv2OjY\n5GvpYgihqqd2qCXdJheo+v4ncDI1nAWu2WuwLn0idEjoFhnI55kXhalPAzD87TnE\nD43FWgAcu5tCbL6obwIDAQAB\n-----END PUBLIC KEY-----"
	pemPriv := "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQCjirTvOfa4UrhpApiFWitJTSGvSLQ4IoUk20C1q7G+Zep3PyEW\ntlt00RP9w/fPkaAukMuFjiDG2VTEaQe5Oczv2OjY5GvpYgihqqd2qCXdJheo+v4n\ncDI1nAWu2WuwLn0idEjoFhnI55kXhalPAzD87TnED43FWgAcu5tCbL6obwIDAQAB\nAoGAS0KoXFUI6K9MpSqoJOondG61+zPSl+iu7BSoNVKDlBLTsTfQkuKtuNcEw6n8\n7z1dgUBqIJaVF91pCJArGT7zw4mBhSKbBMTkVPk3KlJUpGHVSNDuSO3hQ/7MuTD7\nbErf2OWAbpEq6e+BJknCp0yckc69+olRNwnZ1GiHVmHfKR0CQQC5GRUPfw0kMKmT\nq0NeWwS61dcgYGm8CQiRqeYfQ8dl0BvQAEeXPRz0eMCws8IHI5lDlpXLDWDD6IVt\n0U7POI3lAkEA4i/NtgBb2YmkHqyebj7JyAGdGz+uIGJAbmPhxRZ1oErEuCkIg27c\n9y3Al2c5aE/diJyUK5Lj0uyKeIzkeY8XwwJAVIvXadefugscOh49TGkItQqeE+TW\nBxSdPGO9gERmXOP9ADpQeQ1qH2TUpyHEm5wwEoZC75exvmqEH9A+TjrH3QJBAN5u\nhk0CQ1FFo2kq9k6SXpraw2ZllFZyaMxmW0MXWCt++7/jUmT2ZESL8Mazk2f6inBr\nEuda98KYLYBphdHpH0MCQBCyMlTdr4O/0GvG7iY12EG8WkhCrKqqpZa4CFw42Ho3\nKkGaXDNQ02ugSWTCLNJL7bPa25j57ncMZMRSSpcFh08=\n-----END RSA PRIVATE KEY-----"

	blockPub, _ := pem.Decode([]byte(pemPub))
	blockPriv, _ := pem.Decode([]byte(pemPriv))

	pubKey, _ := x509.ParsePKIXPublicKey(blockPub.Bytes)
	privKey, _ := x509.ParsePKCS1PrivateKey(blockPriv.Bytes)

	_rsa := new(CryptRSA)
	_rsa.pubKey = pubKey.(*rsa.PublicKey)
	_rsa.privKey = privKey
	return _rsa
}

// GetXorKey ... extracts XOR cry
func (cr *CryptRSA) GetXorKey(raw []byte) uint {
	rng := rand.Reader
	keyXORraw, err := rsa.DecryptPKCS1v15(rng, cr.privKey, raw)
	if err != nil {
		fmt.Println("Error", err)
	}

	head := binary.LittleEndian.Uint32(keyXORraw[:4])
	keyXOR := (head^0x15a0248e)*head ^ 0x070f1f23&0xffffffff

	return uint(keyXOR)
}

// GetAesKey ... extracts AES cry
func (cr *CryptRSA) GetAesKey(raw []byte) []byte {
	rng := rand.Reader
	keyAES, err := rsa.DecryptPKCS1v15(rng, cr.privKey, raw)
	if err != nil {
		fmt.Println("Error", err)
	}
	return keyAES
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

// ClientCrypt ... Decrypt packets from client
func ClientCrypt(aesKey []byte, xorKey uint) *CryptAES {
	_aes := new(CryptAES)
	_aes.aesKey = aesKey
	_aes.xorKey = xorKey * xorKey & 0xffffffff
	_aes.msgKey = map[uint8]uint8{0x30: 1, 0x31: 2, 0x33: 4, 0x34: 5, 0x35: 6, 0x36: 7, 0x37: 8, 0x38: 9, 0x39: 0xa, 0x3b: 0xc, 0x3c: 0xd, 0x3e: 0xf, 0x3f: 0x10}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(err)
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	_aes.mode = cipher.NewCBCDecrypter(block, iv)
	return _aes
}

// DecXor ... Decrypts 1st layer of XOR encryption
func (cr *CryptAES) DecXor(packet []byte, mkey uint) []byte {
	length := len(packet)
	array := make([]byte, length)
	mul := cr.xorKey * mkey
	//cr.mSeq = mul

	cry := mul ^ (uint(makeSeq(&cr.mSeq)) + 0x75a024a4) ^ 0xC3903b6a

	//fmt.Println("==cr.Seq ", strconv.FormatInt(int64(cr.Seq), 16))
	offset := 4
	if cr.Seq != 0 {
		if cr.Seq%3 != 0 {
			if cr.Seq%5 != 0 {
				if cr.Seq%7 != 0 {
					if cr.Seq%9 != 0 {
						if !(cr.Seq%11 != 0) {
							offset = 7
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
	}

	if cr.Seq == 0 {
		offset = 9
	}
	n := offset * (length / offset)
	fmt.Println("off, len, n:", offset, length, n)
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

// Decrypt ...
func (cr *CryptAES) Decrypt(data []byte, size int) []byte {
	if (len(data) - 1) < aes.BlockSize {
		panic("[CRYPT] ciphertext too short")
	}
	// CBC mode always works in whole blocks.
	if (len(data)-1)%aes.BlockSize != 0 {
		panic("[CRYPT] ciphertext is not a multiple of the block size")
	}

	xored := make([]byte, size)
	if _, ok := cr.msgKey[data[0]]; !ok {
		fmt.Println("[CRYPT] No cry in map:", data[0])
	}

	mkey := uint(size/16-1) << 4
	//fmt.Println("Mkey1", mkey)
	mkey += uint(cr.msgKey[data[0]])
	//fmt.Println("Mkey2", mkey)

	msg := data[1 : size-2]
	//fmt.Print("Mkey: ", mkey, " ")

	xored = cr.DecXor(msg, mkey)

	decr := make([]byte, size)
	cr.mode.CryptBlocks(decr, xored)
	return decr
}

func main() {
	rsa := LoadRSA()

	keys, _ := hex.DecodeString("02422645d225e0a0705c8ffd3fc07153c434dce752e614bc38734b1a470b16a5007936955658c2028784d2677203165c7c2270245a4e8a5414b0171dec91c9ac18330285b0815cae7d49e3808e3103ad95ff8664feb2498798f589494832a422a64fb8eb4a6257b2c9d678f129f28d4423f62da14a2985e7f3324030c56a6f24832ef5829af741100d9eeb06cad2067fcb358012fe7cd4b869722c1fdc2af1b7d2fe2ed091b5c2cfaafb670c4edbc336b5791aff8d8dddab8005458cd6f02d5e2e134d5df810891a85d1739e9ffc3777e9f4bedb3ab432ffbcceb14e549ab091bfddc8fdb8c9bbede43a245fc0f2bdeb869af341476567dacc18f9ff910866b7")
	aesKey := rsa.GetAesKey(keys[:128])
	xorKey := rsa.GetXorKey(keys[128:])
	fmt.Println("AES:", hex.EncodeToString(aesKey), " XOR:", xorKey)
	crypt := ClientCrypt(aesKey, xorKey)

	pack1, _ := hex.DecodeString("390134260e5f08d64f4621e003daaf0068")
	dec := crypt.Decrypt(pack1, 19)
	fmt.Println(hex.EncodeToString(dec))

	pack1, _ = hex.DecodeString("39ee448d628e8bbd7273cbe605a2b73341")
	dec = crypt.Decrypt(pack1, 19)
	fmt.Println(hex.EncodeToString(dec))

	pack1, _ = hex.DecodeString("37479a519f06a3e8a7e9de9b9d6fad50e0")
	dec = crypt.Decrypt(pack1, 19)
	fmt.Println(hex.EncodeToString(dec))

	//s, _ := hex.DecodeString("86764616e6b6875727f7c7726430fed1a1714111e2b2825222f7c297623300d4a4d7cea10a8c73ed744eaf94feb25dfcee3a718fb874a843b4250ae1c7e9a35cd769241cd2d2223f40d5c658b6b2da7716a8c6ed7249b7a1eeaa14c977f9282d5e59b9fa16a97b003ba2790405ec31399293fddf0be2e554039ad309ae2ca7c9bdb2147816c7b8b9a688f5372b1d21c33f7e5af70b59612e44095e2ec739985ea996663707d7a7774010e0b0815121f1c192623202d2a3734313e3b4845424f4c595653505d6a6764616e7b7875727fed0a0704011e1b1815122f2c292623303d3a3744414e4b4855525f5c596663606d6a7774717e7b0805020f0c191613101d2a2724212e3b3835323f4c494643405d5a5754516e6b6865727f7c797704110e1ba8050a4aeea")
	//fmt.Println(hex.EncodeToString(ToClientEncr(s)))
}
