package video_source

type VideoSourceType uint16

const (
	Black			VideoSourceType = 0
	Input1			VideoSourceType = 1
	Input2			VideoSourceType = 2
	Input3			VideoSourceType = 3
	Input4			VideoSourceType = 4
	Input5			VideoSourceType = 5
	Input6			VideoSourceType = 6
	Input7			VideoSourceType = 7
	Input8			VideoSourceType = 8
	Input9			VideoSourceType = 9
	Input10			VideoSourceType = 10
	Input11			VideoSourceType = 11
	Input12			VideoSourceType = 12
	Input13			VideoSourceType = 13
	Input14			VideoSourceType = 14
	Input15			VideoSourceType = 15
	Input16			VideoSourceType = 16
	Input17			VideoSourceType = 17
	Input18			VideoSourceType = 18
	Input19			VideoSourceType = 19
	Input20			VideoSourceType = 20
	ColorBars		VideoSourceType = 1000
	Color1			VideoSourceType = 2001
	Color2			VideoSourceType = 2002
	MediaPlayer1	VideoSourceType = 3010
	MediaPlayer1Key	VideoSourceType = 3011
	MediaPlayer2	VideoSourceType = 3020
	MediaPlayer2Key	VideoSourceType = 3021
	Key1Mask		VideoSourceType = 4010
	Key2Mask		VideoSourceType = 4020
	Key3Mask		VideoSourceType = 4030
	Key4Mask		VideoSourceType = 4040
	DSK1Mask		VideoSourceType = 5010
	DSK2Mask		VideoSourceType = 5020
	SuperSource		VideoSourceType = 6000
	CleanFeed1		VideoSourceType = 7001
	CleanFeed2		VideoSourceType = 7002
	Auxilary1		VideoSourceType = 8001
	Auxilary2		VideoSourceType = 8002
	Auxilary3		VideoSourceType = 8003
	Auxilary4		VideoSourceType = 8004
	Auxilary5		VideoSourceType = 8005
	Auxilary6		VideoSourceType = 8006
	ME1Prog			VideoSourceType = 10010
	ME1Prev			VideoSourceType = 10011
	ME2Prog			VideoSourceType = 10020
	ME2Prev			VideoSourceType = 10021
)