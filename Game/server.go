package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"

	"../common/crypt"
	"../common/packet"
	CS "./opcodes/CS"
)

var sessions []*Session
var accounts *AccountsMap

// Listen ... Listens for new connections
func (s *GameServer) Listen() error {
	s.RSA = crypt.LoadRSA()

	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Println("Game server started at", s.Address)

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Println("Cannot establish connection:", err)
			continue
		}

		conn := &Connection{
			Conn:    client,
			Timeout: s.Timeout,
		}

		// TODO: What is this? Kostyl?
		var num uint8
		conn.encSeq = &num
		go handleSession(conn, s)
	}
}

func handleSession(conn *Connection, serv *GameServer) {
	log.Println("Game session with:", conn.RemoteAddr().String())

	// TODO: Remove or rework this
	sess := &Session{conn: conn, kostyl: 1, alive: true, ingame: false}
	defer func() {
		conn.Close()
		sess.alive = false
		sess.ingame = false
	}()

	// TODO: Rework sessions, add lock, create class with methods
	//sessions = append(sessions, sess)

	var (
		err     error
		opcode  uint16
		subtype byte
		reader  *packet.Reader
	)

	for {
		reader, err = packet.GetPacketReader(conn, 30)
		if err != nil {
			log.Println("Error reading packet", err)
			break
		}
		reader.Byte()
		subtype = reader.Byte()

		switch subtype {
		case 1:
			reader.Byte() // crc
			reader.Byte() // counter
			opcode = reader.Short()

			switch opcode {
			case CS.X2EnterWorld:
				sess.CSX2EnterWorld(reader)
			case CS.GetRsaAesKeys:
				sess.CSGetRsaAesKeys(reader, *serv.RSA)
			case 0x12d:
				fmt.Println("[WORLD] Unknows packet:", opcode)
			default:
				fmt.Println("[WORLD] No opcode found:", opcode)
			}

		case 2:
			opcode = reader.Short()
			switch opcode {
			case 1:
				sess.FinishState(reader)
			case 18: // PingPacket
				go sess.Pong(reader)
			default:
				fmt.Println("[PROXY] No opcode found:", opcode)
			}

		case 3:
			opcode = reader.Short()
			//reader := &packet.PacketReader{Pack: packBuf[4 : plen+4], Offset: 0}
			switch opcode {
			default:
				fmt.Println("[COMPRSSED] No opcode found:", opcode)
			}

		case 4:
			opcode = reader.Short()
			//reader := &packet.PacketReader{Pack: packBuf[4 : plen+4], Offset: 0}
			switch opcode {
			default:
				fmt.Println("[COMPR-MULTI] No opcode found:", opcode)
			}

		case 5:
			decr := sess.cr.Decrypt(reader.Pack[2:], len(reader.Pack))
			reader = packet.CreateReader(decr)

			reader.Byte() // seq
			reader.Byte() // hash

			opcode = reader.Short() //binary.LittleEndian.Uint16(decr[2:4])

			switch opcode {
			case CS.CreateCharacter:
				sess.CSCreateCharacter(reader)
			case CS.SecurityReport:
				sess.CSSecurityReport(reader)
			case CS.PremiumServiceMSG:
				sess.CSPremiumServiceMSG(reader)
			case CS.LeaveWorld:
				sess.CSLeaveWorld(reader)
			case CS.RefreshInCharacterList:
				sess.CSRefreshInCharacterList(reader)
			default:
				fmt.Printf("[%x] %v\n", opcode, hex.EncodeToString(decr))
			}
		default:
			fmt.Println("[GAME] No such subtype:", subtype)
		}
	}
}
