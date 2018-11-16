package atem

type VideoMixerConfig struct {
	SupportedVideoModes []*VideoMode
}

func NewVideoMixerConfig(configModes uint16) VideoMixerConfig {
	return VideoMixerConfig{SupportedVideoModes: VideoModes[0:configModes]}
}
