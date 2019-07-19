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
	address  string
	conn     net.Conn
	encSeq   *uint8
	secret   string
	isAuth   bool
	confSave Config
}

// Initialize ... check if there is Login server and establish connection
func (login *LoginConnection) Initialize(config Config) (err error) {
	log.Println("[LoginConnection] Connecting to login server...")
	login.address = utils.MakeAddress(config.Login.IP, config.Login.Port)
	conn, err := net.Dial("tcp", login.address)
	if err != nil {
		return err
	}
	//defer conn.Close()

	login.secret = config.Login.Secret
	login.conn = conn
	num := uint8(0)
	login.encSeq = &num
	accounts = NewAccountsMap()
	login.confSave = config

	// Send authentication message to Login server after 2 seconds
	go func() {
		time.Sleep(time.Second * 3)
		login.glRegister()
	}()

	return nil
}

// Listen ... Listens for messages from Login server
func (login *LoginConnection) Listen() {
	//login.conn, _ = net.Dial("tcp", login.address)
	defer login.conn.Close()
	log.Println("[LoginConnection] Connected to login server at:", login.address)

	var (
		err    error
		opcode uint16
		reader *packet.Reader
	)

	for {
		reader, err = packet.GetEncPacketReader(login.conn)
		if err != nil {
			log.Println("[LoginConnection] Login server is down, reconnection after 2 secs...")
			time.Sleep(time.Second * 2)
			login.Initialize(login.confSave)
			continue
		}
		reader.Byte() // Seq
		reader.Byte() // CRC8

		opcode = reader.Short()
		switch opcode {
		case 0:
			login.lgRegister(reader)
		case 1:
			login.lgPlayerEnter(reader)
		case 2:
			login.lgPlayerReconnect(reader)
		default:
			fmt.Println("[GameConnection] No such opcode:", opcode)
		}
	}
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

//glAuthLogin ... Authentication of Game server on Login side
func (login *LoginConnection) glRegister() {
	wr := packet.CreateEncWriter(0x0, login.encSeq)
	wr.String(login.secret)
	wr.Send(login.conn)
}

// Login server tells that person wants to enter game server
func (login *LoginConnection) lgPlayerEnter(reader *packet.Reader) {
	accID := reader.Long()
	connID := reader.UInt()
	fmt.Println("lgPlayerEnter", accID, connID)
	if _, ok := accounts.Get(accID); ok {
		login.glPlayerEnter(connID, gameServer.ID, 1)
	} else {
		accounts.Set(accID, connID)
		// TODO: kostyl &Session{}
		gameServer.SessConn.Set(connID, &Session{})
		login.glPlayerEnter(connID, gameServer.ID, 0)
	}
}

func (login *LoginConnection) lgPlayerReconnect(reader *packet.Reader) {
	connID := reader.UInt()
	// TODO: kostyl
	fmt.Println("lgPlayerReconnect", connID)
	//sess.SCReconnectAuth(123)
	//if sess, ok := gameServer.SessConn.Get(connID); ok {
	//	fmt.Println(sess.connID, sess.uid)
	//	sess.SCReconnectAuth(sess.connID)
	//}
}

//glPlayerEnter ...
func (login *LoginConnection) glPlayerEnter(connID uint32, gsID byte, result byte) {
	wr := packet.CreateEncWriter(0x1, login.encSeq)
	wr.UInt(connID)
	wr.Byte(gsID)
	wr.Byte(result)
	wr.Send(login.conn)
}

func (login *LoginConnection) glPlayerReconnect(gsID byte, accID uint64, connID uint32) {
	wr := packet.CreateEncWriter(0x2, login.encSeq)
	wr.Byte(gsID)
	wr.Long(accID)
	wr.UInt(connID)
	wr.Send(login.conn)
}
