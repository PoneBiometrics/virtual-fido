//go:build linux || windows

package virtual_fido

import (
	"github.com/bulwarkid/virtual-fido/fido_client"
	"github.com/bulwarkid/virtual-fido/usbip"
)

func startClient(client fido_client.FIDOClient) {
	usbDevice := usbip.NewUSBDevice()
	server := usbip.NewUSBIPServer(usbDevice)
	server.Start()
}
