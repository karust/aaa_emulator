package main

// TODO: Autoreconnect

import (
	"fmt"
	"log"
	"net"
	"time"

	"../common/packet"
	"../common/utils"
)

// LoginConnection ... Connetion of Game server with Login
type LoginConnection struct {
	address string
	conn    net.Conn
	encSeq  *uint8
	secret  string
	isAuth  bool
}

// Initialize ... check if there is Login server and establish connection
func (login *LoginConnection) Initialize(config Config) (err error) {
	log.Println("[LoginConnection] Connecting to login server...")
	login.address = utils.MakeAddress(config.Login.IP, config.Login.Port)
	conn, err := net.Dial("tcp", login.address)
	if err != nil {
		return err
	}
	conn.Close()

	login.secret = config.Login.Secret
	num := uint8(0)
	login.encSeq = &num
	accounts = NewAccountsMap()
	return nil
}

// Listen ... Listens for messages from Login server
func (login *LoginConnection) Listen() {
	log.Println("[LoginConnection] Listenining from login server...")
	login.conn, _ = net.Dial("tcp", login.address)
	defer login.conn.Close()

	// Send authentication message to Login server after 2 seconds
	go func() {
		time.Sleep(time.Second * 2)
		login.glRegister()
	}()

	var (
		err    error
		opcode uint16
		reader *packet.Reader
	)

	for {
		reader, err = packet.GetEncPacketReader(login.conn)
		if err != nil {
			log.Println("[GameConnection] Error reading packet", err)
			break
		}
		reader.Byte() // Seq
		reader.Byte() // CRC8

		opcode = reader.Short()
		switch opcode {
		case 0:
			login.lgRegister(reader)
		case 1:
			login.lgPlayerEnter(reader)
		default:
			fmt.Println("[GameConnection] No such opcode:", opcode)
		}
	}
}

//glAuthLogin ... Authentication of Game server on Login side
func (login *LoginConnection) glRegister() {
	wr := packet.CreateEncWriter(0x0, login.encSeq)
	wr.String(login.secret)
	wr.Send(login.conn)
}

//glPlayerEnter ...
func (login *LoginConnection) glPlayerEnter(connID uint32, gsID byte, result byte) {
	wr := packet.CreateEncWriter(0x1, login.encSeq)
	wr.UInt(connID)
	wr.Byte(gsID)
	wr.Byte(result)
}

//lgAuthLogin ... Result of authentication on Login server
func (login *LoginConnection) lgRegister(reader *packet.Reader) {
	result := reader.Bool()
	gameServer.ID = reader.Byte()

	if result {
		login.isAuth = true
	} else {
		login.isAuth = false
		log.Fatalln("[LoginConnection] Auth error, check secret!")
	}
}

// Login server tells that person wants to enter game server
func (login *LoginConnection) lgPlayerEnter(reader *packet.Reader) {
	accID := reader.Long()
	connID := reader.UInt()
	if _, ok := accounts.Get(accID); ok {
		login.glPlayerEnter(connID, gameServer.ID, 1)
	} else {
		accounts.Set(accID, connID)
		login.glPlayerEnter(connID, gameServer.ID, 0)
	}
}
