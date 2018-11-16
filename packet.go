package atem

const (
	SyncCommand         uint16 = 1
	ConnectCommand      uint16 = 2
	ConnectRetryCommand uint16 = 6
	ResendCommand       uint16 = 0x4
	RequestNextAfter    uint16 = 0x8
	AckCommand          uint16 = 0x10
)

type AtemPacket struct {
	Flag          uint16
	UID           uint16
	AckResponseID uint16
	AckRequestID  uint16
	Header        [4]byte
	Body          []byte
}

func NewPacket(Flag uint16, UID uint16, AckResponseID uint16, AckRequestID uint16, Body []byte) *AtemPacket {
	return &AtemPacket{Flag: Flag, UID: UID, AckResponseID: AckResponseID, AckRequestID: AckRequestID, Header: [4]byte{0, 0, 0, 0}, Body: Body}
}

func NewConnectCmd(UID uint16) *AtemPacket {
	return &AtemPacket{Flag: ConnectCommand, UID: UID, AckResponseID: 0, AckRequestID: 0, Header: [4]byte{0, 0, 0, 0x03}, Body: []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}
}

func NewAckCmd(UID uint16, AckResponseID uint16) *AtemPacket {
	return &AtemPacket{Flag: AckCommand, UID: UID, AckResponseID: AckResponseID, AckRequestID: 0, Header: [4]byte{0, 0, 0, 0}, Body: make([]byte, 0)}
}

func ParsePacket(msg []byte) *AtemPacket {
	return &AtemPacket{
		Flag:          uint16(msg[0] >> 3),
		UID:           uint16((uint16(msg[2]) << 8) | uint16(msg[3])),
		AckResponseID: uint16((uint16(msg[4]) << 8) | uint16(msg[5])),
		AckRequestID:  uint16((uint16(msg[10]) << 8) | uint16(msg[11])),
		Header:        [4]byte{msg[6], msg[7], msg[8], msg[9]},
		Body:          msg[12:]}
}

func (ap *AtemPacket) Is(cmd uint16) bool {
	return (ap.Flag & cmd) == 1
}

func (ap *AtemPacket) Length() uint16 {
	return uint16(12 + len(ap.Body))
}

func (ap *AtemPacket) HasBody() bool {
	return len(ap.Body) > 0
}

func (ap *AtemPacket) ToBytes() []byte {
	var result []byte

	// Set flag & length
	result = append(result, []byte{uint8((ap.Flag << 3) | ((ap.Length() >> 8) & 0x7)), uint8(ap.Length() & 0xFF)}...)

	// Set uid
	result = append(result, []byte{uint8(ap.UID >> 8), uint8(ap.UID & 0xFF)}...)

	// Set ackid
	result = append(result, []byte{uint8(ap.AckResponseID >> 8), uint8(ap.AckResponseID & 0xFF)}...)

	// Set zeros
	result = append(result, ap.Header[0:4]...)

	// Set targetId
	result = append(result, []byte{uint8(ap.AckRequestID >> 8), uint8(ap.AckRequestID & 0xFF)}...)

	// Add body
	result = append(result, ap.Body...)

	return result
}
