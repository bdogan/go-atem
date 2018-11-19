package atem

func (a *Atem) PerformCut() {
	a.SendCommand(NewCommand("DCut", []byte{uint8(a.MixEffectConfig.ME), 0, 0, 0}))
}

func (a *Atem) PerformAutoTransition() {
	a.SendCommand(NewCommand("DAut", []byte{uint8(a.MixEffectConfig.ME), 0, 0, 0}))
}

func (a *Atem) SetPreviewInput(input VideoInputType, meIndex uint8) {
	// TODO: Check if input is supported by the ATEM
	index := uint16(input)
	// Check if the requested input is already on the preview bus
	// or if the requested M/E is outside of the supported range (ex. requesting M/E 2 on 1 M/E)
	if a.PreviewInput.Index == index || meIndex > uint8(a.MixEffectConfig.ME) {
		return
	}
	a.SendCommand(NewCommand("CPvI", []byte{meIndex, 0, uint8(index >> 8), uint8(index & 0xFF)}))
}

func (a *Atem) SetProgramInput(input VideoInputType, meIndex uint8) {
	// TODO: Check if input is supported by the ATEM
	index := uint16(input)
	// Check if the requested input is already on the preview bus
	// or if the requested M/E is outside of the supported range (ex. requesting M/E 2 on 1 M/E)
	if a.ProgramInput.Index == index || meIndex > uint8(a.MixEffectConfig.ME) {
		return
	}

	a.SendCommand(NewCommand("CPgI", []byte{meIndex, 0, uint8(index >> 8), uint8(index & 0xFF)}))
}

func (a *Atem) SetKeyerOnAir(enabled bool, keyerIndex uint8, meIndex uint8) {
	var enabledByte uint8

	if enabled {
		enabledByte = 1
	}

	body := []byte{meIndex, keyerIndex, 0, enabledByte}
	a.SendCommand(NewCommand("CKOn", body))
}

type MacroAction uint8

const (
	MacroRun               MacroAction = 0
	MacroStop              MacroAction = 1
	MacroStopRecording     MacroAction = 2
	MacroInsertWaitForUser MacroAction = 3
	MacroContinue          MacroAction = 4
	MacroDelete            MacroAction = 5
)

func (a *Atem) RunMacro(macroIndex uint8) {
	if macroIndex > a.MacroPool {
		return
	}
	body := []byte{0, macroIndex, uint8(MacroRun)}
	a.SendCommand(NewCommand("MAct", body))
}

func (a *Atem) StopMacro() {
	body := []byte{0xFF, 0xFF, uint8(MacroStop)}
	a.SendCommand(NewCommand("MAct", body))
}
