package main

import (
	"log"
	"net"
	"time"

	"../common/packet"
)

// Listen ... Listens for new connections
func (server *LoginServer) Listen() {
	listener, _ := net.Listen("tcp", server.Address)
	defer listener.Close()

	log.Println("Login Server started at:", server.Address)
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Error establishing connection: ", err)
			connection.Close()
			continue
		}
		go handleSession(connection)
	}
}

func handleSession(client net.Conn) {
	//log.Println("Login session with ", client.RemoteAddr().String())
	defer client.Close()

	session := Session{Client: client}

	var (
		err    error
		opcode uint16
		reader *packet.Reader
	)

	// handle packets
	for {
		reader, err = packet.GetPacketReader(client, 300)
		if err != nil {
			//log.Println("Error reading packet", err)
			break
		}

		opcode = reader.Short()

		switch opcode {
		// case 0x1:
		// 	err = session.RequestAuth(reader) // for 3.5.1.4 tw client
		// case 0x2:
		// 	err = session.RequestAuthTencent(reader)
		// case 0x3:
		// 	err = session.RequestAuthGameOn(reader)
		// case 0x5:
		// 	err = session.RequestAuthTrion(reader)
		case 0x6:
			err = session.ChallengeResponse2(reader) // for 3.0.3.0 client (mail.ru auth)
			// case 0x8:
			// 	err = session.CAOtpNumber(reader)
			// case 0xA:
			// 	err = session.PcCertNumber(reader)
		case 0xC:
			err = session.ListWorld(reader)
		//case 0x11:
		//	err = session.RequestAuthTW(reader)
		case 0xD:
			err = session.X2EnterWorld(reader)
		case 0xF:
			err = session.RequestReconnect(reader)
		default:
			log.Println("Unknown opcode:", opcode)
		}

		// If during session occured some error - end it
		if err != nil {
			log.Printf("%s - %s: %s\n", session.Client.RemoteAddr().String(), session.Username, err.Error())
			// Wait till client gets all messages and displays them before we close connection
			time.Sleep(time.Second * 3)
			break
		}
	}
}
