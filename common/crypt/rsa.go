package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
)

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
	//keyXOR := (head^0x15a0248e)*head ^ 0x070f1f23&0xffffffff // for 3.0
	keyXOR := (head^0x15A0241F)*head ^ 0x70F1F23&0xffffffff // 3.5
	//keyXOR := (head^0xFF217A9E)*head ^ 0x1F23070F&0xffffffff // smth else
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
