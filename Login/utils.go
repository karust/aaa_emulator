package main

import (
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
)

func makeAdress(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}

func convertIPtoBytes(ip string) []byte {
	parts := strings.Split(ip, ".")
	var bytes []byte
	for _, p := range parts {
		convStr, _ := strconv.Atoi(p)
		bytes = append(bytes, byte(convStr))
	}
	return bytes
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
