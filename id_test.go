package edison_test

import (
	"net"

	. "github.com/gopackage/edison"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {

	Describe("AddrToID", func() {

		It("should convert standard 6 byte MAC addresses", func() {

			addr, err := net.ParseMAC("01:23:45:67:89:ab")
			Ω(err).ShouldNot(HaveOccurred())
			id := AddrToID(addr)
			Ω(id).Should(Equal(uint64(0x0123456789ab)))
		})

		It("should convert standard 8 byte MAC addresses", func() {

			addr, err := net.ParseMAC("01:23:45:67:89:ab:cd:ef")
			Ω(err).ShouldNot(HaveOccurred())
			id := AddrToID(addr)
			Ω(id).Should(Equal(uint64(0x0123456789abcdef)))
		})
	})
})
