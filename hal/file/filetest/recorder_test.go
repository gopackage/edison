package filetest_test

import (
	"os"

	. "../../file"
	. "../filetest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {
	Describe("FileRecorder", func() {
		It("should implement DeviceFile", func() {
			var f DeviceFile
			f = &FileRecorder{}
			Ω(f).ShouldNot(BeNil())
		})

		Describe("StaticResponder", func() {
			responder := StaticResponder{Text: "Hello World"}

			It("should respond to writes without error", func() {
				Ω(responder.Write("foo")).Should(BeNil())
			})

			It("should respond to reads with static text", func() {
				Ω(responder.Read()).Should(Equal("Hello World"))
			})
		})

		Describe("ListResponder", func() {
			responder := ListResponder{}
			responder.Add("A")
			responder.Add("B")
			responder.Add("C")

			It("should respond to writes without error", func() {
				Ω(responder.Write("foo")).Should(BeNil())
			})

			It("should respond to reads with texts in order", func() {
				text, err := responder.Read()
				Ω(err).Should(BeNil())
				Ω(text).Should(Equal("A"))
				text, err = responder.Read()
				Ω(err).Should(BeNil())
				Ω(text).Should(Equal("B"))
				text, err = responder.Read()
				Ω(err).Should(BeNil())
				Ω(text).Should(Equal("C"))
				text, err = responder.Read()
				Ω(err).ShouldNot(BeNil())
				Ω(text).Should(Equal(""))
			})
		})

		Describe("recorder Write", func() {
			It("should respond to a write with response", func() {
				recorder := FileRecorder{}
				recorder.Respond("john", &StaticResponder{})
				err := recorder.Write("john", "smith")
				Ω(err).Should(BeNil())
				records := recorder.Records()

				// Verify records
				Ω(records).Should(HaveLen(1))
				record := records[0]
				Ω(record.Stamp).ShouldNot(BeNil())
				Ω(record.Err).Should(BeNil())
				Ω(record.Path).Should(Equal("/john"))
				Ω(record.Text).Should(Equal("smith"))
				Ω(record.Write).Should(BeTrue())
			})

			It("should respond to a write with no response with an error", func() {
				recorder := FileRecorder{}
				err := recorder.Write("foo", "bar")
				Ω(err).Should(Equal(os.ErrNotExist))
				recorder.Respond("john", &StaticResponder{})
				err = recorder.Write("foo", "bar")
				Ω(err).Should(Equal(os.ErrNotExist))
				err = recorder.Write("john", "smith")
				Ω(err).Should(BeNil())
				records := recorder.Records()

				// Verify records
				Ω(records).Should(HaveLen(3))
				record := records[0]
				Ω(record.Stamp).ShouldNot(BeNil())
				Ω(record.Err).Should(Equal(os.ErrNotExist))
				Ω(record.Path).Should(Equal("/foo"))
				Ω(record.Text).Should(Equal("bar"))
				Ω(record.Write).Should(BeTrue())
				record = records[1]
				Ω(record.Stamp).ShouldNot(BeNil())
				Ω(record.Err).Should(Equal(os.ErrNotExist))
				Ω(record.Path).Should(Equal("/foo"))
				Ω(record.Text).Should(Equal("bar"))
				Ω(record.Write).Should(BeTrue())
				record = records[2]
				Ω(record.Stamp).ShouldNot(BeNil())
				Ω(record.Err).Should(BeNil())
				Ω(record.Path).Should(Equal("/john"))
				Ω(record.Text).Should(Equal("smith"))
				Ω(record.Write).Should(BeTrue())
			})
		})

		Describe("recorder Read", func() {
			It("should respond to a read with response", func() {
				recorder := FileRecorder{}
				recorder.Respond("jane", &StaticResponder{Text: "doe"})
				text, err := recorder.Read("jane")
				Ω(err).Should(BeNil())
				Ω(text).Should(Equal("doe"))

				records := recorder.Records()

				// Verify records
				Ω(records).Should(HaveLen(1))
				record := records[0]
				Ω(record.Stamp).ShouldNot(BeNil())
				Ω(record.Err).Should(BeNil())
				Ω(record.Path).Should(Equal("/jane"))
				Ω(record.Text).Should(Equal("doe"))
				Ω(record.Write).Should(BeFalse())
			})

			It("should respond to a read with no response with an error", func() {
				recorder := FileRecorder{}
				text, err := recorder.Read("foo")
				Ω(err).Should(Equal(os.ErrNotExist))
				Ω(text).Should(BeEmpty())
				recorder.Respond("jane", &StaticResponder{Text: "doe"})
				text, err = recorder.Read("foo")
				Ω(err).Should(Equal(os.ErrNotExist))
				Ω(text).Should(BeEmpty())
				text, err = recorder.Read("jane")
				Ω(err).Should(BeNil())
				Ω(text).Should(Equal("doe"))
				text, err = recorder.Read("/jane")
				Ω(err).Should(BeNil())
				Ω(text).Should(Equal("doe"))
				records := recorder.Records()

				// Verify records
				Ω(records).Should(HaveLen(4))
				record := records[0]
				Ω(record.Stamp).ShouldNot(BeNil())
				Ω(record.Err).Should(Equal(os.ErrNotExist))
				Ω(record.Path).Should(Equal("/foo"))
				Ω(record.Text).Should(BeEmpty())
				Ω(record.Write).Should(BeFalse())
				record = records[1]
				Ω(record.Stamp).ShouldNot(BeNil())
				Ω(record.Err).Should(Equal(os.ErrNotExist))
				Ω(record.Path).Should(Equal("/foo"))
				Ω(record.Text).Should(BeEmpty())
				Ω(record.Write).Should(BeFalse())
				record = records[2]
				Ω(record.Stamp).ShouldNot(BeNil())
				Ω(record.Err).Should(BeNil())
				Ω(record.Path).Should(Equal("/jane"))
				Ω(record.Text).Should(Equal("doe"))
				Ω(record.Write).Should(BeFalse())
				record = records[3]
				Ω(record.Stamp).ShouldNot(BeNil())
				Ω(record.Err).Should(BeNil())
				Ω(record.Path).Should(Equal("/jane"))
				Ω(record.Text).Should(Equal("doe"))
				Ω(record.Write).Should(BeFalse())
			})
		})
	})

	Describe("AbsPath", func() {
		It("should generate proper absolute paths", func() {
			Ω(AbsPath("foo")).Should(Equal("/foo"))
			Ω(AbsPath(".")).Should(Equal("/"))
			Ω(AbsPath("..")).Should(Equal("/"))
			Ω(AbsPath("bar/foo/../baz")).Should(Equal("/bar/baz"))
		})
	})
})
