package types

import "fmt"

type ScanType string

const (
	ProgressiveScanType  	= "p"
	InterlaceScanType		= "i"
)

type VideoRegion string

const (
	PALVideoRegion		= "PAL"
	NTSCVideoRegion		= "NTSC"
)

type VideoMode struct {
	Lines uint16
	ScanType ScanType
	FrameRate float32
	VideoRegion VideoRegion
}

func (vm *VideoMode) String() string {
	if vm.Lines != 0 {
		return fmt.Sprintf("%d%s%f", vm.Lines, vm.ScanType, vm.FrameRate)
	}
	return string(vm.VideoRegion)
}
