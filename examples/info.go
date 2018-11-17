/* Example: info.go
 * Connects to an ATEM and reports
* information about the device
*/

package main

import (
	"flag"
	"github.com/bdogan/go-atem"
	"log"
)

var (
	ipAddress 	= flag.String("ip", "", "Atem switcher ipv4 address")
	debug		= flag.Bool("debug", false, "Connection debugging")
)

type app struct {
	atemClient *atem.Atem
}

func (at *app) onAtemConnected() {
	log.Printf("ATEM connected at %s. UID:%d\n", at.atemClient.Ip, at.atemClient.UID)
	log.Printf("Product ID: %s, Protocol Version: %s\n", at.atemClient.ProductId.String(), at.atemClient.ProtocolVersion.String())
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
