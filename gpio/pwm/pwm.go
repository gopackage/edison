package gpio

import (
	"math"
	"strconv"
)

// File path constants. Pin specific paths require formatting.
const (
	pathBase     = "/sys/class/pwm/pwmchip0"
	pathExport   = pathBase + "/export"
	pathUnexport = pathBase + "/unexport"
	pathPin      = pathBase + "/pwm%s"
	pathDuty     = pathPin + "/duty_cycle"
	pathPeriod   = pathPin + "/period"
	pathEnable   = pathPin + "/enable"
)

// Pin manages a single PWM pin.
type Pin struct {
	Number string // Number is the pin number for the PWM pin under control
}

// NewPin returns a Pin ready for export, enabling and controlling.
func NewPin(pin int) *Pin {
	return &Pin{Number: strconv.Itoa(pin)}
}

// enable writes value to pwm enable path
func (p *Pin) enable(val string) (err error) {
	//_, err = writeFile(pwmEnablePath(p.pin), val)
	return
}

// period reads from pwm period path and returns value
func (p *Pin) period() (period string, err error) {
	/*
		buf, err := readFile(pwmPeriodPath(p.pin))
		if err != nil {
			return
		}
		return string(buf[0 : len(buf)-1]), nil
	*/
	return "", nil
}

// writeDuty writes value to pwm duty cycle path
func (p *Pin) writeDuty(duty string) (err error) {
	//_, err = writeFile(pwmDutyCyclePath(p.pin), duty)
	return
}

// export writes pin to pwm export path
func (p *Pin) export() (err error) {
	//_, err = writeFile(pwmExportPath(), p.pin)
	return
}

// export writes pin to pwm unexport path
func (p *Pin) unexport() (err error) {
	//_, err = writeFile(pwmUnExportPath(), p.pin)
	return
}

// FromScale returns a converted input from min, max to 0.0...1.0.
func FromScale(input, min, max float64) float64 {
	return (input - math.Min(min, max)) / (math.Max(min, max) - math.Min(min, max))
}
