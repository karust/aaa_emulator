package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"

	"../common/crypt"
	"../common/packet"
	CS "./opcodes/CS"
)

var sessions []*Session

// Listen ... Listens for new connections
func (s *GameServer) Listen() error {
	s.RSA = crypt.LoadRSA()

	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Println("ArcheAge Game server started at", s.Address)

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
			case 18:
				sess.Pong(reader)
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
			//reader := &packet.PacketReader{Pack: packBuf[4 : plen+4], Offset: 0}
			decr := sess.cr.Decrypt(reader.Pack[2:], len(reader.Pack))
			//seq := decr[0]  // seq?
			//hash := decr[1] // hash?
			opcode = binary.LittleEndian.Uint16(decr[2:4])
			//fmt.Printf("[%v] %v\n", sess.kostyl, hex.EncodeToString(decr))

			switch opcode {
			case 0x84:
				sess.OnMovement(decr)
				print("Movement")
			default:
				//sess.World_6_BigPacket()
				switch sess.kostyl {
				case 1:
					//sess.BeginGame()
					data, _ := hex.DecodeString("2400dd0564f1fc825223f4c495643405d55a754516e6a91e947cf7c797704010e0b081514272")
					sess.conn.Write(data)
					data, _ = hex.DecodeString("1d00dd05107771045f36774517e6bd86214285b4fe1f2e30d1bd8b5dc4f4231d00dd05cd7071045f36774514e6bd86214285b4fe1f2e30d2bd8b5dc4f423")
					sess.conn.Write(data)
					//print("1\n")
				case 3:
					data, _ := hex.DecodeString("0c00dd05f26537116a238351c6f7")
					sess.conn.Write(data)
					//print("3\n")
				case 4:
					data, _ := hex.DecodeString("1000dd05f631c9c797704010e0b0815186b7")
					sess.conn.Write(data)
					//print("4\n")
				case 6:
					sess.World_6_BigPacket()
					//print("6\n")
					//fmt.Println("[WORLD-ENCR] No opcode found:", opcode)
				}
				sess.kostyl++
			}
			fmt.Printf("[%v] %v\n", strconv.FormatInt(int64(sess.cr.Seq), 16), hex.EncodeToString(decr))
		default:
			fmt.Println("[GAME] No such subtype:", subtype)
		}
	}
}