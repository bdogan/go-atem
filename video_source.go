package atem

import (
	"encoding/binary"
	"fmt"
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
	0:   "External",
	1:   "Black",
	2:   "ColorBars",
	3:   "ColorGenerator",
	4:   "MediaPlayerFill",
	5:   "MediaPlayerKey",
	6:   "SuperSource",
	128: "MEOutput",
	129: "Auxilary",
	130: "Mask",
}

var Availability = map[uint8]string{
	0: "Auxilary",
	1: "Multiviewer",
	2: "SuperSourceArt",
	3: "SuperSourceBox",
	4: "KeySources",
}

var MEAvailability = map[uint8]string{
	0: "ME1 + FillSources",
	1: "ME2 + FillSources",
}

var VideoSourceType = map[uint16]string{
	0:     "Black",
	1:     "Input1",
	2:     "Input2",
	3:     "Input3",
	4:     "Input4",
	5:     "Input5",
	6:     "Input6",
	7:     "Input7",
	8:     "Input8",
	9:     "Input9",
	10:    "Input10",
	11:    "Input11",
	12:    "Input12",
	13:    "Input13",
	14:    "Input14",
	15:    "Input15",
	16:    "Input16",
	17:    "Input17",
	18:    "Input18",
	19:    "Input19",
	20:    "Input20",
	1000:  "ColorBars",
	2001:  "Color1",
	2002:  "Color2",
	3010:  "MediaPlayer1",
	3011:  "MediaPlayer1Key",
	3020:  "MediaPlayer2",
	3021:  "MediaPlayer2Key",
	4010:  "Key1Mask",
	4020:  "Key2Mask",
	4030:  "Key3Mask",
	4040:  "Key4Mask",
	5010:  "DSK1Mask",
	5020:  "DSK2Mask",
	6000:  "SuperSource",
	7001:  "CleanFeed1",
	7002:  "CleanFeed2",
	8001:  "Auxilary1",
	8002:  "Auxilary2",
	8003:  "Auxilary3",
	8004:  "Auxilary4",
	8005:  "Auxilary5",
	8006:  "Auxilary6",
	10010: "ME1Prog",
	10011: "ME1Prev",
	10020: "ME2Prog",
	10021: "ME2Prev",
}

type VideoSources struct {
	list *map[uint16]*VideoSource
}

type VideoInputType uint16

const (
	VideoBlack      VideoInputType = 0
	VideoInput1     VideoInputType = 1
	VideoInput2     VideoInputType = 2
	VideoInput3     VideoInputType = 3
	VideoInput4     VideoInputType = 4
	VideoInput5     VideoInputType = 5
	VideoInput6     VideoInputType = 6
	VideoInput7     VideoInputType = 7
	VideoInput8     VideoInputType = 8
	VideoInput9     VideoInputType = 9
	VideoInput10    VideoInputType = 10
	VideoInput11    VideoInputType = 11
	VideoInput12    VideoInputType = 12
	VideoInput13    VideoInputType = 13
	VideoInput14    VideoInputType = 14
	VideoInput15    VideoInputType = 15
	VideoInput16    VideoInputType = 16
	VideoInput17    VideoInputType = 17
	VideoInput18    VideoInputType = 18
	VideoInput19    VideoInputType = 19
	VideoInput20    VideoInputType = 20
	ColorBars       VideoInputType = 1000
	Color1          VideoInputType = 2001
	Color2          VideoInputType = 2002
	MediaPlayer1    VideoInputType = 3010
	MediaPlayer1Key VideoInputType = 3011
	MediaPlayer2    VideoInputType = 3020
	MediaPlayer2Key VideoInputType = 3021
	Key1Mask        VideoInputType = 4010
	Key2Mask        VideoInputType = 4020
	Key3Mask        VideoInputType = 4030
	Key4Mask        VideoInputType = 4040
	DSK1Mask        VideoInputType = 5010
	DSK2Mask        VideoInputType = 5020
	SuperSource     VideoInputType = 6000
	CleanFeed1      VideoInputType = 7001
	CleanFeed2      VideoInputType = 7002
	Auxilary1       VideoInputType = 8001
	Auxilary2       VideoInputType = 8002
	Auxilary3       VideoInputType = 8003
	Auxilary4       VideoInputType = 8004
	Auxilary5       VideoInputType = 8005
	Auxilary6       VideoInputType = 8006
	ME1Prog         VideoInputType = 10010
	ME1Prev         VideoInputType = 10011
	ME2Prog         VideoInputType = 10020
	ME2Prev         VideoInputType = 10021
)

func CreateVideoSourceList() *VideoSources {
	list := map[uint16]*VideoSource{}
	return &VideoSources{list: &list}
}

func (vss *VideoSources) Update(data []byte) {
	inputIndex := binary.BigEndian.Uint16(data[0:2])
	videoSource, exists := (*vss.list)[inputIndex]
	if !exists {
		videoSource = &VideoSource{Index: inputIndex}
		(*vss.list)[inputIndex] = videoSource
	}
	videoSource.Update(data)
}

func (vss *VideoSources) Get(index uint16) *VideoSource {
	videoSource, exists := (*vss.list)[index]
	if !exists {
		return nil
	}
	return videoSource
}

func (vss *VideoSources) String() string {
	var list []string
	for _, vs := range *vss.list {
		list = append(list, vs.String())
	}
	return strings.Join(list, "\n")
}

type VideoSource struct {
	Index                      uint16
	Type                       string
	LongName                   NullTerminatedString
	ShortName                  NullTerminatedString
	AvailableExternalPortTypes []string
	ExternalPortType           string
	PortType                   string
	Availability               []string
	MEAvailability             []string
}

func (vs *VideoSource) String() string {
	return fmt.Sprintf("[Type: %s, LongName: %s, ShortName: %s, AvailableExternalPortTypes: %s, ExternalPortType: %s, PortType: %s, Availability: %s, MEAvailibilty: %s]", vs.Type, vs.LongName.String(), vs.ShortName.String(), vs.AvailableExternalPortTypes, vs.ExternalPortType, vs.PortType, vs.Availability, vs.MEAvailability)
}

func (vs *VideoSource) Update(data []byte) {
	vs.Type = VideoSourceType[vs.Index]
	vs.LongName = NullTerminatedString{Body: data[2:22]}
	vs.ShortName = NullTerminatedString{Body: data[22:26]}
	vs.AvailableExternalPortTypes = []string{}

	// Available Ext Port Types
	for i, v := range VideoSourceAvailableExtPortTypes {
		if data[29]&(1<<i) == (1 << i) {
			vs.AvailableExternalPortTypes = append(vs.AvailableExternalPortTypes, v)
		}
	}

	// Ext Port Type
	vs.ExternalPortType = VideoSourceExtPortTypes[data[31]]

	// Port Type
	vs.PortType = VideoSourcePortTypes[data[32]]

	// Availability
	for i, v := range Availability {
		if data[34]&(1<<i) == (1 << i) {
			vs.Availability = append(vs.Availability, v)
		}
	}

	// MeAvailability
	for i, v := range MEAvailability {
		if data[35]&(1<<i) == (1 << i) {
			vs.MEAvailability = append(vs.MEAvailability, v)
		}
	}
}
