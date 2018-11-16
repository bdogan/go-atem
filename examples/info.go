/* Example: info.go
 * Connects to an ATEM and reports
* information about the device
*/

package main

import (
	"log"

	"github.com/bdogan/go-atem"
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
	// Change this IP address as needed
	ipAddress := "10.246.0.41"

	app := app{
		atemClient: atem.Create(ipAddress, false),
	}

	// Set connected handler
	app.atemClient.On("connected", app.onAtemConnected)

	// Set closed handler
	app.atemClient.On("closed", app.onAtemClosed)

	// Make connection
	log.Fatal(app.atemClient.Connect())
}
