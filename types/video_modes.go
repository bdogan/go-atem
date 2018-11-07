package types

import "fmt"

type ScanType string

const (
	ProgressiveScanType  	= "p"
	InterlaceScanType		= "i"
)

type VideoRegion string

const (
	UndefinedVideoRegion	= ""
	PALVideoRegion			= "PAL"
	NTSCVideoRegion			= "NTSC"
)

type AspectRatio string

const (
	UndefinedAspectRatio	= ""
	WideAscpectRatio		= "16:9"
	SquareAscpectRatio		= "4:3"
)

var VideoModes = []*VideoMode{
	NewVideoMode(0, 525, InterlaceScanType, 59.94, NTSCVideoRegion, SquareAscpectRatio),
	NewVideoMode(1, 625, InterlaceScanType, 50, PALVideoRegion, SquareAscpectRatio),
	NewVideoMode(2, 525, InterlaceScanType, 59.94, NTSCVideoRegion, WideAscpectRatio),
	NewVideoMode(3, 625, InterlaceScanType, 50, PALVideoRegion, WideAscpectRatio),
	NewVideoMode(4, 720, ProgressiveScanType, 50, UndefinedVideoRegion, UndefinedAspectRatio),
	NewVideoMode(5, 720, ProgressiveScanType, 59.94, UndefinedVideoRegion, UndefinedAspectRatio),
	NewVideoMode(6, 1080, InterlaceScanType, 50, UndefinedVideoRegion, UndefinedAspectRatio),
	NewVideoMode(7, 1080, InterlaceScanType, 59.94, UndefinedVideoRegion, UndefinedAspectRatio),
	NewVideoMode(8, 1080, ProgressiveScanType, 23.98, UndefinedVideoRegion, UndefinedAspectRatio),
	NewVideoMode(9, 1080, ProgressiveScanType, 24, UndefinedVideoRegion, UndefinedAspectRatio),
	NewVideoMode(10, 1080, ProgressiveScanType, 25, UndefinedVideoRegion, UndefinedAspectRatio),
	NewVideoMode(11, 1080, ProgressiveScanType, 29.97, UndefinedVideoRegion, UndefinedAspectRatio),
	NewVideoMode(12, 1080, ProgressiveScanType, 50, UndefinedVideoRegion, UndefinedAspectRatio),
	NewVideoMode(13, 1080, ProgressiveScanType, 59.94, UndefinedVideoRegion, UndefinedAspectRatio),
	NewVideoMode(14, 2160, ProgressiveScanType, 23.98, UndefinedVideoRegion, UndefinedAspectRatio),
	NewVideoMode(15, 2160, ProgressiveScanType, 24, UndefinedVideoRegion, UndefinedAspectRatio),
	NewVideoMode(16, 2160, ProgressiveScanType, 25, UndefinedVideoRegion, UndefinedAspectRatio),
	NewVideoMode(17, 2160, ProgressiveScanType, 29.97, UndefinedVideoRegion, UndefinedAspectRatio),
}

type VideoMode struct {
	Lines uint16
	ScanType ScanType
	FrameRate float32
	VideoRegion VideoRegion
	AspectRatio AspectRatio
	index uint16
}

func NewVideoMode(index uint16, lines uint16, scanType ScanType, frameRate float32, videoRegion VideoRegion, aspectRatio AspectRatio) *VideoMode {
	return &VideoMode{ index: index, Lines: lines, ScanType: scanType, FrameRate: frameRate, VideoRegion: videoRegion, AspectRatio: aspectRatio }
}

func (vm *VideoMode) IsSupported(vmode uint16) bool {
	return vm.index < vmode
}

func (vm *VideoMode) String() string {
	toString := fmt.Sprintf("%d%s%.2f", vm.Lines, vm.ScanType, vm.FrameRate)
	if vm.VideoRegion != UndefinedVideoRegion {
		toString += " " + string(vm.VideoRegion)
	}
	if vm.AspectRatio != UndefinedAspectRatio {
		toString += " " + string(vm.AspectRatio)
	}
	return toString
}