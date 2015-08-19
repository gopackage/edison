package file_test

import (
	"io/ioutil"
	"os"
	"path"

	. "../file"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/twinj/uuid"
)

var _ = Describe("Core", func() {

	Describe("DeviceFile", func() {
		It("should implement DeviceFile", func() {
			var f DeviceFile
			f = &HardwareFile{}
			Ω(f).ShouldNot(BeNil())
		})
		It("should read/write to a real file", func() {
			f := HardwareFile{}
			u4 := uuid.NewV4()
			name := uuid.Formatter(u4, uuid.Clean)
			tmp := path.Join(os.TempDir(), name+".gotest.tmp")

			By("using tmp " + tmp)
			// Write to the file
			err := f.Write(tmp, "Hello World")
			Ω(err).Should(BeNil())
			// Read and verify file
			text, err := f.Read(tmp)
			Ω(err).Should(BeNil())
			Ω(text).Should(Equal("Hello World"))
			// Test file existence directly
			found, err := ioutil.ReadFile(tmp)
			Ω(err).Should(BeNil())
			Ω(found).Should(Equal([]byte("Hello World")))
			// Clean up
			err = os.Remove(tmp)
			Ω(err).Should(BeNil())
		})
	})

	Describe("LogFile", func() {
		It("should implement DeviceFile", func() {
			var f DeviceFile
			f = &LogFile{}
			Ω(f).ShouldNot(BeNil())
		})
		It("should respond to reads with error when no proxy present", func() {
			lf := LogFile{}
			text, err := lf.Read("/foo")
			Ω(text).Should(BeEmpty())
			Ω(err).Should(Equal(os.ErrNotExist))
		})
		// TODO add unit tests for writes and regular reads
	})
	// TODO add unit tests for StatsFile
})
