package utils

import (
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
)

// MakeAdress ... Concatenates IP and Port to make Addres
func MakeAdress(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}

// ConvertIPtoBytes ... Converts string IP value to Byte array
func ConvertIPtoBytes(ip string) []byte {
	parts := strings.Split(ip, ".")
	var bytes []byte
	for _, p := range parts {
		convStr, _ := strconv.Atoi(p)
		bytes = append(bytes, byte(convStr))
	}
	return bytes
}

// BoolToByte ... Converts Boolean value to Byte
func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

// RandomHex ... Generates random Hex string of given length
func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
