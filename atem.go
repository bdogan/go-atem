package atem

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"time"
)

var (
	dstPort = 9910
	c       = func() {}
)

type AtemCallback func()

type Atem struct {

	// Public
	Ip    string
	State ConnectionState
	Debug bool
	UID   uint16

	// Atem
	ProtocolVersion  Version
	ProductId        NullTerminatedString
	Warn             NullTerminatedString
	Topology         Topology
	MixEffectConfig  MixEffectConfig
	MediaPlayers     MediaPlayers
	MultiViewCount   uint8
	AudioMixerConfig AudioMixerConfig
	VideoMixerConfig VideoMixerConfig
	MacroPool        uint8
	PowerStatus      PowerStatus
	VideoMode        *VideoMode
	VideoSources     *VideoSources
	ProgramInput     *VideoSource
	PreviewInput     *VideoSource

	// Private
	connection     net.Conn
	bodyBuffer     []byte
	commandBuffer  []*atemCommand
	outPacketQueue chan *atemPacket
	inPacketQueue  chan *atemPacket
	inCmdQueue     chan *atemCommand
	initialized    bool
	ackRequestID   uint16
	listeners      map[string][]AtemCallback
}

type ConnectionState int

const (
	Open       ConnectionState = 1
	Connecting ConnectionState = 2
	Closed     ConnectionState = 3
)

// Public Static Zone Start

func Create(Ip string, Debug bool) *Atem {
	atem := &Atem{Ip: Ip, Debug: Debug, State: Closed, listeners: map[string][]AtemCallback{}}

	// Initials
	atem.VideoSources = CreateVideoSourceList()

	return atem
}

// Public Zone Start

func (a *Atem) Connect() error {
	// Check already connected
	if a.State != Closed {
		return errors.New("already connected to server: " + a.Ip)
	}

	// Init
	a.UID = 0x53AB
	a.bodyBuffer = make([]byte, 0)
	a.commandBuffer = make([]*atemCommand, 0)
	a.initialized = false
	a.ackRequestID = 0

	// Trying to connect
	a.State = Connecting
	var err error
	a.connection, err = net.DialTimeout("udp", a.Ip+":"+strconv.Itoa(dstPort), time.Duration(time.Millisecond*1000))
	if err != nil {
		a.State = Closed
		return err
	}
	a.State = Open

	// Send hello packet
	a.writePacket(newConnectCmd(a.UID))

	// Read hi packet
	p, err := a.readPacket(time.Now().Add(time.Millisecond * 100))
	if err != nil || p == nil || !p.is(connectCommand) || p.body[0] != 0x2 {
		a.State = Closed
		if err != nil {
			return err
		}
		return errors.New("unable to connect device")
	}
	a.UID = p.uid

	// Send OK
	a.writePacket(newAckCmd(a.UID, 0))

	// Increase local request id
	a.ackRequestID++

	// Create chan
	a.inCmdQueue = make(chan *atemCommand)
	a.outPacketQueue = make(chan *atemPacket)
	a.inPacketQueue = make(chan *atemPacket)

	// Go queues
	go a.processInCmdQueue()
	go a.processOutPacketQueue()
	go a.processInPacketQueue()

	// Check read pipe
	a.processReadQueue()

	// Close connection
	if a.connection != nil {
		a.connection.Close()
	}

	// Return success
	return nil
}

func (a *Atem) Connected() bool {
	return a.State == Open && a.connection != nil
}

func (a *Atem) On(event string, callback func()) {
	if _, exists := a.listeners[event]; !exists {
		a.listeners[event] = make([]AtemCallback, 0)
	}
	a.listeners[event] = append(a.listeners[event], callback)
}

func (a *Atem) Close() {
	if a.Connected() {
		a.State = Closed
		a.connection.Close()
		a.connection = nil
	}
	if a.initialized {
		a.emit("closed")
		a.initialized = false
	}

}

// Private Zone Start

func (a *Atem) emit(event string, params ...interface{}) {
	if listeners, exists := a.listeners[event]; exists {
		in := make([]reflect.Value, len(params))
		for k, param := range params {
			in[k] = reflect.ValueOf(param)
		}
		for _, cb := range listeners {
			go reflect.ValueOf(cb).Call(in)
		}
	}
}

func (a *Atem) writePacket(p *atemPacket) error {
	if !a.Connected() {
		return errors.New("connection error on write packet")
	}
	if p.is(syncCommand) && len(a.commandBuffer) > 0 {
		// Append to body
		for _, cmd := range a.commandBuffer {
			p.appendCommand(cmd)
		}
		// Clear command buffer
		a.commandBuffer = make([]*atemCommand, 0)
	}
	_, err := a.connection.Write(p.toBytes())
	if err != nil {
		a.Close()
		return err
	}
	if a.Debug {
		fmt.Printf("Write: \t%x\n", p.toBytes())
	}
	return nil
}

func (a *Atem) readPacket(timeout time.Time) (*atemPacket, error) {
	if !a.Connected() {
		return nil, errors.New("connection error on read packet")
	}
	var packetBuffer [2060]byte
	a.connection.SetReadDeadline(timeout)
	n, err := a.connection.Read(packetBuffer[0:])
	if err != nil {
		return nil, err
	}
	p := parsePacket(packetBuffer[0:n])
	if a.Debug {
		fmt.Printf("Read: \t%x\n", p.toBytes()[0:12])
	}
	return p, nil
}

