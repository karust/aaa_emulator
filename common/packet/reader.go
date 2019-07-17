package packet

import (
	"encoding/binary"
	"net"
	"time"

	"../../common/crypt"
	"github.com/damnever/bitarray"
)

// Reader ... Reads packets
type Reader struct {
	Pack   []byte
	Offset uint16
	Len    uint16
	Err    bool
}

// CreateReader ... Creates packet reader
func CreateReader(buffer []byte) *Reader {
	return &Reader{Pack: buffer, Offset: 0, Len: uint16(len(buffer)), Err: false}
}

// GetPacketReader ... Reads length of packet, creates buffer and returns reader
// Read timeout is in seconds
func GetPacketReader(client net.Conn, timeout time.Duration) (reader *Reader, err error) {
	// Set read timout if user will hold connection too long
	if timeout != 0 {
		client.SetReadDeadline(time.Now().Add(timeout * time.Second))
	}
	// Read size of packet
	packLenBuf := make([]byte, 2)
	_, err = client.Read(packLenBuf)
	if err != nil {
		return
	}

	// Read packet
	packLen := binary.LittleEndian.Uint16(packLenBuf)
	packBuf := make([]byte, packLen)
	_, err = client.Read(packBuf)
	if err != nil {
		return
	}

	reader = CreateReader(packBuf)
	return
}

// GetEncPacketReader ... Get encrypted packet reader
func GetEncPacketReader(client net.Conn) (reader *Reader, err error) {
	// Read size of packet
	packLenBuf := make([]byte, 2)
	_, err = client.Read(packLenBuf)
	if err != nil {
		return
	}

	// Read packet
	packLen := binary.LittleEndian.Uint16(packLenBuf)
	packBuf := make([]byte, packLen)
	_, err = client.Read(packBuf)
	if err != nil {
		return
	}

	decr := crypt.ToClientEncr(packBuf[2:])
	reader = CreateReader(decr)
	return
}

// Byte .. read byte
func (pr *Reader) Byte() byte {
	if pr.Offset+1 > pr.Len {
		pr.Err = true
		return 0
	}
	defer func() { pr.Offset++ }()
	return byte(pr.Pack[pr.Offset])
}

// Short ... read short
func (pr *Reader) Short() uint16 {
	if pr.Offset+2 > pr.Len {
		pr.Err = true
		return 0
	}
	defer func() { pr.Offset += 2 }()
	return binary.LittleEndian.Uint16(pr.Pack[pr.Offset : pr.Offset+2])
}

// Int ... read integer
func (pr *Reader) Int() int {
	if pr.Offset+4 > pr.Len {
		pr.Err = true
		return 0
	}
	defer func() { pr.Offset += 4 }()
	return int(binary.LittleEndian.Uint32(pr.Pack[pr.Offset : pr.Offset+4]))
}

// UInt ... read integer
func (pr *Reader) UInt() uint32 {
	if pr.Offset+4 > pr.Len {
		pr.Err = true
		return 0
	}
	defer func() { pr.Offset += 4 }()
	return binary.LittleEndian.Uint32(pr.Pack[pr.Offset : pr.Offset+4])
}

// Float ... read integer
func (pr *Reader) Float() float32 {
	if pr.Offset+4 > pr.Len {
		pr.Err = true
		return 0
	}
	defer func() { pr.Offset += 4 }()
	return float32(binary.LittleEndian.Uint32(pr.Pack[pr.Offset : pr.Offset+4]))
}

// Int24 ... read integer24
func (pr *Reader) Int24() int {
	if pr.Offset+3 > pr.Len {
		pr.Err = true
		return 0
	}
	b1 := pr.Byte()
	b2 := pr.Byte()
	b3 := pr.Byte()
	num := int(b3)<<16 | int(b2)<<8 | int(b1)
	return num
}

// Long ... read long integer
func (pr *Reader) Long() uint64 {
	if pr.Offset+8 > pr.Len {
		pr.Err = true
		return 0
	}
	defer func() { pr.Offset += 8 }()
	return binary.LittleEndian.Uint64(pr.Pack[pr.Offset : pr.Offset+8])
}

// String ... read string of provided length
func (pr *Reader) String() string {
	len := pr.Short()
	if pr.Offset+len > pr.Len {
		pr.Err = true
		return ""
	}
	defer func() { pr.Offset += len }()
	return string(pr.Pack[pr.Offset : pr.Offset+len])
}

// BytesLen ... return bytes of required length
func (pr *Reader) BytesLen(len uint16) []byte {
	if pr.Offset+len > pr.Len {
		pr.Err = true
		return nil
	}
	defer func() { pr.Offset += len }()
	return pr.Pack[pr.Offset : pr.Offset+len]
}

// Bytes ... return bytes
func (pr *Reader) Bytes() []byte {
	len := pr.Short()
	return pr.BytesLen(len)
}

// Bool .. reads boolean
func (pr *Reader) Bool() bool {
	b := pr.Byte()
	if b == 0 {
		return false
	}
	return true
}

// Pisc ...
func (pr *Reader) Pisc(count int) []uint64 {
	result := make([]uint64, count)
	pish := bitarray.New(int(pr.Byte()))
	for index := 0; index < count*2; index += 2 {
		v, _ := pish.Get(index)
		v1, _ := pish.Get(index + 1)
		if v != 0 && v1 != 0 { // uint
			result[index/2] = uint64(pr.UInt())
		} else if v1 != 0 { // bc
			result[index/2] = uint64(pr.Int24())
		} else if v != 0 { // ushort
			result[index/2] = uint64(pr.Short())
		} else { // byte
			result[index/2] = uint64(pr.Byte())
		}
	}
	return result
}
