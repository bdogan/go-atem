package atem

func (a *Atem) PerformCut() {
	a.sendCommand(newCommand("DCut", []byte{uint8(a.MixEffectConfig.ME), 0, 0, 0}))
}

func (a *Atem) PerformAutoTransition() {
	a.sendCommand(newCommand("DAut", []byte{uint8(a.MixEffectConfig.ME), 0, 0, 0}))
}

func (a *Atem) SetPreviewInput(input VideoInputType, mixEffectIndex uint8) {
	// TODO: Check if input is supported by the ATEM
	index := uint16(input)
	// Check if the requested input is already on the preview bus
	// or if the requested M/E is outside of the supported range (ex. requesting M/E 2 on 1 M/E)
	if a.ProgramInput.index == index || mixEffectIndex > uint8(a.MixEffectConfig.ME) {
		return
	}
	a.sendCommand(newCommand("CPvI", []byte{uint8(a.MixEffectConfig.ME), 0, uint8(index >> 8), uint8(index & 0xFF)}))
}

func (a *Atem) SetProgramInput(input VideoInputType, mixEffectIndex uint8) {
	// TODO: Check if input is supported by the ATEM
	index := uint16(input)
	// Check if the requested input is already on the preview bus
	// or if the requested M/E is outside of the supported range (ex. requesting M/E 2 on 1 M/E)
	if a.ProgramInput.index == index || mixEffectIndex > uint8(a.MixEffectConfig.ME) {
		return
	}

	a.sendCommand(newCommand("CPgI", []byte{uint8(a.MixEffectConfig.ME), 0, uint8(index >> 8), uint8(index & 0xFF)}))
}
