package gpio_test

import (
	. "../gpio"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {

	Describe("gpio", func() {

		It("should do something awesome", func() {
			p := Placeholder{}
			var err error
			Ω(err).ShouldNot(HaveOccurred())
			Ω(p).ShouldNot(BeNil())
		})
	})
})
