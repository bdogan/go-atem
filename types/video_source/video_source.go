package video_source

type VideoSource uint16

const (
	Black			VideoSource = 0
	Input1			VideoSource = 1
	Input2			VideoSource = 2
	Input3			VideoSource = 3
	Input4			VideoSource = 4
	Input5			VideoSource = 5
	Input6			VideoSource = 6
	Input7			VideoSource = 7
	Input8			VideoSource = 8
	Input9			VideoSource = 9
	Input10			VideoSource = 10
	Input11			VideoSource = 11
	Input12			VideoSource = 12
	Input13			VideoSource = 13
	Input14			VideoSource = 14
	Input15			VideoSource = 15
	Input16			VideoSource = 16
	Input17			VideoSource = 17
	Input18			VideoSource = 18
	Input19			VideoSource = 19
	Input20			VideoSource = 20
	ColorBars		VideoSource = 1000
	Color1			VideoSource = 2001
	Color2			VideoSource = 2002
	MediaPlayer1	VideoSource = 3010
	MediaPlayer1Key	VideoSource = 3011
	MediaPlayer2	VideoSource = 3020
	MediaPlayer2Key	VideoSource = 3021
	Key1Mask		VideoSource = 4010
	Key2Mask		VideoSource = 4020
	Key3Mask		VideoSource = 4030
	Key4Mask		VideoSource = 4040
	DSK1Mask		VideoSource = 5010
	DSK2Mask		VideoSource = 5020
	SuperSource		VideoSource = 6000
	CleanFeed1		VideoSource = 7001
	CleanFeed2		VideoSource = 7002
	Auxilary1		VideoSource = 8001
	Auxilary2		VideoSource = 8002
	Auxilary3		VideoSource = 8003
	Auxilary4		VideoSource = 8004
	Auxilary5		VideoSource = 8005
	Auxilary6		VideoSource = 8006
	ME1Prog			VideoSource = 10010
	ME1Prev			VideoSource = 10011
	ME2Prog			VideoSource = 10020
	ME2Prev			VideoSource = 10021
)