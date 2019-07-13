package main

import (
	"net"
	"time"

	"../common/crypt"
	"github.com/jmoiron/sqlx"
)

// GameServer ... Holds game server data
type GameServer struct {
	Address string
	Timeout time.Duration
	Online  uint
	RSA     *crypt.CryptRSA
	DB      *sqlx.DB
}

// Connection ... Class for Connection
type Connection struct {
	net.Conn
	IdleTimeout time.Duration
	buffSize    int16
	encSeq      *uint8
	//proxySeq    *uint8
}

// Session ... Session
type Session struct {
	conn *Connection
	//db   *gorm.DB
	accountID    int
	uid          uint
	cr           *crypt.CryptAES
	kostyl       int
	alive        bool
	ingame       bool
	visibleChars []int
}
