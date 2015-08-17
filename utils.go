package edison

import (
	"errors"
	"fmt"
	"net"
)

// DeviceID returns the current device ID. Currently this is the MAC address
// of the wifi interface on the device.
func DeviceID() (uint64, error) {
	// Grab hardware ID for the WiFi interface
	// we assume it is `wlan0`
	interfaces, err := net.Interfaces()
	if err != nil {
		return 0, err
	}
	for _, intf := range interfaces {
		if intf.Name == "wlan0" {
			return AddrToID(intf.HardwareAddr), nil
		}
	}

	return 0, errors.New("Missing WiFi interface wlan0")
}

// DeviceIDx returns the current device ID as a hex encoded string. Most
// APIs will expect the hex encoded stringversion of the device ID so this
// is a convenience method to give you that directly.
func DeviceIDx() (string, error) {
	id, err := DeviceID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", id), nil
}

// AddrToID converts a network hardware address to a device ID.
func AddrToID(addr net.HardwareAddr) uint64 {
	id := uint64(0)
	size := len(addr) - 1
	for pos, val := range addr {
		id |= uint64(uint64(val) << uint((size-pos)*8))
	}

	return id
}
