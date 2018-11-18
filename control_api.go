package atem

func (a *Atem) PerformCut() {
	a.sendCommand(newCommand("DCut", []byte{uint8(a.MixEffectConfig.ME), 0, 0, 0}))
}

func (a *Atem) SetPreviewInput(input VideoInputType) {
	// TODO: Check if input is supported by the ATEM
	index := uint16(input)
	if a.PreviewInput.index == index {
		return
	}
	a.sendCommand(newCommand("CPvI", []byte{uint8(a.MixEffectConfig.ME), 0, uint8(index >> 8), uint8(index & 0xFF)}))
}

func (a *Atem) SetProgramInput(input VideoInputType) {
	// TODO: Check if input is supported by the ATEM
	index := uint16(input)
	if a.PreviewInput.index == index {
		return
	}
	a.sendCommand(newCommand("CPgI", []byte{uint8(a.MixEffectConfig.ME), 0, uint8(index >> 8), uint8(index & 0xFF)}))
}
