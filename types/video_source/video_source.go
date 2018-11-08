package video_source

import (
	"encoding/binary"
	"fmt"
	"go-atem/types"
	"strings"
)

var VideoSourceAvailableExtPortTypes = map[uint8]string{
	0: "SDI",
	1: "HDMI",
	2: "Component",
	3: "Composite",
	4: "SVideo",
}

var VideoSourceExtPortTypes = map[uint8]string{
	0: "Internal",
	1: "SDI",
	2: "HDMI",
	3: "Composite",
	4: "Component",
	5: "SVideo",
}

var VideoSourcePortTypes = map[uint8]string{
	0: "External",
	1: "Black",
	2: "ColorBars",
	3: "ColorGenerator",
	4: "MediaPlayerFill",
	5: "MediaPlayerKey",
	6: "SuperSource",
	128: "MEOutput",
	129: "Auxilary",
	130: "Mask",
}

type Availability struct {
	Auxilary bool
	Multiviewer bool
	SuperSourceArt bool
	SuperSourceBox bool
	KeySources bool
}

func (a *Availability) String() string {
	return fmt.Sprintf("Auxilary: %t, Multiviewer: %t, SuperSourceArt: %t, SuperSourceBox: %t, KeySources: %t", a.Auxilary, a.Multiviewer, a.SuperSourceArt, a.SuperSourceBox, a.KeySources)
}

type MEAvailability struct {
	ME1FillSources bool
	ME2FillSources bool
}

func (mea *MEAvailability) String() string {
	return fmt.Sprintf("ME1FillSources: %t, ME2FillSources: %t", mea.ME1FillSources, mea.ME2FillSources)
}

var VideoSourceType = map[uint16]string{
	0		: "Black",
	1		: "Input1",
	2		: "Input2",
	3		: "Input3",
	4		: "Input4",
	5		: "Input5",
	6		: "Input6",
	7		: "Input7",
	8		: "Input8",
	9		: "Input9",
	10		: "Input10",
	11		: "Input11",
	12		: "Input12",
	13		: "Input13",
	14		: "Input14",
	15		: "Input15",
	16		: "Input16",
	17		: "Input17",
	18		: "Input18",
	19		: "Input19",
	20		: "Input20",
	1000	: "ColorBars",
	2001	: "Color1",
	2002	: "Color2",
	3010	: "MediaPlayer1",
	3011	: "MediaPlayer1Key",
	3020	: "MediaPlayer2",
	3021	: "MediaPlayer2Key",
	4010	: "Key1Mask",
	4020	: "Key2Mask",
	4030	: "Key3Mask",
	4040	: "Key4Mask",
	5010	: "DSK1Mask",
	5020	: "DSK2Mask",
	6000	: "SuperSource",
	7001	: "CleanFeed1",
	7002	: "CleanFeed2",
	8001	: "Auxilary1",
	8002	: "Auxilary2",
	8003	: "Auxilary3",
	8004	: "Auxilary4",
	8005	: "Auxilary5",
	8006	: "Auxilary6",
	10010	: "ME1Prog",
	10011	: "ME1Prev",
	10020	: "ME2Prog",
	10021	: "ME2Prev",
}

type VideoSources struct {
	list *map[uint16]*VideoSource
}

func CreateVideoSourceList() *VideoSources {
	list := map[uint16]*VideoSource{}
	return &VideoSources{ list: &list }
}

func (vss *VideoSources) Update(data []byte) {
	inputIndex := binary.BigEndian.Uint16(data[0:2])
	videoSource, exists := (*vss.list)[inputIndex]
	if !exists {
		videoSource = &VideoSource{ index: inputIndex }
		(*vss.list)[inputIndex] = videoSource
	}
	videoSource.Update(data)
}

func (vss *VideoSources) String() string {
	var list []string
	for _, vs := range *vss.list {
		list = append(list, vs.String())
	}
	return strings.Join(list, "\n")
}

type VideoSource struct {
	index uint16

	Type string
	LongName types.NullTerminatedString
	ShortName types.NullTerminatedString
	AvailableExternalPortTypes []string
	ExternalPortType string
	PortType string
	Availability Availability
	MEAvailability MEAvailability
}

func (vs *VideoSource) String() string {
	return fmt.Sprintf("[Type: %s, LongName: %s, ShortName: %s, AvailableExternalPortTypes: %s, ExternalPortType: %s, PortType: %s, Availability: %s, MEAvailibilty: %s]", vs.Type, vs.LongName.String(), vs.ShortName.String(), vs.AvailableExternalPortTypes, vs.ExternalPortType, vs.PortType, vs.Availability.String(), vs.MEAvailability.String())
}

func (vs *VideoSource) Update(data []byte) {
	vs.Type = VideoSourceType[vs.index]
	vs.LongName = types.NullTerminatedString{ Body: data[2:22] }
	vs.ShortName = types.NullTerminatedString{ Body: data[22:26] }
	vs.AvailableExternalPortTypes = []string{}

	// Available Ext Port Types
	for i, v := range VideoSourceAvailableExtPortTypes {
		if data[27] & (1 << i) == (1 << i) {
			vs.AvailableExternalPortTypes = append(vs.AvailableExternalPortTypes, v)
		}
	}

	// Ext Port Type
	vs.ExternalPortType = VideoSourceExtPortTypes[data[29]]

	// Port Type
	vs.PortType = VideoSourcePortTypes[data[30]]
}