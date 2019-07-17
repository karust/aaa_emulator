package main

import (
	"net"
	"time"

	"../common/crypt"
	"github.com/jmoiron/sqlx"
)

// GameServer ... Holds game server data
type GameServer struct {
	Address     string
	PubModulus  []byte
	PubExponent []byte
	Timeout     time.Duration
	Online      uint
	RSA         *crypt.CryptRSA
	DB          *sqlx.DB
	ID          byte
	LoginConn   *LoginConnection
	SessConn    *SessionMap
}

// Connection ... Class for Connection
type Connection struct {
	net.Conn
	Timeout  time.Duration
	buffSize int16

	encSeq *uint8
	//proxySeq    *uint8
}

// Session ... Session
type Session struct {
	conn *Connection
	//db   *gorm.DB
	accountID    uint64
	connID       uint32
	uid          uint
	cr           *crypt.CryptAES
	kostyl       int
	alive        bool
	ingame       bool
	visibleChars []int
}
