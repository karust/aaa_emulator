package main

import (
	"fmt"
	"log"
	"net"

	"../common/packet"
	"../common/utils"
)

// GameConnection ... Communication with Game servers
type GameConnection struct {
	address  string
	secret   string
	encSeq   *uint8
	IDtoConn map[byte]net.Conn
}

// Initialize ... Initialize GameConnection
func (game *GameConnection) Initialize(config Config) error {
	log.Println("[GameConnection] Launching Game Connection listener...")
	game.address = utils.MakeAddress(config.GameListener.IP, config.GameListener.Port)
	game.secret = config.GameListener.Secret
	num := uint8(0)
	game.encSeq = &num
	game.IDtoConn = make(map[byte]net.Conn)
	server, err := net.Listen("tcp", game.address)
	if err != nil {
		return err
	}
	server.Close()
	return nil
}

// Listen ... Listen for game server connections
func (game *GameConnection) Listen() {
	server, _ := net.Listen("tcp", game.address)
	defer server.Close()

	for {
		connection, err := server.Accept()
		if err != nil {
			log.Println("[GameConnection] Error establishing GameConnection: ", err)
			continue
		}

		allowed, id := game.checkIP(connection)
		if !allowed {
			log.Println("[GameConnection] IP not allowed: ", err)
			connection.Close()
		}

		go game.handleConnection(connection, id)
	}
}

func (game *GameConnection) handleConnection(conn net.Conn, id byte) {
	//defer conn.Close()

	var (
		err    error
		opcode uint16
		reader *packet.Reader
	)

	for {
		reader, err = packet.GetEncPacketReader(conn)
		if err != nil {
			log.Println("[GameConnection] Error reading packet", err)
			break
		}
		reader.Byte() // hash
		reader.Byte() // CRC

		opcode = reader.Short()
		switch opcode {
		case 0:
			game.glRegisterGameServer(conn, reader, id)
		case 1:
			game.glPlayerEnter(conn, reader)
		case 2:
			game.glPlayerReconnect(conn, reader)
		case 3:
			game.glServerState(conn, reader)
		default:
			fmt.Println("[GameConnection] No such opcode:", opcode)
		}
	}
}

// Check if Game server IP is in allowed
func (game *GameConnection) checkIP(connection net.Conn) (bool, byte) {
	remoteIP := connection.RemoteAddr().(*net.TCPAddr).IP
	ip := utils.ConvertIPfromBytes(remoteIP)
	fmt.Println(ip)
	for _, gS := range loginServer.GameServers {
		if ip == gS.IP || ip == "127.0.0.1" {
			fmt.Println(gS)
			game.IDtoConn[gS.ID] = connection
			return true, gS.ID
		}
	}
	return false, 0
}

// Register Game Server
func (game *GameConnection) glRegisterGameServer(conn net.Conn, reader *packet.Reader, id byte) {
	secret := reader.String()
	if secret == game.secret {
		fmt.Println(id)
		fmt.Println(loginServer.GameServers)
		loginServer.GameServers[id].IsOnline = 1
		game.lgRegisterResponse(true, id)
		return
	}
	game.lgRegisterResponse(false, id)
	conn.Close()
	return
}

// Player entered game server
func (game *GameConnection) glPlayerEnter(conn net.Conn, reader *packet.Reader) {
	connID := reader.UInt()
	gsID := reader.Byte()
	result := reader.Byte()

	var sess *Session
	var ok bool
	if sess, ok = loginServer.Clients.Get(connID); !ok {
		fmt.Println("[glPlayerEnter] no such connID:", connID)
		return
	}

	if result == 0 {
		var gS *GameServer
		if gS, ok = loginServer.GameServers[gsID]; !ok {
			fmt.Println("[glPlayerEnter] no such game server:", gsID)
			return
		}
		sess.ACWorldCookiePacket(connID, gS)
	} else if result == 1 {
		sess.LoginDenied("Currently active", 33)
	} else {
		sess.LoginDenied("Unknown result", 25)
		fmt.Println("[glPlayerEnter] unknown result:", result)
	}
}

// Player reconnected login server
func (game *GameConnection) glPlayerReconnect(conn net.Conn, reader *packet.Reader) {

}

// Register Game Server
func (game *GameConnection) glServerState(conn net.Conn, reader *packet.Reader) {
	gsID := reader.Byte()
	gsLoad := reader.Byte()

	log.Println("[glServerState] load changed GS:", gsID, gsLoad)
	if _, ok := loginServer.GameServers[gsID]; ok {
		loginServer.GameServers[gsID].Load = gsLoad
	} else {
		log.Println("[glServerState] error")
	}
}

// Register Game Server
func (game *GameConnection) lgRegisterResponse(result bool, gsID byte) {
	wr := packet.CreateEncWriter(0x0, game.encSeq)
	wr.Bool(result)
	wr.Byte(gsID)
	wr.Send(game.IDtoConn[gsID])
}

// Player wants to enter game server
func (game *GameConnection) lgPlayerEnter(accID uint64, connID uint32, gsID byte) {
	wr := packet.CreateEncWriter(0x1, game.encSeq)
	wr.Long(accID)
	wr.UInt(connID)
	fmt.Println(gsID, game.IDtoConn)
	wr.Send(game.IDtoConn[gsID])
}

// Player reconnected login server
func (game *GameConnection) lgPlayerReconnect(conn net.Conn, token uint32, gsID byte) {
	wr := packet.CreateEncWriter(0x2, game.encSeq)
	wr.UInt(token)
	wr.Send(game.IDtoConn[gsID])
}
