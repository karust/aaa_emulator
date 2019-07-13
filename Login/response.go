package main

import "../common/packet"

// Responses

func (sess *Session) loginDenied(responseVerbose string, reason byte) error {
	serial := packet.CreateWriter(12)
	serial.Byte(reason)
	serial.Short(0)
	serial.String(responseVerbose)
	serial.Send(sess.Client)
	err := serial.Send(sess.Client)
	return err
}

func (sess *Session) joinResponse() error {
	serial := packet.CreateWriter(0)
	serial.Byte(1)       // AuthID
	serial.Short(0)      // Reason
	serial.Long(4719366) // "afs" from archerage
	err := serial.Send(sess.Client)
	return err
}

func (sess *Session) authResponse() error {
	serial := packet.CreateWriter(3)
	if sess.Username == "admin" {
		serial.Long(1) // AccountID
	} else {
		serial.Long(2)
	}
	serial.String("FE4E6C87FB6C1625CA3832B478E2E2F0")
	serial.Byte(5)
	err := serial.Send(sess.Client)
	if err != nil {
		return err
	}
	return nil
}

func (sess *Session) worldListPacket() error {
	serial := packet.CreateWriter(8)
	serial.Byte(byte(len(loginServer.GameServers)))
	for i := range loginServer.GameServers {
		serial.Byte(loginServer.GameServers[i].ID)
		serial.Byte(loginServer.GameServers[i].Type)
		serial.Byte(loginServer.GameServers[i].Color)
		serial.String(loginServer.GameServers[i].Name)
		serial.Byte(loginServer.GameServers[i].IsOnline)
		serial.Byte(loginServer.GameServers[i].Load)
		serial.Byte(3) // ?
		serial.Byte(0) // Humans
		serial.Byte(3) // ?
		serial.Byte(0) // Dwarfs
		serial.Byte(0) // Elfs
		serial.Byte(0) // Hari...
		serial.Byte(0) // Cats
		serial.Byte(3) // ?
		serial.Byte(0) // Warlocks
	}
	serial.Byte(0) // Char Count
	err := serial.Send(sess.Client)
	//should be characters info

	return err
}

func (sess *Session) worldCookiePacket(cookie uint32, gameServer *GameServer) error {
	serial := packet.CreateWriter(0xA)
	serial.UInt(cookie)
	serial.Byte(gameServer.IP[3])
	serial.Byte(gameServer.IP[2])
	serial.Byte(gameServer.IP[1])
	serial.Byte(gameServer.IP[0])
	serial.Short(gameServer.Port)
	serial.Long(0)
	serial.Long(0)
	serial.Short(0)

	err := serial.Send(sess.Client)
	return err
}
