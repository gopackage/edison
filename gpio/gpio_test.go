package gpio_test

import (
	. "github.com/gopackage/edison/gpio"
	"github.com/gopackage/sysfs/sysfstest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("gpio", func() {

	Describe("PWM", func() {
		It("should initialize properly", func() {
			recorder := &sysfstest.FileRecorder{}
			// We want need to supply a period responder
			recorder.Respond("/sys/class/pwm/pwmchip0/pwm1/period", &sysfstest.StaticResponder{Text: "555"})
			pins := NewPins(recorder)
			pwm, err := pins.PWM(1)
			Ω(err).Should(BeNil())
			Ω(pwm).ShouldNot(BeNil())
			Ω(pwm.Period()).Should(Equal(555))
			records := recorder.Records()
			Ω(records).Should(HaveLen(2))
		})
	})
	Describe("Digital Out", func() {

	})
	Describe("Digital In", func() {

	})
	Describe("FromScale", func() {
		It("should calculate the scale of an input", func() {
			scale := FromScale(128, 0, 255)
			Ω(scale).Should(BeNumerically("~", .5, .01))
		})
	})
	Describe("Limits", func() {
		It("should validate", func() {
			l := Limits{0, 255}
			l.Validate()
			Ω(l.Min).Should(Equal(0))
			Ω(l.Max).Should(Equal(255))
			l = Limits{255, 0}
			l.Validate()
			Ω(l.Min).Should(Equal(0))
			Ω(l.Max).Should(Equal(255))
		})
	})
})
