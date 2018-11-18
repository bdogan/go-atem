/* Example: info.go
 * Connects to an ATEM and reports
* information about the device
*/

package main

import (
	"flag"
	"log"
	"time"

	"github.com/bdogan/go-atem"
)

var (
	ipAddress = flag.String("ip", "", "Atem switcher ipv4 address")
	debug     = flag.Bool("debug", false, "Connection debugging")
)

type app struct {
	atemClient *atem.Atem
}

func (at *app) onAtemConnected() {
	log.Printf("ATEM connected at %s. UID:%d\n", at.atemClient.Ip, at.atemClient.UID)

	at.atemClient.SetProgramInput(atem.VideoBlack, 0)
	at.atemClient.SetPreviewInput(atem.VideoBlack, 0)

	setPgm := true
	input := atem.VideoInput1

	for at.atemClient.Connected() {
		time.Sleep(time.Millisecond * 100)
		if setPgm {
			at.atemClient.SetProgramInput(input, 0)
		} else {
			at.atemClient.SetPreviewInput(input, 0)
		}

		input++

		if input == atem.VideoInput9 {
			input = atem.VideoInput1
			setPgm = !setPgm
		}
	}
}

func (at *app) onAtemClosed() {
	log.Println("Connection closed")
}

func main() {
	// Parse flag
	flag.Parse()

	// Create app
	app := app{
		atemClient: atem.Create(*ipAddress, *debug),
	}

	// Set connected handler
	app.atemClient.On("connected", app.onAtemConnected)

	// Set closed handler
	app.atemClient.On("closed", app.onAtemClosed)

	// Make connection
	app.atemClient.Connect()
}
