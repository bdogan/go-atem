package atem

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"time"

	"github.com/bdogan/go-atem/cmd"
	"github.com/bdogan/go-atem/packet"
	"github.com/bdogan/go-atem/types"
	"github.com/bdogan/go-atem/types/video_source"
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
	ProtocolVersion  types.Version
	ProductId        types.NullTerminatedString
	Warn             types.NullTerminatedString
	Topology         types.Topology
	MixEffectConfig  types.MixEffectConfig
	MediaPlayers     types.MediaPlayers
	MultiViewCount   uint8
	AudioMixerConfig types.AudioMixerConfig
	VideoMixerConfig types.VideoMixerConfig
	MacroPool        uint8
	PowerStatus      types.PowerStatus
	VideoMode        *types.VideoMode
	VideoSources     *video_source.VideoSources
	ProgramInput     *video_source.VideoSource
	PreviewInput     *video_source.VideoSource

	// Private
	connection     net.Conn
	bodyBuffer     []byte
	outPacketQueue chan *packet.AtemPacket
	inPacketQueue  chan *packet.AtemPacket
	inBodyQueue    chan []byte
	inCmdQueue     chan *cmd.AtemCmd
	initialized    bool
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
	atem.VideoSources = video_source.CreateVideoSourceList()

	return atem
}

// Public Zone Start

func (a *Atem) Connected() bool {
	return a.State == Open && a.connection != nil
}

func (a *Atem) On(event string, callback func()) {
	if _, exists := a.listeners[event]; !exists {
		a.listeners[event] = make([]AtemCallback, 0)
	}
	a.listeners[event] = append(a.listeners[event], callback)
}

func (a *Atem) Connect() error {
	for {
		err := a.connect()
		if a.Debug {
			fmt.Println(err)
		}
		if err != nil {
			return err
		}
	}
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
			reflect.ValueOf(cb).Call(in)
		}
	}
}

func (a *Atem) connect() error {
	// Check already connected
	if a.State != Closed {
		return errors.New("already connected to server: " + a.Ip)
	}

	// Init
	a.UID = 0x53AB
	a.bodyBuffer = make([]byte, 0)
	a.initialized = false

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
	a.writePacket(packet.NewConnectCmd(a.UID))

	// Read hi packet
	p, err := a.readPacket(time.Now().Add(time.Millisecond * 100))
	if err != nil || p == nil || p.Is(packet.ConnectCommand) || p.Body[0] != 0x2 {
		a.State = Closed
		if err != nil {
			return err
		}
		return errors.New("unable to connect device")
	}
	a.UID = p.UID

	// Send OK
	a.writePacket(packet.NewAckCmd(a.UID, 0))

	// Create chan
	a.inBodyQueue = make(chan []byte)
	a.inCmdQueue = make(chan *cmd.AtemCmd)
	a.outPacketQueue = make(chan *packet.AtemPacket)
	a.inPacketQueue = make(chan *packet.AtemPacket)

	// Go queues
	go a.processInCmdQueue()
	go a.processInBodyQueue()
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

func (a *Atem) writePacket(p *packet.AtemPacket) error {
	if !a.Connected() {
		return errors.New("connection error on write packet")
	}
	_, err := a.connection.Write(p.ToBytes())
	if err != nil {
		a.Close()
		return err
	}
	if a.Debug {
		fmt.Printf("Send: \t\t%x\n", p.ToBytes()[0:12])
	}
	return nil
}

func (a *Atem) readPacket(timeout time.Time) (*packet.AtemPacket, error) {
	if !a.Connected() {
		return nil, errors.New("connection error on read packet")
	}
	var packetBuffer [2060]byte
	a.connection.SetReadDeadline(timeout)
	n, err := a.connection.Read(packetBuffer[0:])
	if err != nil {
		return nil, err
	}
	p := packet.Parse(packetBuffer[0:n])
	if a.Debug {
		fmt.Printf("Receive: \t%x\n", p.ToBytes()[0:12])
	}
	return p, nil
}

func (a *Atem) writePacketQueue(p *packet.AtemPacket) {
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
			a.ProtocolVersion = types.Version{Major: binary.BigEndian.Uint16(c.Body[0:2]), Minor: binary.BigEndian.Uint16(c.Body[2:4])}
		case "_pin":
			a.ProductId = types.NullTerminatedString{Body: c.Body}
		case "Warn":
			a.Warn = types.NullTerminatedString{Body: c.Body}
		case "_top":
			a.Topology = types.Topology{MEs: c.Body[0], Sources: c.Body[1], ColorGenerators: c.Body[2], AUXBusses: c.Body[3], DownstreamKeyes: c.Body[4], Stringers: c.Body[5], DVEs: c.Body[6], SuperSources: c.Body[7], UnknownByte8: c.Body[8], HasSDOutput: (c.Body[9] & 1) == 1, UnknownByte10: c.Body[10]}
		case "_MeC":
			a.MixEffectConfig = types.MixEffectConfig{ME: types.AtemMeModel(c.Body[0]), KeyersOnME: c.Body[1]}
		case "_mpl":
			a.MediaPlayers = types.MediaPlayers{StillBanks: c.Body[0], ClipBanks: c.Body[1]}
		case "_MvC":
			a.MultiViewCount = c.Body[0]
		case "_AMC":
			a.AudioMixerConfig = types.AudioMixerConfig{AudioChannels: c.Body[0], HasMonitor: (c.Body[1] & 1) == 1}
		case "_VMC":
			a.VideoMixerConfig = types.NewVideoMixerConfig(binary.BigEndian.Uint16(c.Body[0:2]))
		case "_MAC":
			a.MacroPool = c.Body[0]
		case "Powr":
			a.PowerStatus = types.PowerStatus{MainPower: c.Body[0]&1 == 1, BackupPower: c.Body[0]&(1<<1) == (1 << 1)}
		case "VidM":
			a.VideoMode = types.NewVideoModeByIndex(c.Body[0])
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

func (a *Atem) processInBodyQueue() {
	for a.Connected() {
		// Get []byte from queue
		b := <-a.inBodyQueue

		// Check size
		if len(b) == 0 {
			continue
		}

		// Read body buffer
		byteCursor := uint16(0)
		totalBytes := uint16(len(b))
		for totalBytes > byteCursor {
			packetLength := binary.BigEndian.Uint16(b[byteCursor : byteCursor+2])
			a.inCmdQueue <- cmd.Parse(b[byteCursor : byteCursor+packetLength])
			byteCursor = byteCursor + packetLength
		}

		// Trigger connected
		if !a.initialized {
			a.initialized = true
			a.emit("connected")
		}

	}
}

func (a *Atem) processInPacketQueue() {
	for a.Connected() {
		// Get packet from queue
		p := <-a.inPacketQueue

		// Change uid given
		a.UID = p.UID

		// Inspect packet
		switch true {

		// Sync command
		case p.Is(packet.SyncCommand):
			if p.HasBody() {
				// Append to body buffer
				a.bodyBuffer = append(a.bodyBuffer, p.Body...)
			} else {
				a.inBodyQueue <- a.bodyBuffer
				// Clean body buffer
				a.bodyBuffer = make([]byte, 0)
				// Send ack
				a.writePacketQueue(packet.NewAckCmd(a.UID, p.AckRequestID))
			}

		// Else is close
		default:
			if a.Debug {
				fmt.Printf("Unknown packet received:\t%xb\n", p.ToBytes())
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
			}
			a.Close()
			return
		}
		a.inPacketQueue <- p
	}
}