func (a *Atem) SendCommand(c *atemCommand) {
	a.commandBuffer = append(a.commandBuffer, c)
}

func (a *Atem) writePacketQueue(p *atemPacket) {
	// Send packet to queue
	a.outPacketQueue <- p
}

func (a *Atem) processInCmdQueue() {
	for a.Connected() {
		// Get command from queue
		c := <-a.inCmdQueue

		// Debug
		if a.Debug {
			fmt.Println(c)
		}

		// Save command
		switch c.Name {
		case "_ver":
			a.ProtocolVersion = Version{Major: binary.BigEndian.Uint16(c.Body[0:2]), Minor: binary.BigEndian.Uint16(c.Body[2:4])}
		case "_pin":
			a.ProductId = NullTerminatedString{Body: c.Body}
		case "Warn":
			a.Warn = NullTerminatedString{Body: c.Body}
		case "_top":
			a.Topology = Topology{MEs: c.Body[0], Sources: c.Body[1], ColorGenerators: c.Body[2], AUXBusses: c.Body[3], DownstreamKeyes: c.Body[4], Stringers: c.Body[5], DVEs: c.Body[6], SuperSources: c.Body[7], UnknownByte8: c.Body[8], HasSDOutput: (c.Body[9] & 1) == 1, UnknownByte10: c.Body[10]}
		case "_MeC":
			a.MixEffectConfig = MixEffectConfig{ME: AtemMeModel(c.Body[0]), KeyersOnME: c.Body[1]}
		case "_mpl":
			a.MediaPlayers = MediaPlayers{StillBanks: c.Body[0], ClipBanks: c.Body[1]}
		case "_MvC":
			a.MultiViewCount = c.Body[0]
		case "_AMC":
			a.AudioMixerConfig = AudioMixerConfig{AudioChannels: c.Body[0], HasMonitor: (c.Body[1] & 1) == 1}
		case "_VMC":
			a.VideoMixerConfig = NewVideoMixerConfig(binary.BigEndian.Uint16(c.Body[0:2]))
		case "_MAC":
			a.MacroPool = c.Body[0]
		case "Powr":
			a.PowerStatus = PowerStatus{MainPower: c.Body[0]&1 == 1, BackupPower: c.Body[0]&(1<<1) == (1 << 1)}
		case "VidM":
			a.VideoMode = NewVideoModeByIndex(c.Body[0])
		case "InPr":
			a.VideoSources.Update(c.Body)
		case "PrgI":
			a.ProgramInput = a.VideoSources.Get(binary.BigEndian.Uint16(c.Body[2:4]))
		case "PrvI":
			a.PreviewInput = a.VideoSources.Get(binary.BigEndian.Uint16(c.Body[2:4]))
		}

		// Trigger change command
		a.emit(c.Name + ".change")
	}
}

func (a *Atem) parseBodyBuffer() int {

	// Total command count
	parsedCommandTotal := 0

	// Check has bytes in buffer
	if len(a.bodyBuffer) == 0 {
		return parsedCommandTotal
	}

	// Read body start to end
	byteCursor := uint16(0)
	totalBytes := uint16(len(a.bodyBuffer))
	for totalBytes > byteCursor {
		packetLength := binary.BigEndian.Uint16(a.bodyBuffer[byteCursor : byteCursor+2])
		a.inCmdQueue <- parseCommand(a.bodyBuffer[byteCursor : byteCursor+packetLength])
		byteCursor = byteCursor + packetLength
		parsedCommandTotal++
	}

	// Clean body buffer
	a.bodyBuffer = make([]byte, 0)

	return parsedCommandTotal
}

func (a *Atem) processInPacketQueue() {
	for a.Connected() {
		// Get packet from queue
		p := <-a.inPacketQueue

		// Change uid given
		a.UID = p.uid

		// Inspect packet
		switch true {

		// Sync command
		case p.is(syncCommand):
			if p.hasBody() {
				// Append to body buffer
				a.bodyBuffer = append(a.bodyBuffer, p.body...)
			}
			if a.initialized || !p.hasBody() {
				// Parse body buffer
				a.parseBodyBuffer()

				// Send ack
				a.writePacketQueue(newAckCmd(a.UID, p.ackRequestID))
				a.writePacketQueue(newSyncCommand(a.UID, a.ackRequestID))
				a.ackRequestID++

				// Trigger connected
				if !a.initialized {
					a.initialized = true
					a.emit("connected")
				}
			}

		// Ack command
		case p.is(ackCommand):
			// To Do: Check from memory

		// Else is close
		default:
			if a.Debug {
				fmt.Printf("Unknown packet received:\t%xb\n", p.toBytes())
			}
			a.Close()
		}
	}
}

func (a *Atem) processOutPacketQueue() {
	for a.Connected() {
		p := <-a.outPacketQueue
		a.writePacket(p)
	}
}

func (a *Atem) processReadQueue() {
	for a.Connected() {
		p, err := a.readPacket(time.Now().Add(time.Millisecond * 1000))
		if err != nil {
			if a.Debug {
				fmt.Println("Connection closed on read")
				fmt.Println(err.Error())
			}
			a.Close()
			return
		}
		a.inPacketQueue <- p
	}
}
