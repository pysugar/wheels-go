package codec

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

// Define message types
const (
	MSG_TYPE_AUTH      = 1
	MSG_TYPE_TEXT      = 2
	MSG_TYPE_HEARTBEAT = 3
	MSG_TYPE_ACK       = 4
)

// Message struct
type Message struct {
	Version  byte
	Type     byte
	Length   uint16 // Payload length
	MsgID    uint32
	Payload  []byte
	Checksum uint32
}

// Encode the message into bytes
func (m *Message) Encode() ([]byte, error) {
	m.Length = uint16(len(m.Payload)) // Only payload length
	buf := make([]byte, 8+m.Length+4) // Header (8) + Payload + Checksum (4)

	// Set header fields
	buf[0] = m.Version
	buf[1] = m.Type
	binary.BigEndian.PutUint16(buf[2:4], m.Length)
	binary.BigEndian.PutUint32(buf[4:8], m.MsgID)

	// Set payload
	copy(buf[8:8+m.Length], m.Payload)

	// Calculate checksum for the entire message excluding the checksum field itself
	m.Checksum = crc32.ChecksumIEEE(buf[:8+m.Length])
	binary.BigEndian.PutUint32(buf[8+m.Length:], m.Checksum)

	return buf, nil
}

// Decode bytes into a Message struct
func Decode(data []byte) (*Message, error) {
	if len(data) < 12 { // 8 bytes header + 4 bytes checksum
		return nil, errors.New("data too short to decode header")
	}

	length := binary.BigEndian.Uint16(data[2:4])
	msg := &Message{
		Version:  data[0],
		Type:     data[1],
		Length:   length,
		MsgID:    binary.BigEndian.Uint32(data[4:8]),
		Payload:  data[8 : 8+length],
		Checksum: binary.BigEndian.Uint32(data[len(data)-4:]),
	}

	// Validate the total length
	expectedLength := 8 + msg.Length + 4 // Header + Payload + Checksum
	if uint16(len(data)) != expectedLength {
		return nil, errors.New("invalid message length")
	}

	// Verify checksum
	calculated := crc32.ChecksumIEEE(data[:8+msg.Length])
	if msg.Checksum != calculated {
		return nil, errors.New("checksum mismatch")
	}

	return msg, nil
}

//func CalculateChecksum(data []byte) uint16 {
//	var sum uint32
//	for _, b := range data {
//		sum += uint32(b)
//	}
//	return uint16(sum & 0xFFFF)
//}
