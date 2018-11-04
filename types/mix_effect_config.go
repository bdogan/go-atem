package types

type AtemMeModel uint8

const (
	ME1 = 0
	ME2 = 1
)

type MixEffectConfig struct {
	ME AtemMeModel
	KeyersOnME uint8
}
