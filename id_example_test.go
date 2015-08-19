package edison_test

import (
	"fmt"
	"net"

	. "../edison"
)

func ExampleAddrToID() {
	addr, err := net.ParseMAC("01:23:45:67:89:ab")
	if err != nil {
		// fail
	}
	id := AddrToID(addr)
	fmt.Printf("%x", id)
	// Output: 123456789ab
}
