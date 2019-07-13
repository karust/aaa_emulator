package main

import (
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
