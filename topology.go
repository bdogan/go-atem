package atem

type Topology struct {
	MEs             uint8
	Sources         uint8
	ColorGenerators uint8
	AUXBusses       uint8
	DownstreamKeyes uint8
	Stringers       uint8
	DVEs            uint8
	SuperSources    uint8
	UnknownByte8    uint8
	HasSDOutput     bool
	UnknownByte10   uint8
}
