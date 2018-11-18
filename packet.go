package atem

const (
	syncCommand         uint16 = 1
	connectCommand      uint16 = 2
	connectRetryCommand uint16 = 6
	resendCommand       uint16 = 0x4
	requestNextAfter    uint16 = 0x8
	ackCommand          uint16 = 0x10
)

type atemPacket struct {
	flag          uint16
	uid           uint16
	ackResponseID uint16
	ackRequestID  uint16
	header        [4]byte
	body          []byte
}

func newSyncCommand(uid uint16, requestId uint16) *atemPacket {
	return &atemPacket{
		flag:          syncCommand,
		uid:           uid,
		ackResponseID: 0,
		ackRequestID:  requestId,
		header:        [4]byte{0, 0, 0, 0},
	}
}

func newConnectCmd(uid uint16) *atemPacket {
	return &atemPacket{
		flag:          connectCommand,
		uid:           uid,
		ackResponseID: 0,
		ackRequestID:  0,
		header:        [4]byte{0, 0, 0, 0x03},
		body:          []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
}

func newAckCmd(uid uint16, ackResponseID uint16) *atemPacket {
	return &atemPacket{
		flag:          ackCommand,
		uid:           uid,
		ackResponseID: ackResponseID,
		ackRequestID:  0,
		header:        [4]byte{0, 0, 0, 0},
		body:          make([]byte, 0),
	}
}

func parsePacket(msg []byte) *atemPacket {
	return &atemPacket{
		flag:          uint16(msg[0] >> 3),
		uid:           uint16((uint16(msg[2]) << 8) | uint16(msg[3])),
		ackResponseID: uint16((uint16(msg[4]) << 8) | uint16(msg[5])),
		ackRequestID:  uint16((uint16(msg[10]) << 8) | uint16(msg[11])),
		header:        [4]byte{msg[6], msg[7], msg[8], msg[9]},
		body:          msg[12:]}
}

func (ap *atemPacket) appendCommand(cmd *AtemCommand) {
	// Add sync flag if not
	if !ap.is(syncCommand) {
		ap.addFlag(syncCommand)
	}
	// Append to body
	ap.body = append(ap.body, (*cmd).toBytes()...)
}

func packetFromCommand(cmd *AtemCommand, uid uint16, requestID uint16) *atemPacket {
	return &atemPacket{
		flag:          syncCommand,
		uid:           uid,
		ackResponseID: 0,
		ackRequestID:  requestID,
		header:        [4]byte{0, 0, 0, 0},
		body:          cmd.toBytes(),
	}
}

func (ap *atemPacket) addFlag(flag uint16) {
	ap.flag = ap.flag | flag
}

func (ap *atemPacket) is(cmd uint16) bool {
	return (ap.flag & cmd) == cmd
}

func (ap *atemPacket) length() uint16 {
	return uint16(12 + len(ap.body))
}

func (ap *atemPacket) hasBody() bool {
	return len(ap.body) > 0
}

func (ap *atemPacket) toBytes() []byte {
	var result []byte

	// Set flag & length
	result = append(result, []byte{uint8((ap.flag << 3) | ((ap.length() >> 8) & 0x7)), uint8(ap.length() & 0xFF)}...)

	// Set uid
	result = append(result, []byte{uint8(ap.uid >> 8), uint8(ap.uid & 0xFF)}...)

	// Set ackid
	result = append(result, []byte{uint8(ap.ackResponseID >> 8), uint8(ap.ackResponseID & 0xFF)}...)

	// Set zeros
	result = append(result, ap.header[0:4]...)

	// Set targetId
	result = append(result, []byte{uint8(ap.ackRequestID >> 8), uint8(ap.ackRequestID & 0xFF)}...)

	// Add body
	result = append(result, ap.body...)

	return result
}
