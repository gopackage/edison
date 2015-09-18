package ble_test

import (
	. "github.com/gopackage/edison/ble"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {

	Describe("Bluetooth", func() {

		It("should do something awesome", func() {
			b := Bluetooth{}
			var err error
			Ω(err).ShouldNot(HaveOccurred())
			Ω(b).ShouldNot(BeNil())
		})
	})
})
