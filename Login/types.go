package main

import (
	"net"
	"time"

	"github.com/jmoiron/sqlx"
)

// LoginServer ... Holds login server data
type LoginServer struct {
	Address          string
	Timeout          time.Duration
	MaxCharacters    byte
	CharExpanderItem byte
	GameServers      map[byte]*GameServer
	DB               *sqlx.DB
	Autologin        bool
}

// GameServer ... Holds information about game server
type GameServer struct {
	Name     string // verbose
	ID       byte   // sid
	Type     byte   // stype
	Color    byte   // scolor
	Load     byte   // 0 - low, 1 - avg, 2 - high
	IsOnline byte
	IP       []byte
	Port     uint16
}

// Session ... Holds information about login session with client
type Session struct {
	Client    net.Conn
	Username  string
	AccountID uint
}
