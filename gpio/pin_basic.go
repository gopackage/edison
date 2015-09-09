package gpio

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/gopackage/sysfs"
)

// BasicOutputPin is a simple implementation of the DigitalOutputPin interface. The BasicOutputPin
// does not cache or reuse file handles or any other optimizations for
// higher performance. High pin access rates may want to rely on a more
// complex higher performance implementation.
type BasicOutputPin struct {
	pin   int
	rPin  int
	label string
	file  sysfs.DeviceFile

	// Cached paths
	pathExport   string
	pathUnexport string
	pathValue    string
}

// Export sets up the pin so it is useable for reads/writes.
func (p *BasicOutputPin) Export() (err error) {
	fmt.Printf("Exporting pin: %d\n", p.pin)
	if err := p.file.Write(GPIOPATH+"/export", strconv.Itoa(p.pin)); err != nil {
		// If EBUSY then the pin has already been exported
		if err.(*os.PathError).Err != syscall.EBUSY {
			return err
		}
	}
	return nil
}

// Unexport tears down the pin so it is no longer useable for reads/writes.
func (p *BasicOutputPin) Unexport() (err error) {
	if err := p.file.Write(GPIOPATH+"/unexport", strconv.Itoa(p.pin)); err != nil {
		// If EINVAL then the pin is reserved in the system and can't be unexported
		if err.(*os.PathError).Err != syscall.EINVAL {
			return err
		}
	}
	return nil
}

// SetHigh sets the digital pin value high.
func (p *BasicOutputPin) SetHigh() (err error) {
	str := fmt.Sprintf("%v/%s/value", GPIOPATH, p.label)
	err = p.file.Write(str, "1")
	return err
}

// SetLow sets the digital pin value low.
func (p *BasicOutputPin) SetLow() (err error) {
	str := fmt.Sprintf("%v/%s/value", GPIOPATH, p.label)
	err = p.file.Write(str, "0")
	return err
}

// BasicInputPin is a simple implementation of the DigitalInputPin interface. The BasicInputPin
// does not cache or reuse file handles or any other optimizations for
// higher performance. High pin access rates may want to rely on a more
// complex higher performance implementation.
type BasicInputPin struct {
	pin   int
	rPin  int
	label string
	file  sysfs.DeviceFile

	// Cached paths
	pathExport   string
	pathUnexport string
	pathValue    string
}

// Export sets up the pin so it is useable for reads/writes.
func (p *BasicInputPin) Export() (err error) {
	fmt.Printf("Exporting pin: %d\n", p.pin)
	if err := p.file.Write(GPIOPATH+"/export", strconv.Itoa(p.pin)); err != nil {
		// If EBUSY then the pin has already been exported
		if err.(*os.PathError).Err != syscall.EBUSY {
			return err
		}
	}
	return nil
}

// Unexport tears down the pin so it is no longer useable for reads/writes.
func (p *BasicInputPin) Unexport() (err error) {
	if err := p.file.Write(GPIOPATH+"/unexport", strconv.Itoa(p.pin)); err != nil {
		// If EINVAL then the pin is reserved in the system and can't be unexported
		if err.(*os.PathError).Err != syscall.EINVAL {
			return err
		}
	}
	return nil
}

// High determines whether the pin is high (true).
func (p *BasicInputPin) High() (isHigh bool, err error) {
	buf, err := p.file.Read(fmt.Sprintf("%v/%d/value", GPIOPATH, p.pin))
	if err != nil {
		return false, err
	}
	return string(buf[0]) == "0", err
}

// Low determines whether the pin is low (true).
func (p *BasicInputPin) Low() (isLow bool, err error) {
	buf, err := p.file.Read(fmt.Sprintf("%v/%d/value", GPIOPATH, p.pin))
	if err != nil {
		return false, err
	}
	return string(buf[0]) == "0", err
}
