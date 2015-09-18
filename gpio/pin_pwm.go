package gpio

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gopackage/sysfs"
)

// File path constants. BasicPWMPin specific paths require formatting.
const (
	pwmPathBase     = "/sys/class/pwm/pwmchip0"
	pwmPathExport   = pwmPathBase + "/export"
	pwmPathUnexport = pwmPathBase + "/unexport"
)

// BasicPWMPin manages a single PWM pin.
type BasicPWMPin struct {
	Number int              // Number is the pin number for the PWM pin under control
	pin    string           // The pin number as as string
	file   sysfs.DeviceFile // The device file to write to
	// Pre-calculated file path names for pin settings
	pathDuty   string // file path for setting duty cycle
	pathPeriod string // file path for setting pwm period
	pathEnable string // file path for enabling/disabling pwm

	// Automatic scaling of the duty cycle
	limits Limits  // The limits of scaled values
	span   float64 // The distance between the scale limits
	min    float64 // The minimum value of the limits as a float
	period float64 // The current period
}

// NewBasicPWMPin returns a Pin ready for export, enabling and controlling.
func NewBasicPWMPin(pin int, file sysfs.DeviceFile) (*BasicPWMPin, error) {
	p := &BasicPWMPin{Number: pin, file: file}
	p.pin = strconv.Itoa(pin)
	pathPin := fmt.Sprintf(pwmPathBase+"/pwm%d", pin)
	p.pathDuty = pathPin + "/duty_cycle"
	p.pathPeriod = pathPin + "/period"
	p.pathEnable = pathPin + "/enable"
	p.SetLimits(Limits{0, 255})

	// We need to capture the current settings for the period
	// so we can easily scale duty cycles without a lot of
	// i/o.
	period, err := p.Period()
	if err != nil {
		return nil, err
	}
	p.period = float64(period)
	return p, nil
}

// Enable starts the pin actively sending a PWM pulse train.
func (p *BasicPWMPin) Enable() (err error) {
	return p.file.Write(p.pathEnable, "1")
}

// Disable stops the pin BasicPWMPin pulse train.
func (p *BasicPWMPin) Disable() (err error) {
	return p.file.Write(p.pathEnable, "0")
}

// Period returns the BasicPWMPin period.
func (p *BasicPWMPin) Period() (period int, err error) {
	buf, err := p.file.Read(p.pathPeriod)
	if err != nil {
		return 0, err
	}
	per, err := strconv.Atoi(strings.TrimSpace(string(buf)))
	if err != nil {
		return 0, err
	}
	p.period = float64(per)
	return per, nil
}

// SetDuty changes the pwm duty cycle of the pin.
func (p *BasicPWMPin) SetDuty(duty int) (err error) {
	// TODO check for legal bounds on duty value
	return p.file.Write(p.pathDuty, strconv.Itoa(duty))
}

// SetLimits sets the limits that scaled PWM outputs are based on.
// By default, the scale is 0 - 255 to support unsigned byte values
// common for PWM (e.g. for scaling 24-bit RGB colors).
func (p *BasicPWMPin) SetLimits(l Limits) {
	l.Validate()
	p.limits = l
	p.min = float64(l.Min)
	p.span = float64(l.Max) - p.min
}

// Scale sets the duty cycle base on the currently set Limits for the PWM.
func (p *BasicPWMPin) Scale(value int) error {
	duty := (float64(value) - p.min) / p.span
	return p.SetDuty(int(p.period * duty))
}

// Export adds the pin to the exported list making it available for
// reading and writing.
func (p *BasicPWMPin) Export() (err error) {
	return p.file.Write(pwmPathExport, p.pin)
}

// Unexport removes the pin from the exported list making it no longer
// able to read or write.
func (p *BasicPWMPin) Unexport() (err error) {
	return p.file.Write(pwmPathUnexport, p.pin)
}
